package helpers

import (
	"context"
	"fmt"

	cmtPb "github.com/voonik/goConnect/api/go/cmt/product"
	spb "github.com/voonik/goConnect/api/go/ss2/seller"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/utils"
)

func PerformApproveProductFunc(ctx context.Context, ids []uint64) *spb.BasicApiResponse {
	resp := &spb.BasicApiResponse{}
	seller := GetSellerById(ctx, *utils.GetCurrentUserID(ctx))
	vendorAddress, _, _ := GetVendorAddressBySellerID(ctx, seller.ID)
	if seller.ID != utils.Zero && len(vendorAddress) > utils.Zero && seller.PanNumber != utils.EmptyString && (seller.ActivationState != utils.HOLD_OFF) {
		itemCountResp := getAPIHelperInstance().CmtApproveItems(ctx, &cmtPb.ApproveItemParams{ProductIds: ids, State: uint64(seller.ActivationState), UserId: seller.UserID})
		if itemCountResp.Count > utils.Zero {
			resp.Status, resp.Message = utils.Success, fmt.Sprintf("The total number of products approved are %d", itemCountResp.Count)
		} else {
			resp.Message = fmt.Sprintf("The total number of products approved are %d", itemCountResp.Count)
		}
	} else {
		resp.Message = "Pick Up Address or Pan number is missing"
	}
	return resp
}

func GetSellerById(ctx context.Context, userID uint64) *models.Seller {
	sellerData := models.Seller{}
	database.DBAPM(ctx).Model(&models.Seller{}).Find("id = ?", userID).Find(&sellerData)
	return &sellerData
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
