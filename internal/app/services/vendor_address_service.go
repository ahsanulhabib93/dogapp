package services

import (
	"context"
	"fmt"

	"github.com/shopuptech/go-libs/logger"
	vapb "github.com/voonik/goConnect/api/go/ss2/vendor_address"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/utils"
)

type VendorAddressService struct{}

// def get_data
// 	params.permit!
// 	seller_id = params[:seller_user_id] ? Seller.select("id").find_by(user_id: params[:seller_user_id]) : params[:seller_id]
// 	vendor_addresses = VendorAddress.where(:seller_id => seller_id) if seller_id.present?
// 	vendor_addresses = VendorAddress.unscoped.where(:id => params[:id]) if params[:id].present?
// 	vendor_addresses = VendorAddress.where(:id => params[:scoped_id]) if params[:scoped_id].present?
// 	vendor_addresses
// end

func (vapb *VendorAddressService) GetData(ctx context.Context, params *vapb.GetDataParams) (*vapb.GetDataResponse, error) {
	response := vapb.GetDataResponse{Status: utils.Success}
	userIds := params.GetUserIds()
	sellerIds := params.GetSellerIds()
	ids := params.GetIds()

	if len(userIds) == 0 && len(sellerIds) == 0 && len(ids) == 0 {
		response.Status = utils.Failure
		response.Message = "param not specified"
		return &response, nil
	}

	// query := database.DBAPM(ctx).Model(&models.Seller{}).Where("user_id = ?", userId)
	query := database.DBAPM(ctx).Model(&models.VendorAddress{})
	if len(ids) != 0 {
		query = query.Where("id in (?)", ids)
	}
	if len(sellerIds) != 0 {
		query = query.Where("seller_id in (?)", sellerIds)
	}
	if len(userIds) != 0 {
		sellers := []*models.Seller{}
		err := database.DBAPM(ctx).Model(&models.Seller{}).Where("user_id in (?)", userIds).Scan(&sellers).Error
		if err != nil {
			logger.FromContext(ctx).Info("Error in vendor address service GetData API")
			response.Status = utils.Failure
			response.Message = fmt.Sprint("Error in vendor address service GetData API: ", err.Error())
			return &response, nil
		}
		if len(sellers) == 0 {
			response.Message = "seller not found with the user id"
			return &response, nil
		}
	}
	vendorAddress := []*models.VendorAddress{}
	err := query.Scan(vendorAddress).Error
	if err != nil {
		logger.FromContext(ctx).Info("Error in vendor address service GetData API: ", err.Error())
		response.Status = utils.Failure
		response.Message = fmt.Sprint("Error in vendor address service GetData API: ", err.Error())
		return &response, nil
	}
	if len(vendorAddress) == 0 {
		response.Message = "vendor address not found"
		return &response, nil
	}
	response.Message = "fetched vendor address successfully"
	return &response, nil
}

func (vapb *VendorAddressService) VerifyAddress(ctx context.Context, params *vapb.VerifyAddressParams) (*vapb.BasicApiResponse, error) {
	return nil, nil
}
