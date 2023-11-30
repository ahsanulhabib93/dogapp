package services

import (
	"context"

	sbdpb "github.com/voonik/goConnect/api/go/ss2/seller_bank_detail"
)

type SellerBankDetailService struct{}

func (sbds *SellerBankDetailService) VerifyBankDetail(ctx context.Context, params *sbdpb.VerifyBankDetailParams) (*sbdpb.BasicApiResponse, error) {
	return &sbdpb.BasicApiResponse{}, nil
}
