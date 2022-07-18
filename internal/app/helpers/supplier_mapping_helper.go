package helpers

import (
	"context"
	"fmt"
	"log"
	"time"

	supplierpb "github.com/voonik/goConnect/api/go/ss2/supplier"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/models"
)

func UpdateSupplierOpcMapping(ctx context.Context, id, opcId uint64, delete bool) *supplierpb.BasicApiResponse {
	if err := IsOpcListValid(ctx, []uint64{opcId}); delete == false && err != nil {
		return &supplierpb.BasicApiResponse{Message: err.Error()}
	}

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
	responseId := GetParentCategories(ctx, newIds)
	log.Printf("CMT Response:: categories = %v, response = %v\n", newIds, responseId)

	for _, id := range responseId {
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
	if mapToRestore != nil {
		database.DBAPM(ctx).Model(&models.SupplierCategoryMapping{}).Unscoped().Where("id IN (?)", mapToRestore).Update("deleted_at", nil)
	}
	if mapToDelete != nil {
		database.DBAPM(ctx).Model(&models.SupplierCategoryMapping{}).Unscoped().Where("id IN (?)", mapToDelete).Update("deleted_at", &currentTime)
	}

	newIds = []uint64{}
	for k, v := range categoryToCreateMap {
		if v {
			newIds = append(newIds, k)
		}
	}

	return PrepareCategoreMapping(newIds)
}
