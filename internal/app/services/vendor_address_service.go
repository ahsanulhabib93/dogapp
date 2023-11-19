package services

import (
	"context"

	vapb "github.com/voonik/goConnect/api/go/ss2/vendor_address"
)

type VendorAddressService struct{}

func (vapb *VendorAddressService) GetData(ctx context.Context, params *vapb.GetDataParams) (*vapb.GetDataResponse, error) {
	return nil, nil
}

func (vapb *VendorAddressService) VerifyAddress(ctx context.Context, params *vapb.VerifyAddressParams) (*vapb.BasicApiResponse, error) {
	return nil, nil
}
