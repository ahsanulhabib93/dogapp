package services

import (
	"context"

	spdpb "github.com/voonik/goConnect/api/go/ss2/seller_pricing_detail"
)

type SellerPricingDetailService struct{}

func (spdpb *SellerPricingDetailService) VerifyDetail(ctx context.Context, params *spdpb.VerifyDetailParams) (*spdpb.BasicApiResponse, error) {
	return &spdpb.BasicApiResponse{}, nil
}
