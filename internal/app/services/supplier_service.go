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

func (ss *SupplierService) Get(ctx context.Context, params *supplierpb.GetSupplierParam) (*supplierpb.SupplierObject, error) {
	log.Printf("GetSupplierParam: %+v", params)
	supplier := models.Supplier{}
	result := database.DBAPM(ctx).Model(&models.Supplier{}).Preload("SupplierAddresses").
		Preload("PaymentAccountDetails").First(&supplier, params.GetId())
	if result.RecordNotFound() {
		return nil, NotFound
	}

	resp := supplierDBResponse{}
	database.DBAPM(ctx).Model(&models.Supplier{}).Where("suppliers.id = ?", supplier.ID).
		Joins(" left join supplier_category_mappings on supplier_category_mappings.supplier_id=suppliers.id").
		Joins(" left join supplier_sa_mappings on supplier_sa_mappings.supplier_id=suppliers.id").Group("id").
		Select(ss.getResponseField()).Scan(&resp)

	resp.SupplierAddresses = supplier.SupplierAddresses
	resp.PaymentAccountDetails = supplier.PaymentAccountDetails
	log.Printf("GetSupplierResponse: %+v", resp)
	return ss.prepareSupplierResponse(resp), nil
}

// List ...
func (ss *SupplierService) List(ctx context.Context, params *supplierpb.ListParams) (*supplierpb.ListResponse, error) {
	log.Printf("ListSupplierParams: %+v", params)
	suppliers := []supplierDBResponse{}
	database.DBAPM(ctx).Model(&models.Supplier{}).
		Joins(" left join supplier_category_mappings on supplier_category_mappings.supplier_id=suppliers.id").Group("id").
		Joins(" left join supplier_sa_mappings on supplier_sa_mappings.supplier_id=suppliers.id").Group("id").
		Select(ss.getResponseField()).Scan(&suppliers)
	resp := ss.prepareListResponse(suppliers)
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
	if len(params.GetSupplierIds()) != 0 {
		query = query.Where("suppliers.id IN (?)", params.GetSupplierIds())
	}
	if params.GetName() != "" {
		query = query.Where("suppliers.name LIKE ?", fmt.Sprintf("%s%%", params.GetName()))
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
		SupplierSaMappings:       ss.prepareSaMapping(params.GetSaIds()),
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
		resp.Id = supplier.ID
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
	if params.GetSaIds() != nil {
		query = query.Preload("SupplierSaMappings")
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
			SupplierSaMappings:       ss.updateSaMapping(ctx, supplier.ID, params.GetSaIds()),
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

func (ss *SupplierService) updateSaMapping(ctx context.Context, supplierId uint64, newIds []uint64) []models.SupplierSaMapping {
	if len(newIds) == 0 {
		return nil
	}

	sourcingAssociateMappings := []models.SupplierSaMapping{}
	database.DBAPM(ctx).Model(&models.SupplierSaMapping{}).Unscoped().Where("supplier_id = ?", supplierId).Find(&sourcingAssociateMappings)
	saCreateMap := map[uint64]bool{}
	for _, id := range newIds {
		saCreateMap[id] = true
	}

	mapToDelete := []uint64{}
	mapToRestore := []uint64{}
	for _, sMap := range sourcingAssociateMappings {
		_, inNewList := saCreateMap[sMap.SourcingAssociateId]
		if !inNewList {
			mapToDelete = append(mapToDelete, sMap.ID)
		} else {
			saCreateMap[sMap.SourcingAssociateId] = false
			mapToRestore = append(mapToRestore, sMap.ID)
		}
	}

	currentTime := time.Now()
	database.DBAPM(ctx).Model(&models.SupplierSaMapping{}).Unscoped().Where("id IN (?)", mapToRestore).Update("deleted_at", nil)
	database.DBAPM(ctx).Model(&models.SupplierSaMapping{}).Unscoped().Where("id IN (?)", mapToDelete).Update("deleted_at", &currentTime)
	newIds = []uint64{}
	for k, v := range saCreateMap {
		if v {
			newIds = append(newIds, k)
		}
	}

	return ss.prepareSaMapping(newIds)
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
func (ss *SupplierService) prepareSaMapping(ids []uint64) []models.SupplierSaMapping {
	sourcing_associates := []models.SupplierSaMapping{}
	for _, id := range ids {
		sourcing_associates = append(sourcing_associates, models.SupplierSaMapping{
			SourcingAssociateId: id,
		})
	}
	return sourcing_associates
}

func (ss *SupplierService) getResponseField() string {
	s := []string{
		"suppliers.id",
		"suppliers.status",
		"suppliers.supplier_type",
		"suppliers.name",
		"suppliers.email",
		"GROUP_CONCAT( DISTINCT supplier_category_mappings.category_id) as category_ids",
		"GROUP_CONCAT( DISTINCT supplier_sa_mappings.sourcing_associate_id) as sa_ids",
	}

	return strings.Join(s, ",")
}

func (ss *SupplierService) prepareListResponse(suppliers []supplierDBResponse) supplierpb.ListResponse {
	data := []*supplierpb.SupplierObject{}
	for _, supplier := range suppliers {
		data = append(data, ss.prepareSupplierResponse(supplier))
	}

	return supplierpb.ListResponse{Data: data}
}

func (ss *SupplierService) prepareSupplierResponse(supplier supplierDBResponse) *supplierpb.SupplierObject {
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
	so.SaIds = []uint64{}
	for _, saId := range strings.Split(supplier.SaIds, ",") {
		saId = strings.TrimSpace(saId)
		if saId == "" {
			continue
		}
		v, _ := strconv.Atoi(saId)
		so.SaIds = append(so.SaIds, uint64(v))
	}

	return so
}

type supplierDBResponse struct {
	models.Supplier
	CategoryIds string `json:"category_ids,omitempty"`
	SaIds       string `json:"sa_ids,omitempty"`
}
