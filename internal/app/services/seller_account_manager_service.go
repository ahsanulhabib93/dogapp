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

	resp.Status = "success"
	resp.AccountManager = helpers.GetAndFormatSellerAccountManager(ctx, params.SellerId)

	return resp, nil
}
