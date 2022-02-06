package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
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
	suppliers := []supplierDBResponse{}
	database.DBAPM(ctx).Model(&models.Supplier{}).Joins("left join supplier_category_mappings on supplier_category_mappings.supplier_id=suppliers.id").Group("id").Select(ss.getResponseField()).Scan(&suppliers)
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
		Status:                   params.GetStatus(),
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
			Status:                   params.GetStatus(),
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
	database.DBAPM(ctx).Model(&models.SupplierCategoryMapping{}).Unscoped().Where("supplier_id = ?", supplierId).Find(&supplierCategoryMappings)
	categoryToCreateMap := map[uint64]bool{}
	for _, id := range newIds {
		categoryToCreateMap[id] = true
	}

	mapToDelete := []uint64{}
	mapToRestore := []uint64{}
	for _, cMap := range supplierCategoryMappings {
		_, inNewList := categoryToCreateMap[cMap.CategoryID]
		if !inNewList {
			mapToDelete = append(mapToDelete, cMap.ID)
		} else {
			categoryToCreateMap[cMap.CategoryID] = false
			mapToRestore = append(mapToRestore, cMap.ID)
		}
	}

	currentTime := time.Now()
	database.DBAPM(ctx).Model(&models.SupplierCategoryMapping{}).Unscoped().Where("id IN (?)", mapToRestore).Update("deleted_at", nil)
	database.DBAPM(ctx).Model(&models.SupplierCategoryMapping{}).Unscoped().Where("id IN (?)", mapToDelete).Update("deleted_at", &currentTime)
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

func (ss *SupplierService) getResponseField() string {
	s := []string{
		"suppliers.id",
		"suppliers.supplier_type",
		"suppliers.name",
		"suppliers.email",
		"GROUP_CONCAT(supplier_category_mappings.category_id) as category_ids",
	}

	return strings.Join(s, ",")
}

func (ss *SupplierService) prepareResponse(suppliers []supplierDBResponse) supplierpb.ListResponse {
	data := []*supplierpb.SupplierObject{}
	for _, supplier := range suppliers {
		temp, _ := json.Marshal(supplier)
		so := &supplierpb.SupplierObject{}
		json.Unmarshal(temp, so)
		so.CategoryIds = []uint64{}
		for _, cId := range strings.Split(supplier.CategoryIds, ",") {
			cId = strings.TrimSpace(cId)
			if cId == "" {
				continue
			}

			v, _ := strconv.Atoi(cId)
			so.CategoryIds = append(so.CategoryIds, uint64(v))
		}

		data = append(data, so)
	}

	return supplierpb.ListResponse{Data: data}
}

type supplierDBResponse struct {
	models.Supplier
	CategoryIds string `json:"category_ids,omitempty"`
}
