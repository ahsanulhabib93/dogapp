package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	supplierpb "github.com/voonik/goConnect/api/go/ss2/supplier"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/utils"
)

func UpdateSupplierOpcMapping(ctx context.Context, id, opcId uint64, delete bool) *supplierpb.BasicApiResponse {
	resp := &supplierpb.BasicApiResponse{Success: true, Message: "Supplier Mapped with OPC"}

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
		resp.Message = "Supplier Unmapped with OPC"
	}

	if err := database.DBAPM(ctx).Unscoped().Save(&opcMap).Error; err != nil {
		resp.Success = false
		resp.Message = fmt.Sprintf("Error while processing Supplier-OPC mapping: %s", err.Error())
	}

	return resp
}

func UpdateSupplierCategoryMapping(ctx context.Context, supplierId uint64, newIds []uint64) []models.SupplierCategoryMapping {
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

	return PrepareCategoreMapping(newIds)
}

func PrepareFilter(query *gorm.DB, params *supplierpb.ListParams) *gorm.DB {
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

func SetPage(query *gorm.DB, params *supplierpb.ListParams) {
	if params.GetPerPage() <= 0 || params.GetPerPage() > utils.DEFAULT_PER_PAGE {
		params.PerPage = utils.DEFAULT_PER_PAGE
	}

	params.Page = utils.Int64Max(utils.DEFAULT_PAGE, params.GetPage())
	offset := (params.GetPage() - 1) * params.GetPerPage()
	*query = *query.Offset(offset).Limit(params.GetPerPage())
}

func PrepareCategoreMapping(ids []uint64) []models.SupplierCategoryMapping {
	categories := []models.SupplierCategoryMapping{}
	for _, id := range ids {
		categories = append(categories, models.SupplierCategoryMapping{
			CategoryID: id,
		})
	}

	return categories
}

func PrepareOpcMapping(ids []uint64) []models.SupplierOpcMapping {
	processCenters := []models.SupplierOpcMapping{}
	for _, id := range ids {
		processCenters = append(processCenters, models.SupplierOpcMapping{
			ProcessingCenterID: id,
		})
	}
	return processCenters
}

func PrepareListResponse(suppliers []SupplierDBResponse, total uint64) supplierpb.ListResponse {
	data := []*supplierpb.SupplierObject{}
	for _, supplier := range suppliers {
		data = append(data, PrepareSupplierResponse(supplier))
	}

	return supplierpb.ListResponse{Data: data, TotalCount: total}
}

func PrepareSupplierResponse(supplier SupplierDBResponse) *supplierpb.SupplierObject {
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

func PrepareSupplierAddress(params *supplierpb.SupplierParam) []models.SupplierAddress {
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

type SupplierDBResponse struct {
	models.Supplier
	CategoryIds string `json:"category_ids,omitempty"`
	OpcIds      string `json:"opc_ids,omitempty"`
}
