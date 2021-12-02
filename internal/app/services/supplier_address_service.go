package services

import (
	"context"
	"fmt"
	"log"

	addresspb "github.com/voonik/goConnect/api/go/ss2/supplier_address"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/helpers"
	"github.com/voonik/ss2/internal/app/models"
)

// SupplierAddressService ...
type SupplierAddressService struct{}

// List ...
func (sas *SupplierAddressService) List(ctx context.Context, params *addresspb.ListSupplierAddressParams) (*addresspb.ListSupplierAddressResponse, error) {
	resp := addresspb.ListSupplierAddressResponse{}
	database.DBAPM(ctx).Model(&models.SupplierAddress{}).Where("supplier_id = ?", params.GetSupplierId()).Scan(&resp.Data)
	return &resp, nil
}

// Add ...
func (sas *SupplierAddressService) Add(ctx context.Context, params *addresspb.SupplierAddressParam) (*addresspb.BasicApiResponse, error) {
	resp := addresspb.BasicApiResponse{Success: false}

	supplier := &models.Supplier{}
	result := database.DBAPM(ctx).Model(&models.Supplier{}).First(supplier, params.GetSupplierId())
	if result.RecordNotFound() {
		resp.Message = "Supplier Not Found"
	} else {
		supplierAddress := models.SupplierAddress{
			Supplier:  *supplier,
			Firstname: params.GetFirstname(),
			Lastname:  params.GetLastname(),
			Address1:  params.GetAddress1(),
			Address2:  params.GetAddress2(),
			Landmark:  params.GetLandmark(),
			City:      params.GetCity(),
			State:     params.GetState(),
			Country:   params.GetCountry(),
			Zipcode:   params.GetZipcode(),
			Phone:     params.GetPhone(),
			GstNumber: params.GetGstNumber(),
			IsDefault: params.GetIsDefault(),
		}

		err := database.DBAPM(ctx).Model(&models.SupplierAddress{}).Create(&supplierAddress)
		if err != nil && err.Error != nil {
			errorMsg := fmt.Sprintf("Error while creating Supplier Address: %s", err.Error)
			log.Println(errorMsg)
			resp.Message = errorMsg
		} else {
			helpers.UpdateOtherAddress(ctx, &supplierAddress)
			resp.Message = "SupplierAddress Added Successfully"
			resp.Success = true
		}
	}
	return &resp, nil
}
