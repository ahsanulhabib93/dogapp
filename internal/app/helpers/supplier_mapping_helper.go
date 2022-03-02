package helpers

import (
	"context"
	"fmt"
	"log"
	"time"

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

func GetOPCListForCurrentUser(ctx context.Context) []uint64 {
	opcList := []uint64{}

	userId := *utils.GetCurrentUserID(ctx)
	resp, err := getOpcClient().ProcessingCenterList(ctx, userId)
	if err != nil {
		log.Printf("GetOPCListForCurrentUser: Failed to fetch OPC list. Error: %v\n", err)
		return opcList
	}

	for _, opc := range resp.Data {
		opcList = append(opcList, opc.OpcId)
	}

	log.Printf("GetOPCListForCurrentUser: opc list = %v\n", opcList)
	return opcList
}
