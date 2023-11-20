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
