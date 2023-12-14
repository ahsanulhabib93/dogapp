package services

import (
	"context"

	sampb "github.com/voonik/goConnect/api/go/ss2/seller_account_manager"
)

type SellerAccountManagerService struct{}

func (sams *SellerAccountManagerService) List(ctx context.Context, params *sampb.ListParams) (*sampb.ListResponse, error) {
	resp := &sampb.ListResponse{}
	return resp, nil
}
