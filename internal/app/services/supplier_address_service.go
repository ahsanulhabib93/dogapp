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
	log.Printf("ListAddressParams: %+v", params)
	resp := addresspb.ListSupplierAddressResponse{}
	query := database.DBAPM(ctx).Model(&models.SupplierAddress{})
	if params.GetSupplierId() != 0 {
		query = query.Where("supplier_id = ?", params.GetSupplierId())
	}
	if params.GetId() != 0 {
		query = query.Where("id = ?", params.GetId())
	}
	query.Scan(&resp)
	log.Printf("ListAddressResponse: %+v", resp)
	return &resp, nil
}

// Add ...
func (sas *SupplierAddressService) Add(ctx context.Context, params *addresspb.SupplierAddressParam) (*addresspb.BasicApiResponse, error) {
	log.Printf("AddAddressParams: %+v", params)
	resp := addresspb.BasicApiResponse{Success: false}

	supplier := &models.Supplier{}
	result := database.DBAPM(ctx).Model(&models.Supplier{}).First(supplier, params.GetSupplierId())
	if result.RecordNotFound() {
		resp.Message = "Supplier Not Found"
	} else if !supplier.IsChangeAllowed(ctx) {
		resp.Message = "Change Not Allowed"
	} else {
		supplierAddress := models.SupplierAddress{
			SupplierID: supplier.ID,
			Firstname:  params.GetFirstname(),
			Lastname:   params.GetLastname(),
			Address1:   params.GetAddress1(),
			Address2:   params.GetAddress2(),
			Landmark:   params.GetLandmark(),
			City:       params.GetCity(),
			State:      params.GetState(),
			Country:    params.GetCountry(),
			Zipcode:    params.GetZipcode(),
			Phone:      params.GetPhone(),
			GstNumber:  params.GetGstNumber(),
			IsDefault:  params.GetIsDefault(),
		}

		err := database.DBAPM(ctx).Save(&supplierAddress)
		if err != nil && err.Error != nil {
			resp.Message = fmt.Sprintf("Error while creating Supplier Address: %s", err.Error)
		} else {
			helpers.UpdateDefaultAddress(ctx, &supplierAddress)
			resp.Message = "Supplier Address Added Successfully"
			resp.Success = true
		}
	}
	log.Printf("AddAddressResponse: %+v", resp)
	return &resp, nil
}

// Edit ...
func (sas *SupplierAddressService) Edit(ctx context.Context, params *addresspb.SupplierAddressObject) (*addresspb.BasicApiResponse, error) {
	log.Printf("EditAddressParams: %+v", params)
	resp := addresspb.BasicApiResponse{Success: false}

	supplierAddress := models.SupplierAddress{}
	result := database.DBAPM(ctx).Model(&models.SupplierAddress{}).First(&supplierAddress, params.GetId())
	if result.RecordNotFound() {
		resp.Message = "Supplier Address Not Found"
	} else {
		supplier := models.Supplier{}
		database.DBAPM(ctx).Model(&models.Supplier{}).First(&supplier, supplierAddress.SupplierID)
		if !supplier.IsChangeAllowed(ctx) {
			resp.Message = "Change Not Allowed"
		} else if supplierAddress.IsDefault && !params.GetIsDefault() {
			resp.Message = "Default address is required"
		} else {
			err := database.DBAPM(ctx).Model(&supplierAddress).Updates(models.SupplierAddress{
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
			})

			if err != nil && err.Error != nil {
				resp.Message = fmt.Sprintf("Error while updating Supplier Address: %s", err.Error)
			} else {
				helpers.UpdateDefaultAddress(ctx, &supplierAddress)
				resp.Message = "Supplier Address Edited Successfully"
				resp.Success = true
			}
		}
	}
	log.Printf("EditAddressResponse: %+v", resp)
	return &resp, nil
}
