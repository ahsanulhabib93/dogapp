package helpers

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	supplierpb "github.com/voonik/goConnect/api/go/ss2/supplier"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/utils"
)

func GetPaymentAccountDetails(ctx context.Context, supplier models.Supplier, warehouseID uint64) []*supplierpb.PaymentAccountDetailObject {
	type dbResponse struct {
		*supplierpb.PaymentAccountDetailObject
		DhCode string
	}
	paymentDetails := []*dbResponse{}
	query := database.DBAPM(ctx).Model(&models.PaymentAccountDetail{}).
		Joins(models.GetBankJoinStr()).Where("supplier_id = ?", supplier.ID)
	selectQuery := "payment_account_details.*, banks.name bank_name"
	if warehouseID != 0 {
		query = query.Joins(models.JoinPaymentAccountDetailWarehouseMappings()).Where("warehouse_id = ?", warehouseID)
		selectQuery = "payment_account_details.*, banks.name bank_name, payment_account_detail_warehouse_mappings.dh_code dh_code"
	}
	query.Select(selectQuery).Scan(&paymentDetails)

	var paymentDetailIds []uint64
	for _, paymentDetail := range paymentDetails {
		paymentDetailIds = append(paymentDetailIds, paymentDetail.Id)
	}

	paymentResponse := []*supplierpb.PaymentAccountDetailObject{}

	warehouses := GetWarehousesForPaymentAccountDetails(ctx, paymentDetailIds)
	for _, paymentDetail := range paymentDetails {
		resp := paymentDetail.PaymentAccountDetailObject
		resp.Warehouses = warehouses[paymentDetail.Id]
		dhCodes := strings.Split(paymentDetail.DhCode, ",")
		for _, code := range dhCodes {
			dhCode, _ := strconv.Atoi(code)
			resp.DhCode = append(resp.DhCode, uint64(dhCode))
		}
		paymentResponse = append(paymentResponse, resp)
	}

	return paymentResponse
}

func GetWarehousesForPaymentAccountDetails(ctx context.Context, paymentDetailIds []uint64) map[uint64][]uint64 {
	warehouses := make(map[uint64][]uint64)
	var paymentAccountDetailWarehouseMappings []*models.PaymentAccountDetailWarehouseMapping
	database.DBAPM(ctx).Model(&models.PaymentAccountDetailWarehouseMapping{}).Where(
		"payment_account_detail_id IN (?)", paymentDetailIds,
	).Find(&paymentAccountDetailWarehouseMappings)
	for _, paymentDetailWarehouseMapping := range paymentAccountDetailWarehouseMappings {
		paymentAccountDetailID := paymentDetailWarehouseMapping.PaymentAccountDetailID
		warehouses[paymentAccountDetailID] = append(warehouses[paymentAccountDetailID], paymentDetailWarehouseMapping.WarehouseID)
	}
	return warehouses
}

func UpdatePaymentAccountDetailWarehouseMapping(ctx context.Context, paymentAccountDetailId uint64, warehouseIds []uint64) error {
	paymentAccountDetail := models.PaymentAccountDetail{}
	result := database.DBAPM(ctx).Model(&models.PaymentAccountDetail{}).First(&paymentAccountDetail, "id = ?", paymentAccountDetailId)
	if result.RecordNotFound() {
		return fmt.Errorf("PaymentAccountDetail Not Found")
	} else if err := result.Error; err != nil {
		return err
	}

	// get existingMappings :: map[warehouseId]PaymentAccountDetailWarehouseMapping
	paymentAccountDetailWarehouseMappings := []*models.PaymentAccountDetailWarehouseMapping{}
	database.DBAPM(ctx).Model(&paymentAccountDetail).Association("PaymentAccountDetailWarehouseMappings").Find(&paymentAccountDetailWarehouseMappings)
	existingMappings := make(map[uint64]*models.PaymentAccountDetailWarehouseMapping)
	var existingWarehouseIds []uint64
	for _, pADWhMapping := range paymentAccountDetailWarehouseMappings {
		existingMappings[pADWhMapping.WarehouseID] = pADWhMapping
		existingWarehouseIds = append(existingWarehouseIds, pADWhMapping.WarehouseID)
	}

	// Deleting mappings for warehouse_ids not given
	if warehousesToDelete, err := utils.SliceDifference(existingWarehouseIds, warehouseIds); err == nil && warehousesToDelete != nil {
		for _, warehouseId := range warehousesToDelete.([]uint64) {
			database.DBAPM(ctx).Delete(&models.PaymentAccountDetailWarehouseMapping{}, existingMappings[warehouseId])
		}
	}

	// insert new warehouse_ids
	if warehousesToMap, err := utils.SliceDifference(warehouseIds, existingWarehouseIds); err == nil && warehousesToMap != nil {
		for _, warehouseId := range warehousesToMap.([]uint64) {
			database.DBAPM(ctx).Model(&paymentAccountDetail).Association("PaymentAccountDetailWarehouseMappings").Append(&models.PaymentAccountDetailWarehouseMapping{
				WarehouseID: warehouseId,
			})
		}
	}

	return nil
}
