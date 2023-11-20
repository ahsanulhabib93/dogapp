package services

import (
	"context"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/shopuptech/go-libs/logger"
	vapb "github.com/voonik/goConnect/api/go/ss2/vendor_address"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/utils"
)

type VendorAddressService struct{}

func (va *VendorAddressService) GetData(ctx context.Context, params *vapb.GetDataParams) (*vapb.GetDataResponse, error) {
	return nil, nil
}

func (vas *VendorAddressService) VerifyAddress(ctx context.Context, params *vapb.VerifyAddressParams) (*vapb.BasicApiResponse, error) {
	response := vapb.BasicApiResponse{Status: utils.Failure}
	uuid := params.GetId()
	if uuid == "" {
		response.Message = "param not specified"
		return &response, nil
	}
	vendorAddress := models.VendorAddress{}
	query := database.DBAPM(ctx).Model(&models.VendorAddress{}).Where("uuid = ?", uuid)
	err := query.Scan(&vendorAddress).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			response.Message = "vendor address not found"
			return &response, nil
		}
		response.Message = fmt.Sprint("error in vendor address service VerifyAddress API1: ", err.Error())
		logger.FromContext(ctx).Error(response.Message)
		return &response, nil
	}
	err = query.Update("verification_status", utils.Verified).Error
	if err != nil {
		response.Message = fmt.Sprint("not able to update verification status: ", err.Error())
		logger.FromContext(ctx).Error(response.Message)
		return &response, nil
	}
	sellerID := vendorAddress.SellerID
	seller := models.Seller{}
	err = database.DBAPM(ctx).Model(&models.Seller{}).Where("id = ?", sellerID).Scan(&seller).Error
	if err != nil {
		response.Message = fmt.Sprint("error in vendor address service VerifyAddress API2: ", err.Error())
		logger.FromContext(ctx).Error(response.Message)
		return &response, nil
	}
	sellerActivityLog := models.SellerActivityLog{
		SellerID: seller.UserID,
		Action:   "verify_address",
		Notes:    `{"status":"verified"}`,
	}
	if uid := utils.GetCurrentUserID(ctx); uid != nil {
		sellerActivityLog.UserID = *uid
	}
	err = database.DBAPM(ctx).Model(&models.SellerActivityLog{}).Create(&sellerActivityLog).Error
	if err != nil {
		response.Message = fmt.Sprint("error in vendor address service VerifyAddress API3: ", err.Error())
		logger.FromContext(ctx).Error(response.Message)
		return &response, nil
	}
	response.Status = utils.Success
	response.Message = "vendor address verified successfully"
	return &response, nil
}
