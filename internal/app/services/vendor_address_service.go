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

func ParamCount(params ...[]uint64) int {
	count := 0
	for _, param := range params {
		if len(param) > 0 {
			count += 1
		}
	}
	return count
}

func (vas *VendorAddressService) GetData(ctx context.Context, params *vapb.GetDataParams) (*vapb.GetDataResponse, error) {
	response := vapb.GetDataResponse{Status: utils.Failure}
	userIds := params.GetUserIds()
	sellerIds := params.GetSellerIds()
	ids := params.GetIds()

	if count := ParamCount(userIds, sellerIds, ids); count == 0 {
		response.Message = "param not specified"
		return &response, nil
	} else if count > 1 {
		response.Message = fmt.Sprint("specify any one param")
		return &response, nil
	}
	query := database.DBAPM(ctx).Model(&models.VendorAddress{})
	if len(ids) != 0 {
		query = query.Where("id in (?)", ids)
	} else if len(sellerIds) != 0 {
		query = query.Where("seller_id in (?)", sellerIds)
	} else if len(userIds) != 0 {
		sellers := []*models.Seller{}
		err := database.DBAPM(ctx).Model(&models.Seller{}).Where("user_id in (?)", userIds).Scan(&sellers).Error
		if err != nil {
			response.Message = fmt.Sprint("error in vendor address service GetData API: ", err.Error())
			logger.FromContext(ctx).Error(response.Message)
			return &response, nil
		}
		if len(sellers) == 0 {
			response.Message = "seller not found with the user id"
			return &response, nil
		}
		var sellerIDs []int
		for _, seller := range sellers {
			sellerIDs = append(sellerIDs, int(seller.ID))
		}
		query = query.Where("seller_id in (?)", sellerIDs)
	}
	vendorAddress := []*vapb.VendorAddressObject{}
	err := query.Scan(&vendorAddress).Error
	if err != nil {
		response.Message = fmt.Sprint("error in vendor address service GetData API: ", err.Error())
		logger.FromContext(ctx).Error(response.Message)
		return &response, nil
	}
	response.Status = utils.Success
	if len(vendorAddress) == 0 {
		response.Message = "vendor address not found"
		return &response, nil
	}
	response.VendorAddress = vendorAddress
	response.Message = "fetched vendor address successfully"
	return &response, nil
}

func (vapb *VendorAddressService) VerifyAddress(ctx context.Context, params *vapb.VerifyAddressParams) (*vapb.BasicApiResponse, error) {
	return nil, nil
}
