package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	supplierpb "github.com/voonik/goConnect/api/go/ss2/supplier"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/utils"
)

// SupplierService ...
type SupplierService struct{}

func (ss *SupplierService) Get(ctx context.Context, params *supplierpb.GetSupplierParam) (*supplierpb.SupplierResponse, error) {
	log.Printf("GetSupplierParam: %+v", params)
	supplier := models.Supplier{}
	result := database.DBAPM(ctx).Model(&models.Supplier{}).Preload("SupplierAddresses").
		Preload("PaymentAccountDetails").First(&supplier, params.GetId())
	if result.RecordNotFound() {
		return &supplierpb.SupplierResponse{Success: false}, nil
	}

	resp := supplierDBResponse{}
	database.DBAPM(ctx).Model(&models.Supplier{}).Where("suppliers.id = ?", supplier.ID).
		Joins(" left join supplier_category_mappings on supplier_category_mappings.supplier_id=suppliers.id").
		Joins(" left join supplier_opc_mappings on supplier_opc_mappings.supplier_id=suppliers.id").Group("id").
		Select(ss.getResponseField()).Scan(&resp)

	resp.SupplierAddresses = supplier.SupplierAddresses
	resp.PaymentAccountDetails = supplier.PaymentAccountDetails
	log.Printf("GetSupplierResponse: %+v", resp)
	return &supplierpb.SupplierResponse{
		Success: true,
		Data:    ss.prepareSupplierResponse(resp)}, nil
}

// List ...
func (ss *SupplierService) List(ctx context.Context, params *supplierpb.ListParams) (*supplierpb.ListResponse, error) {
	log.Printf("ListSupplierParams: %+v", params)
	suppliers := []supplierDBResponse{}
	query := database.DBAPM(ctx).Model(&models.Supplier{})
	query = ss.prepareFilter(query, params)
	var total uint64
	query.Count(&total)

	ss.setPage(query, params)
	query.Joins(" left join supplier_category_mappings on supplier_category_mappings.supplier_id=suppliers.id").Group("id").
		Joins(" left join supplier_opc_mappings on supplier_opc_mappings.supplier_id=suppliers.id").Group("id").
		Select(ss.getResponseField()).Scan(&suppliers)

	resp := ss.prepareListResponse(suppliers, total)
	log.Printf("ListSupplierResponse: %+v", resp)
	return &resp, nil
}

