package services

import (
	"context"

	spdspb "github.com/voonik/goConnect/api/go/ss2/seller_pricing_detail"
)

type SellerPricingDetailService struct{}

func (spdpb *SellerPricingDetailService) VerifyDetail(ctx context.Context, params *spdspb.VerifyDetailParams) (*spdspb.BasicApiResponse, error) {
	return &spdspb.BasicApiResponse{}, nil
}
