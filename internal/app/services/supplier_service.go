package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	supplierpb "github.com/voonik/goConnect/api/go/ss2/supplier"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/utils"
)

// SupplierService ...
type SupplierService struct{}

// List ...
func (ss *SupplierService) List(ctx context.Context, params *supplierpb.ListParams) (*supplierpb.ListResponse, error) {
	log.Printf("ListSupplierParams: %+v", params)
	suppliers := []models.Supplier{}
	database.DBAPM(ctx).Model(&models.Supplier{}).Preload("SupplierCategoryMappings").Find(&suppliers)
	resp := ss.prepareResponse(suppliers)
	log.Printf("ListSupplierResponse: %+v", resp)
	return &resp, nil
}

// ListWithSupplierAddresses ...
func (ss *SupplierService) ListWithSupplierAddresses(ctx context.Context, params *supplierpb.ListParams) (*supplierpb.ListResponse, error) {
	log.Printf("ListwithAddressParams: %+v", params)
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
	log.Printf("ListwithAddressResponse: %+v", resp)
	return &resp, nil
}

// Add ...
func (ss *SupplierService) Add(ctx context.Context, params *supplierpb.SupplierParam) (*supplierpb.BasicApiResponse, error) {
	log.Printf("AddSupplierParams: %+v", params)
	resp := supplierpb.BasicApiResponse{Success: false}
	supplier := models.Supplier{
		Name:                     params.GetName(),
		Email:                    params.GetEmail(),
		SupplierType:             utils.SupplierType(params.GetSupplierType()),
		SupplierCategoryMappings: ss.prepareCategoreMapping(params.GetCategoryIds()),
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
	err := database.DBAPM(ctx).Save(&supplier)

	if err != nil && err.Error != nil {
		resp.Message = fmt.Sprintf("Error while creating Supplier: %s", err.Error)
	} else {
		resp.Message = "Supplier Added Successfully"
		resp.Success = true
	}
	log.Printf("AddSupplierResponse: %+v", resp)
	return &resp, nil
}

// Edit ...
func (ss *SupplierService) Edit(ctx context.Context, params *supplierpb.SupplierObject) (*supplierpb.BasicApiResponse, error) {
	log.Printf("EditSupplierParams: %+v", params)
	resp := supplierpb.BasicApiResponse{Success: false}

	supplier := models.Supplier{}
	query := database.DBAPM(ctx).Model(&models.Supplier{})
	if params.GetCategoryIds() != nil {
		query = query.Preload("SupplierCategoryMappings")
	}

	result := query.First(&supplier, params.GetId())
	if result.RecordNotFound() {
		resp.Message = "Supplier Not Found"
	} else {
		err := database.DBAPM(ctx).Model(&supplier).Updates(models.Supplier{
			Name:                     params.GetName(),
			Email:                    params.GetEmail(),
			SupplierType:             utils.SupplierType(params.GetSupplierType()),
			SupplierCategoryMappings: ss.updateCategoryMapping(ctx, supplier.ID, params.GetCategoryIds()),
		})
		if err != nil && err.Error != nil {
			resp.Message = fmt.Sprintf("Error while updating Supplier: %s", err.Error)
		} else {
			resp.Message = "Supplier Edited Successfully"
			resp.Success = true
		}
	}
	log.Printf("EditSupplierResponse: %+v", resp)
	return &resp, nil
}

func (ss *SupplierService) updateCategoryMapping(ctx context.Context, supplierId uint64, newIds []uint64) []models.SupplierCategoryMapping {
	if len(newIds) == 0 {
		return nil
	}

	supplierCategoryMappings := []models.SupplierCategoryMapping{}
	database.DBAPM(ctx).Model(&models.SupplierCategoryMapping{}).Unscoped().Where("supplier_category_mappings.supplier_id = ?", supplierId).Find(&supplierCategoryMappings)
	categoryToCreateMap := map[uint64]bool{}
	for _, id := range newIds {
		categoryToCreateMap[id] = true
	}

	for _, cMap := range supplierCategoryMappings {
		query := database.DBAPM(ctx).Model(&models.SupplierCategoryMapping{}).Unscoped().Where("id = ?", cMap.ID)
		if _, ok := categoryToCreateMap[cMap.ID]; ok {
			query.Update("deleted_at", nil)
			categoryToCreateMap[cMap.ID] = false
		} else if cMap.DeletedAt == nil {
			t := time.Now()
			query.Update("deleted_at", &t)
		}
	}

	newIds = []uint64{}
	for k, v := range categoryToCreateMap {
		if v {
			newIds = append(newIds, k)
		}
	}

	return ss.prepareCategoreMapping(newIds)
}

func (ss *SupplierService) prepareCategoreMapping(ids []uint64) []models.SupplierCategoryMapping {
	categories := []models.SupplierCategoryMapping{}
	for _, id := range ids {
		categories = append(categories, models.SupplierCategoryMapping{
			CategoryID: id,
		})
	}

	return categories
}

func (ss *SupplierService) prepareResponse(suppliers []models.Supplier) supplierpb.ListResponse {
	data := []supplierpb.SupplierObject{}
	for _, supplier := range suppliers {
		temp, _ := json.Marshal(supplier)
		so := supplierpb.SupplierObject{}
		json.Unmarshal(temp, &so)
		for _, cMap := range supplier.SupplierCategoryMappings {
			so.CategoryIds = append(so.CategoryIds, cMap.CategoryID)
		}

		data = append(data, so)

	}

	resp := supplierpb.ListResponse{}
	temp, _ := json.Marshal(data)
	json.Unmarshal(temp, &resp.Data)
	return resp
}
