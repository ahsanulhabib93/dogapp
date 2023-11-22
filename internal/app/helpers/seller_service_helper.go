package helpers

import (
	"context"
	"fmt"

	cmtPb "github.com/voonik/goConnect/api/go/cmt/product"
	"strconv"
	"strings"

	"github.com/shopuptech/go-libs/logger"
	spb "github.com/voonik/goConnect/api/go/ss2/seller"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/utils"
)

func PerformSendActivationMail(ctx context.Context, params *spb.SendActivationMailParams) *spb.BasicApiResponse {
	resp := &spb.BasicApiResponse{Status: utils.Failure}
	sellerDetails := GetSellerByIds(ctx, params.GetIds())
	var noAccess []uint64
	if len(sellerDetails) == utils.Zero {
		resp.Message = "Seller not found"
	} else {
		for _, seller := range sellerDetails {
			sellerBankDetails := GetSellerBankDetails(ctx, seller)
			if seller.PanNumber != utils.EmptyString && seller.EmailConfirmed && seller.MouAgreed && len(sellerBankDetails) > utils.Zero {
				var successfulStateChanges int
				resp, successfulStateChanges = VerifyVendorAddress(ctx, seller, params.GetAction())
				if resp.Status == utils.Success && successfulStateChanges > utils.One {
					resp.Message = fmt.Sprintf("%d Seller accounts activated successfully.", successfulStateChanges)
				}
				if seller.StateReason > utils.Zero && seller.ActivationState > utils.Zero {
					noAccess = FindNonAccessSellers(params, seller)
					if len(noAccess) > utils.Zero {
						noAccessStr := utils.GetArrIntToArrStr(noAccess)
						resp.Message += " You don't have access to activate this Seller(s) - " + strings.Join(noAccessStr, ",")
					}
				}
			} else {
				resp.Message += strconv.Itoa(int(seller.UserID)) + ": Seller Pan Number, Bank Detail, MOU and Email should be confirmed."
			}
		}
	}
	return resp
}

func VerifyVendorAddress(ctx context.Context, seller *models.Seller, action string) (*spb.BasicApiResponse, int) {
	resp := &spb.BasicApiResponse{Status: utils.Failure}
	successfulStateChanges := utils.Zero
	vendorAddresses, verifiedCount, defaultCount := GetVendorAddressBySellerID(ctx, seller.ID)
	addressCount := len(vendorAddresses)
	if verifiedCount == utils.Zero && addressCount > utils.One {
		resp.Message += fmt.Sprintf("%s: Make at least one address as verified.", strconv.Itoa(int(seller.UserID)))
	} else if defaultCount == utils.Zero && addressCount > utils.One {
		resp.Message += fmt.Sprintf("%s: Make at least one address as default.", strconv.Itoa(int(seller.UserID)))
	} else if addressCount == utils.Zero {
		resp.Message += fmt.Sprintf("%s: At least one address should be present.", strconv.Itoa(int(seller.UserID)))
	} else if len(seller.SellerPricingDetails) == utils.Zero {
		resp.Message += fmt.Sprintf("%s: Seller pricing details are not present.", strconv.Itoa(int(seller.UserID)))
	} else if IsSellerPricingDetailsNotVerified(ctx, seller.SellerPricingDetails[utils.Zero]) {
		resp.Message += fmt.Sprintf("%s: Seller pricing details are not verified.", strconv.Itoa(int(seller.UserID)))
	} else {
		if addressCount == utils.One {
			vendorAddresses[utils.Zero].VerificationStatus = "VERIFIED"
			vendorAddresses[utils.Zero].DefaultAddress = true
			database.DBAPM(ctx).Save(&vendorAddresses[utils.Zero])
		}
		var err error
		resp, err = ActivateSeller(ctx, *seller)
		if err != nil {
			// NewRelic::Agent.notice_error(err)
			logger.Log().Errorf("Error during seller Activation for %s. Issue - %s\n", seller.UserID, err.Error())
			resp.Message += fmt.Sprintf("%s activation failed - %s", strconv.Itoa(int(seller.UserID)), err.Error())
		} else {
			if resp.Status == utils.Success {
				successfulStateChanges += 1
				CreateSellerActivityLog(ctx, seller.ID, action)
			} else {
				resp.Message += fmt.Sprintf("%s: %s", strconv.Itoa(int(seller.UserID)), resp.Message)
			}
		}
	}
	return resp, successfulStateChanges
}

func FindNonAccessSellers(params *spb.SendActivationMailParams, seller *models.Seller) []uint64 {
	var noAccess []uint64
	activationState, stateReason := seller.ActivationState, seller.StateReason
	isQualityTeam := params.GetIsQualityTeam()
	isRiskTeam := params.GetIsRiskTeam()
	isSellerOnboardingTeam := params.GetIsSellerOnboardingTeam()
	if isQualityTeam || isRiskTeam {
		condition := (stateReason == utils.PRODUCT_QUALITY && isQualityTeam) || (stateReason != utils.PRODUCT_QUALITY && isRiskTeam)
		condition = condition && CheckRestrictiveSellerState(activationState)
		if !condition {
			noAccess = append(noAccess, seller.ID)
		}
	} else if isSellerOnboardingTeam && (!isQualityTeam || !isRiskTeam) {
		if !SellerIsOnboardingState(activationState) && !SellerIsOnboardingStateReason(stateReason) {
			noAccess = append(noAccess, seller.ID)
		}
	}
	return noAccess
}