// ListWithSupplierAddresses ...
func (ss *SupplierService) ListWithSupplierAddresses(ctx context.Context, params *supplierpb.ListParams) (*supplierpb.ListResponse, error) {
	log.Printf("ListwithAddressParams: %+v", params)
	resp := supplierpb.ListResponse{}

	query := database.DBAPM(ctx).Model(&models.Supplier{})
	query = ss.prepareFilter(query, params)

	var total uint64
	query.Count(&total)
	ss.setPage(query, params)
	suppliersWithAddresses := []models.Supplier{{}}
	query.Joins("join supplier_addresses on supplier_addresses.supplier_id=suppliers.id").
		Select("distinct suppliers.*").Preload("SupplierAddresses").Find(&suppliersWithAddresses)

	temp, _ := json.Marshal(suppliersWithAddresses)
	json.Unmarshal(temp, &resp.Data)
	resp.TotalCount = total
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
		SupplierOpcMappings:      ss.prepareOpcMapping(params.GetOpcIds()),
		SupplierAddresses:        ss.prepareSupplierAddress(params),
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

func (ss *SupplierService) Map(ctx context.Context, params *supplierpb.SupplierMappingParams) (*supplierpb.BasicApiResponse, error) {
	supplier := &models.Supplier{}
	result := database.DBAPM(ctx).Model(&models.Supplier{}).First(supplier, params.GetSupplierId())
	if result.RecordNotFound() {
		return &supplierpb.BasicApiResponse{Message: "Supplier Not Found"}, nil
	}

	isDeleting := false
	if strings.ToLower(params.OperationType) == "delete" {
		isDeleting = true
	}

	switch strings.ToLower(params.MapWith) {
	case "opc":
		return ss.opcMapping(ctx, params.SupplierId, params.Id, isDeleting), nil
	}

	return &supplierpb.BasicApiResponse{Message: "Invalid mapping pair"}, nil
}

func (ss *SupplierService) opcMapping(ctx context.Context, id, opcId uint64, delete bool) *supplierpb.BasicApiResponse {
	resp := &supplierpb.BasicApiResponse{Success: true}

	opcMap := &models.SupplierOpcMapping{}
	result := database.DBAPM(ctx).Model(&opcMap).Unscoped().First(opcMap, "processing_center_id = ? AND supplier_id = ?", opcId, id)
	if notFound := result.RecordNotFound(); (notFound && delete) ||
		(!notFound && opcMap.DeletedAt != nil && delete) {
		return resp
	}

	opcMap.DeletedAt = nil
	opcMap.SupplierID = id
	opcMap.ProcessingCenterID = opcId
	if now := time.Now(); delete {
		opcMap.DeletedAt = &now
	}

	if err := database.DBAPM(ctx).Unscoped().Save(&opcMap).Error; err != nil {
		resp.Success = false
		resp.Message = fmt.Sprintf("Error while processing Supplier-OPC mapping: %s", err.Error())
	}

	return resp
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

func (ss *SupplierService) prepareFilter(query *gorm.DB, params *supplierpb.ListParams) *gorm.DB {
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
	if status := params.GetStatus(); status == models.SupplierStatusActive || status == models.SupplierStatusPending {
		query = query.Where("suppliers.status = ?", params.GetStatus())
	}
	if params.GetPhone() != "" {
		query = query.Where("supplier_addresses.phone = ?", params.GetPhone())
	}
	if params.GetCity() != "" {
		query = query.Where("supplier_addresses.city = ?", params.GetCity())
	}

	return query
}

func (ss *SupplierService) setPage(query *gorm.DB, params *supplierpb.ListParams) {
	if params.GetPerPage() <= 0 || params.GetPerPage() > utils.DEFAULT_PER_PAGE {
		params.PerPage = utils.DEFAULT_PER_PAGE
	}

	offset := params.GetPage() * params.GetPerPage()
	*query = *query.Offset(offset).Limit(params.GetPerPage())
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
func (ss *SupplierService) prepareOpcMapping(ids []uint64) []models.SupplierOpcMapping {
	processCenters := []models.SupplierOpcMapping{}
	for _, id := range ids {
		processCenters = append(processCenters, models.SupplierOpcMapping{
			ProcessingCenterID: id,
		})
	}
	return processCenters
}

func (ss *SupplierService) getResponseField() string {
	s := []string{
		"suppliers.id",
		"suppliers.status",
		"suppliers.supplier_type",
		"suppliers.name",
		"suppliers.email",
		"GROUP_CONCAT( DISTINCT supplier_category_mappings.category_id) as category_ids",
		"GROUP_CONCAT( DISTINCT supplier_opc_mappings.processing_center_id) as opc_ids",
	}

	return strings.Join(s, ",")
}

func (ss *SupplierService) prepareListResponse(suppliers []supplierDBResponse, total uint64) supplierpb.ListResponse {
	data := []*supplierpb.SupplierObject{}
	for _, supplier := range suppliers {
		data = append(data, ss.prepareSupplierResponse(supplier))
	}

	return supplierpb.ListResponse{Data: data, TotalCount: total}
}

func (ss *SupplierService) prepareSupplierResponse(supplier supplierDBResponse) *supplierpb.SupplierObject {
	temp, _ := json.Marshal(supplier)
	so := &supplierpb.SupplierObject{}
	json.Unmarshal(temp, so)

	so.CategoryIds = []uint64{}
	for _, cId := range strings.Split(supplier.CategoryIds, ",") {
		if cId := strings.TrimSpace(cId); cId != "" {
			v, _ := strconv.Atoi(cId)
			so.CategoryIds = append(so.CategoryIds, uint64(v))
		}
	}

	so.OpcIds = []uint64{}
	for _, saId := range strings.Split(supplier.OpcIds, ",") {
		if opcId := strings.TrimSpace(saId); opcId != "" {
			v, _ := strconv.Atoi(saId)
			so.OpcIds = append(so.OpcIds, uint64(v))
		}
	}

	return so
}

func (ss *SupplierService) prepareSupplierAddress(params *supplierpb.SupplierParam) []models.SupplierAddress {
	if params.GetFirstname() == "" && params.GetLastname() == "" && params.GetAddress1() == "" && params.GetAddress2() == "" &&
		params.GetLandmark() == "" && params.GetCity() == "" && params.GetState() == "" && params.GetCountry() == "" &&
		params.GetZipcode() == "" && params.GetPhone() == "" && params.GetGstNumber() == "" {
		return []models.SupplierAddress{}
	}

	return []models.SupplierAddress{{
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
	}}
}

type supplierDBResponse struct {
	models.Supplier
	CategoryIds string `json:"category_ids,omitempty"`
	OpcIds      string `json:"opc_ids,omitempty"`
}
