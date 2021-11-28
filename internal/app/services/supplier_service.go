package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	supplierpb "github.com/voonik/goConnect/api/go/ss2/supplier"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/utils"
)

type SupplierService struct{}

func (ss *SupplierService) ListSupplier(ctx context.Context, params *supplierpb.ListParams) (*supplierpb.ListResponse, error) {
	resp := supplierpb.ListResponse{}
	database.DBAPM(ctx).Model(&models.Supplier{}).Scan(&resp.Data)
	return &resp, nil
}

func (ss *SupplierService) ListWithSupplierAddresses(ctx context.Context, params *supplierpb.ListParams) (*supplierpb.ListResponse, error) {
	resp := supplierpb.ListResponse{}

	query := database.DBAPM(ctx).Model(&models.Supplier{}).Joins("join supplier_addresses on supplier_addresses.supplier_id=suppliers.id")

	if params.GetId() != 0 {
		query = query.Where("suppliers.id = ?", params.GetId())
	}
	if params.GetName() != "" {
		query = query.Where("suppliers.name like ?", fmt.Sprintf("%s%%", params.GetName()))
	}
	if params.GetEmail() != "" {
		query = query.Where("suppliers.email = ?", params.GetEmail())
	}
	if params.GetPhone() != "" {
		query = query.Where("supplier_addresses.phone = ?", params.GetPhone())
	}
	if params.GetCity() != "" {
		query = query.Where("supplier_addresses.city = ?", params.GetCity())
	}

	suppliersWithAddresses := []models.Supplier{{}}
	query.Select("distinct suppliers.*").Preload("SupplierAddresses").Find(&suppliersWithAddresses)

	temp, _ := json.Marshal(suppliersWithAddresses)
	json.Unmarshal(temp, &resp.Data)

	return &resp, nil
}

func (ss *SupplierService) AddSupplier(ctx context.Context, params *supplierpb.SupplierParam) (*supplierpb.BasicApiResponse, error) {
	resp := supplierpb.BasicApiResponse{Success: false}
	supplier := models.Supplier{
		Name:         params.GetName(),
		Email:        params.GetEmail(),
		SupplierType: utils.SupplierType(params.GetSupplierType()),
		SupplierAddresses: []models.SupplierAddress{{
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
			IsDefault: true,
		}},
	}
	err := database.DBAPM(ctx).Model(&models.Supplier{}).Create(&supplier)

	if err != nil && err.Error != nil {
		errorMsg := fmt.Sprintf("Error while creating Supplier: %s", err.Error)
		log.Println(errorMsg)
		resp.Message = errorMsg
	} else {
		resp.Message = "Supplier Added Successfully"
		resp.Success = true
	}
	return &resp, nil
}

func (ss *SupplierService) ListSupplierAddresses(ctx context.Context, params *supplierpb.ListSupplierAddressParams) (*supplierpb.ListSupplierAddressResponse, error) {
	resp := supplierpb.ListSupplierAddressResponse{}
	database.DBAPM(ctx).Model(&models.SupplierAddress{}).Where("supplier_id = ?", params.GetSupplierId()).Scan(&resp.Data)
	return &resp, nil
}

func (ss *SupplierService) AddSupplierAddress(ctx context.Context, params *supplierpb.SupplierAddressParam) (*supplierpb.BasicApiResponse, error) {
	resp := supplierpb.BasicApiResponse{Success: false}

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
			resp.Message = "SupplierAddress Added Successfully"
			resp.Success = true
		}
	}
	return &resp, nil
}
