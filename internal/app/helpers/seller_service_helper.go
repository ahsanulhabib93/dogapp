package helpers

import (
	"context"
	"encoding/json"
	"fmt"

	"strconv"
	"strings"

	cmtPb "github.com/voonik/goConnect/api/go/cmt/product"
	omsPb "github.com/voonik/goConnect/api/go/oms/seller"
	"github.com/voonik/goConnect/api/go/vigeon/notify"

	"github.com/shopuptech/go-libs/logger"
	spb "github.com/voonik/goConnect/api/go/ss2/seller"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/utils"
)

type ReturnExchangePolicy struct {
	Return   ExchangeDetails `json:"return"`
	Exchange ExchangeDetails `json:"exchange"`
}

// ExchangeDetails represents the details for return or exchange.
type ExchangeDetails struct {
	DefaultDuration      int    `json:"default_duration"`
	ReturnDaysStartsFrom string `json:"return_days_starts_from"`
}

func PerformSendActivationMail(ctx context.Context, sellerIDs []uint64, params *spb.SendActivationMailParams) *spb.BasicApiResponse {
	resp := &spb.BasicApiResponse{Status: utils.Failure}
	sellerDetails := GetSellerByIds(ctx, sellerIDs)
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

func PerformApproveProductFunc(ctx context.Context, param *spb.ApproveProductsParams) *spb.BasicApiResponse {
	resp := &spb.BasicApiResponse{Status: utils.Failure}
	seller := GetSellerByUserId(ctx, param.GetId())
	if seller.ID == utils.Zero {
		resp.Message = "Seller Not Found"
	} else {
		if len(seller.VendorAddresses) > utils.Zero && seller.PanNumber != utils.EmptyString && seller.ActivationState != 5 {
			itemCountResp := getAPIHelperInstance().CmtApproveItems(ctx, &cmtPb.ApproveItemParams{ProductIds: param.GetIds(), State: uint64(seller.ActivationState), UserId: seller.UserID})
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

func GetArrayIdsFromString(id string) (string, []uint64) {
	params := map[string]string{"id": id}
	stringIDs := strings.Split(params["id"], ",")

	sellerIDs := make([]uint64, len(stringIDs))
	for i, strID := range stringIDs {
		trimmedStrID := strings.TrimSpace(strID)
		id, err := strconv.ParseUint(trimmedStrID, utils.Ten, utils.SixtyFour)
		if err != nil {
			return fmt.Sprintf("Error converting string to uint64: %+v", err), []uint64{}
		}
		sellerIDs[i] = id
	}
	return "", sellerIDs
}

func createSeller(ctx context.Context, params *spb.CreateParams) (*models.Seller, error) {
	returnExchangePolicy, _ := json.Marshal(DefaultsellerReturnExchangePolicy())
	jsonDataMapping, _ := json.Marshal(utils.SellerDataMapping)
	sellerPricingDetails := &models.SellerPricingDetail{}

	seller := &models.Seller{
		UserID:                 params.Seller.UserId,
		BrandName:              params.Seller.BrandName,
		CompanyName:            params.Seller.CompanyName,
		PrimaryEmail:           params.Seller.PrimaryEmail,
		PrimaryPhone:           params.Seller.PrimaryPhone,
		ActivationState:        utils.ActivationState(params.Seller.ActivationState),
		Slug:                   params.Seller.Slug,
		Hub:                    params.Seller.Hub,
		Slot:                   params.Seller.Slot,
		DeliveryType:           int(params.Seller.DeliveryType),
		ProcessingType:         int(params.Seller.ProcessingType),
		BusinessUnit:           int(params.Seller.BusinessUnit),
		FullfillmentType:       int(params.Seller.FullfillmentType),
		ColorCode:              utils.ColorCode(params.Seller.ColorCode),
		TinNumber:              params.Seller.TinNumber,
		SellerCloseDay:         params.Seller.SellerCloseDay,
		AcceptedPaymentMethods: params.Seller.AcceptedPaymentMethods,
		AffiliateURL:           utils.DefaultAffiliateURL,
		IsDirect:               true,
		ReturnExchangePolicy:   returnExchangePolicy,
		DataMapping:            jsonDataMapping,
		AggregatorID:           int(params.Seller.UserId),
		SellerPricingDetails:   []*models.SellerPricingDetail{sellerPricingDetails}, // Taking values from DB defaults
		SellerConfig:           createSellerDefaultSellerConfig(),
	}

	if len(params.Seller.VendorAddresses) != utils.Zero {
		seller.VendorAddresses = assignVendorAddressData(params.Seller.VendorAddresses)
	}

	err := database.DBAPM(ctx).Model(&models.Seller{}).Create(seller).Error
	if err != nil {
		return nil, err
	}

	return seller, nil
}

func assignVendorAddressData(vendorAddressesObjects []*spb.VendorAddressObject) []*models.VendorAddress {
	vendorAddresses := []*models.VendorAddress{}
	for _, vendorAddress := range vendorAddressesObjects {
		vendorAddresses = append(vendorAddresses, &models.VendorAddress{
			Firstname:            vendorAddress.Firstname,
			Lastname:             vendorAddress.Lastname,
			Address1:             vendorAddress.Address1,
			Address2:             vendorAddress.Address2,
			City:                 vendorAddress.City,
			Zipcode:              vendorAddress.Zipcode,
			AlternativePhone:     vendorAddress.AlternativePhone,
			Company:              vendorAddress.Company,
			State:                vendorAddress.State,
			Country:              vendorAddress.Country,
			AddressType:          int(vendorAddress.AddressType),
			DefaultAddress:       vendorAddress.DefaultAddress,
			AddressProofFileName: vendorAddress.AddressProofFileName,
			VerificationStatus:   utils.VerificationStatus(vendorAddress.VerificationStatus),
			UUID:                 vendorAddress.Uuid,
			ExtraData:            vendorAddress.ExtraData,
		})
	}
	return vendorAddresses
}

func createSellerDefaultSellerConfig() *models.SellerConfig {
	refundPolicy, _ := json.Marshal(map[string]int{
		"cod":           1,
		"payu_redirect": 1,
	})
	sellerConfig := &models.SellerConfig{
		ItemsPerPackage:       int(utils.DefaultSellerItemsPerPackage),
		MaxQuantity:           int(utils.DefaultSellerMaxQuantity),
		SellerStockEnabled:    true,
		CODConfirmationNeeded: true,
		AllowPriceUpdate:      true,
		PickupType:            int(utils.DefaultSellerPickupType),
		AllowVendorCoupons:    true,
		RefundPolicy:          refundPolicy,
	}
	return sellerConfig
}

func GetDefaultSellerConfig() *spb.SellerConfig {
	return &spb.SellerConfig{
		ItemsPerPackage:       utils.DefaultSellerItemsPerPackage,
		MaxQuantity:           utils.DefaultSellerMaxQuantity,
		SellerStockEnabled:    true,
		CodConfirmationNeeded: true,
		AllowPriceUpdate:      true,
		PickupType:            utils.DefaultSellerPickupType,
		AllowVendorCoupons:    true,
	}
}

func defaultsellerReturnExchangePolicyConfig() *spb.SellerReturnExchangePolicyConfig {
	return &spb.SellerReturnExchangePolicyConfig{
		ReturnDaysStartsFrom: "delivery",
		DefaultDuration:      uint64(15),
	}
}

func DefaultsellerReturnExchangePolicy() *spb.ReturnExchangePolicy {
	exchangeConfig := defaultsellerReturnExchangePolicyConfig()
	return &spb.ReturnExchangePolicy{
		Return:   exchangeConfig,
		Exchange: exchangeConfig,
	}
}

func SendEmailNotification(ctx context.Context, seller *models.Seller, agentEmail string, err error, response *spb.CreateResponse) *notify.EmailResp {
	var subject, content string
	toEmail := "Mokam<noreply@shopf.co>"
	fromEmail := "smk@shopf.co"

	if err != nil {
		subject = "New Seller Registration Failed"
		content = fmt.Sprintf("Seller Registration failed because of the error <br><br> %s <br><br> %s <br><br> with the response <br><br> %s", err.Error(), strings.Join(strings.Split(fmt.Sprintf("%+v", err), "\n"), "<br>"), response)
	} else {
		subject = "New Seller Registered successfully"
		content = fmt.Sprintf("Seller Registered (<b>email:</b> %s, <b>agent_email:</b> %s ) with the response <br><br> %s", seller.PrimaryEmail, agentEmail, response)
	}

	emailParam := notify.EmailParam{
		ToEmail:   toEmail,
		FromEmail: fromEmail,
		Subject:   subject,
		Content:   content,
	}

	return getVigeonAPIHelperInstance().SendEmailAPI(ctx, emailParam)
}

func ValidateSellerParams(params *spb.CreateParams) error {
	if params.Seller == nil {
		return fmt.Errorf("Missing Seller Params")
	}
	return nil
}

func ProcessSellerRegistration(ctx context.Context, params *spb.CreateParams) (*models.Seller, string, error) {
	registrationMessage := "Seller registered successfully."
	existingSeller := GetSellerByUserId(ctx, params.Seller.UserId)
	if existingSeller.ID != utils.Zero {
		registrationMessage = "Seller already registered."
		return existingSeller, registrationMessage, nil
	}
	newSeller, err := createSeller(ctx, params)
	return newSeller, registrationMessage, err
}

func getSellerPricingDetailsSum(sellerPricingDetails []*models.SellerPricingDetail) (uint64, uint64) {
	var totalLeadShippingDays, totalCommissionPercent uint64
	for _, pricingDetail := range sellerPricingDetails {
		totalLeadShippingDays += uint64(pricingDetail.LeadShippingDays)
		totalCommissionPercent += uint64(pricingDetail.CommissionPercent)
	}
	return totalLeadShippingDays, totalCommissionPercent
}

func getVendorAddressData(params *models.Seller) []*omsPb.VendorAddressObject {
	vendorAddresses := params.VendorAddresses
	vendorAddressData := []*omsPb.VendorAddressObject{}
	for _, vendorAddress := range vendorAddresses {
		vendorAddressData = append(vendorAddressData, &omsPb.VendorAddressObject{
			SellerId:           params.UserID,
			Firstname:          vendorAddress.Firstname,
			Lastname:           vendorAddress.Lastname,
			Address1:           vendorAddress.Address1,
			Address2:           vendorAddress.Address2,
			City:               vendorAddress.City,
			Zipcode:            vendorAddress.Zipcode,
			Phone:              fmt.Sprint(vendorAddress.Phone),
			AlternativePhone:   vendorAddress.AlternativePhone,
			Company:            vendorAddress.Company,
			State:              vendorAddress.State,
			Country:            vendorAddress.Country,
			AddressType:        uint64(vendorAddress.AddressType),
			DefaultAddress:     vendorAddress.DefaultAddress,
			VerificationStatus: string(vendorAddress.VerificationStatus),
			ExtraData:          vendorAddress.ExtraData,
		})
	}
	return vendorAddressData
}

func CreateOMSSellerSync(ctx context.Context, params *models.Seller) (err error) {
	var returnExchangePolicy ReturnExchangePolicy
	if err := json.Unmarshal([]byte(params.ReturnExchangePolicy), &returnExchangePolicy); err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return err
	}
	totalLeadShippingDays, totalCommissionPercent := getSellerPricingDetailsSum(params.SellerPricingDetails)
	sellerParam := omsPb.SellerParams{
		Seller: &omsPb.Seller{
			UserId:           params.UserID,
			FullfillmentType: uint64(params.FullfillmentType),
			BrandName:        params.BrandName,
			CompanyName:      params.CompanyName,
			Slug:             params.Slug,
			AggregatorId:     uint64(params.AggregatorID),
			PrimaryEmail:     params.PrimaryEmail,
			PrimaryPhone:     params.PrimaryPhone,
			ReturnExchangePolicy: &omsPb.ReturnExchangePolicy{
				Return: &omsPb.SellerReturnExchangePolicyConfig{
					ReturnDaysStartsFrom: returnExchangePolicy.Return.ReturnDaysStartsFrom,
					DefaultDuration:      uint64(returnExchangePolicy.Return.DefaultDuration),
				},
				Exchange: &omsPb.SellerReturnExchangePolicyConfig{
					ReturnDaysStartsFrom: returnExchangePolicy.Exchange.ReturnDaysStartsFrom,
					DefaultDuration:      uint64(returnExchangePolicy.Exchange.DefaultDuration),
				},
			},
			TinNumber:             params.TinNumber,
			PanNumber:             params.PanNumber,
			SellerInvoiceNumber:   uint64(params.SellerInvoiceNumber),
			SellerType:            uint64(params.SellerType),
			SellerRating:          float32(params.SellerRating),
			SellerId:              params.ID,
			DeliveryType:          uint64(params.DeliveryType),
			ProcessingType:        uint64(params.ProcessingType),
			BusinessUnit:          uint64(params.BusinessUnit),
			CodConfirmationNeeded: params.SellerConfig.CODConfirmationNeeded,
			RefundPolicy:          params.SellerConfig.RefundPolicy.String(),
			PenaltyPolicy:         params.SellerConfig.PenaltyPolicy,
			ItemsPerPackage:       uint64(params.SellerConfig.ItemsPerPackage),
			PickupType:            uint64(params.SellerConfig.PickupType),
			QcFrequency:           uint64(params.SellerConfig.QCFrequency),
			LeadShippingDays:      totalLeadShippingDays,
			CommissionPercent:     totalCommissionPercent,
		},
		VendorAddresses: getVendorAddressData(params),
	}
	resp := getAPIHelperInstance().CreateOmsSeller(ctx, &sellerParam)
	if !resp.Success {
		return fmt.Errorf(resp.Message)
	}
	return nil
}