func GetSellerByIds(ctx context.Context, userIds []uint64) []*models.Seller {
	sellerDetails := []*models.Seller{}
	database.DBAPM(ctx).Preload("SellerPricingDetails").Model(&models.Seller{}).Where("user_id in (?)", userIds).Find(&sellerDetails)
	return sellerDetails
}

func GetSellerBankDetails(ctx context.Context, seller *models.Seller) []*models.SellerBankDetail {
	sellerBankDetails := []*models.SellerBankDetail{}
	database.DBAPM(ctx).Model(&models.SellerBankDetail{}).Where("seller_id = ? and deleted_at is NULL", seller.ID).Scan(&sellerBankDetails)
	return sellerBankDetails
}

func GetVendorAddressBySellerID(ctx context.Context, sellerID uint64) ([]models.VendorAddress, uint64, uint64) {
	vendorAddress := []models.VendorAddress{}
	query := database.DBAPM(ctx).Model(models.VendorAddress{}).Where(
		"seller_id = ? and gst_status is not NULL and deleted_at is NULL", sellerID)
	query.Scan(&vendorAddress)

	var defaultAddressCount, verifiedStatusCount uint64
	query.Where("default_address = ?", true).Count(&defaultAddressCount)
	query.Where("verification_status = ?", utils.Verified).Count(&verifiedStatusCount)
	return vendorAddress, verifiedStatusCount, defaultAddressCount
}

func CheckRestrictiveSellerState(sellerState utils.ActivationState) bool {
	return sellerState == utils.SUSPENDED || sellerState == utils.BLOCKED || sellerState == utils.UNDER_REVIEW || sellerState == utils.FRAUD || sellerState == utils.ON_HOLD || sellerState == utils.HOLD_OFF
}

func SellerIsOnboardingState(activationState utils.ActivationState) bool {
	return activationState == utils.NOT_ACTIVATED || activationState == utils.VERIFICATION_PENDING || activationState == utils.HOLD_OFF || activationState == utils.VACATION_PENDING || activationState == utils.GST_PENDING || activationState == utils.UNDER_REVIEW

}

func SellerIsOnboardingStateReason(stateReason utils.StateReason) bool {
	return stateReason == utils.PENDING_CONTACT_WITH_SS || stateReason == utils.VACATION_MODE
}

func IsSellerPricingDetailsNotVerified(ctx context.Context, sellerPrice *models.SellerPricingDetail) bool {
	return sellerPrice.Verified == utils.SellerPriceVerified(utils.NotVerified)
}

func ActivateSeller(ctx context.Context, seller models.Seller) (*spb.BasicApiResponse, error) {
	resp := spb.BasicApiResponse{Status: utils.Success, Message: "Seller account activated successfully."}
	seller.ActivationState, seller.StateReason = utils.ACTIVATED, 0
	database.DBAPM(ctx).Save(&seller)
	return &resp, nil
}

func CreateSellerActivityLog(ctx context.Context, sellerID uint64, action string) {
	var currentUserId uint64
	if v := utils.GetCurrentUserID(ctx); v != nil {
		currentUserId = *v
	}

	activityLog := models.SellerActivityLog{
		UserID:   currentUserId,
		SellerID: sellerID,
		Action:   action,
		Notes:    `"reason": "Activation without mail"`,
	}
	database.DBAPM(ctx).Create(&activityLog)
}

func PerformApproveProductFunc(ctx context.Context, ids []uint64) *spb.BasicApiResponse {
	resp := &spb.BasicApiResponse{Status: utils.Failure}
	seller := GetSellerByUserId(ctx, *utils.GetCurrentUserID(ctx))
	if seller.ID == utils.Zero {
		resp.Message = "Seller Not Found"
	} else {
		if len(seller.VendorAddresses) > utils.Zero && seller.PanNumber != utils.EmptyString && seller.ActivationState != 5 {
			itemCountResp := getAPIHelperInstance().CmtApproveItems(ctx, &cmtPb.ApproveItemParams{ProductIds: ids, State: uint64(seller.ActivationState), UserId: seller.UserID})
			resp.Status, resp.Message = utils.Success, fmt.Sprintf("The total number of products approved are %d", itemCountResp.GetCount())
		} else {
			resp.Message = "Pick Up Address or Pan number is missing"
		}
	}
	return resp
}

func GetSellerByUserId(ctx context.Context, userID uint64) *models.Seller {
	sellerData := models.Seller{}
	database.DBAPM(ctx).Preload("VendorAddresses").Model(&models.Seller{}).Where("user_id = ?", userID).Find(&sellerData)
	return &sellerData
}
