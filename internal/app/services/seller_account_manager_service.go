package services

import (
	"context"

	sampb "github.com/voonik/goConnect/api/go/ss2/seller_account_manager"
	"github.com/voonik/ss2/internal/app/helpers"
	"github.com/voonik/ss2/internal/app/utils"
)

type SellerAccountManagerService struct{}

func (sams *SellerAccountManagerService) List(ctx context.Context, params *sampb.ListParams) (*sampb.ListResponse, error) {
	resp := &sampb.ListResponse{}

	if params.SellerId == utils.Zero {
		return resp, nil
	}

	var samList []models.SellerAccountManager
	var accountManagers []*sampb.AccountManagerObject
	database.DBAPM(ctx).Model(&models.SellerAccountManager{}).Where(`seller_id =?`, sellerID).Order("role, priority").Scan(&samList)

	for _, sam := range samList {
		accountManagers = append(accountManagers, &sampb.AccountManagerObject{
			Id:       sam.ID,
			Email:    sam.Email,
			Phone:    uint64(sam.Phone),
			Name:     sam.Name,
			Priority: uint64(sam.Priority),
			Role:     sam.Role,
		})
	}

	resp.Status = "success"
	resp.AccountManager = accountManagers
	return resp, nil
}
