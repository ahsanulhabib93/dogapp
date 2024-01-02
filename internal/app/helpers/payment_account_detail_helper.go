package helpers

import (
	"context"
	"fmt"
	"strings"

	"github.com/shopuptech/go-libs/logger"
	paywellPb "github.com/voonik/goConnect/api/go/paywell_token/payment_gateway"
	paymentpb "github.com/voonik/goConnect/api/go/ss2/payment_account_detail"

	supplierpb "github.com/voonik/goConnect/api/go/ss2/supplier"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/utils"
)

func GetPaymentAccountDetails(ctx context.Context, supplier models.Supplier, warehouseID uint64) []*supplierpb.PaymentAccountDetailObject {
	type dbResponse struct {
		*supplierpb.PaymentAccountDetailObject
		DhCodeStr string `json:"dh_code_str,omitempty"`
	}
	paymentDetails := []*dbResponse{}
	query := database.DBAPM(ctx).Model(&models.PaymentAccountDetail{}).
		Joins(models.GetBankJoinStr()).Where("supplier_id = ?", supplier.ID)
	selectQuery := "payment_account_details.*, banks.name bank_name"
	if warehouseID != 0 {
		query = query.Joins(models.JoinPaymentAccountDetailWarehouseMappings()).Where("warehouse_id = ?", warehouseID)
		selectQuery = "payment_account_details.*, banks.name bank_name, payment_account_detail_warehouse_mappings.dh_code dh_code_str"
	}
	query.Select(selectQuery).Scan(&paymentDetails)
	fmt.Println("logger here paymentDetails ", paymentDetails)

	var paymentDetailIds []uint64
	for _, paymentDetail := range paymentDetails {
		paymentDetailIds = append(paymentDetailIds, paymentDetail.Id)
	}

	paymentResponse := []*supplierpb.PaymentAccountDetailObject{}

	warehouseDhCodeMap := GetWarehouseDhCodeForPaymentAccountDetails(ctx, paymentDetailIds)
	for _, paymentDetail := range paymentDetails {
		resp := paymentDetail.PaymentAccountDetailObject
		warehouses := []uint64{}
		for whId := range warehouseDhCodeMap[paymentDetail.Id] {
			warehouses = append(warehouses, whId)
		}
		fmt.Println("logger here paymentDetail.DhCodeStr ", paymentDetail.DhCodeStr)
		resp.Warehouses = warehouses
		resp.DhCode = []string{}
		if strings.TrimSpace(paymentDetail.DhCodeStr) != utils.EmptyString {
			resp.DhCode = strings.Split(paymentDetail.DhCodeStr, ",")
		}
		resp.WarehouseDhCodeMap = warehouseDhCodeMap[paymentDetail.Id]
		paymentResponse = append(paymentResponse, resp)
	}

	return paymentResponse
}

func GetWarehouseDhCodeForPaymentAccountDetails(ctx context.Context, paymentDetailIds []uint64) map[uint64]map[uint64]*supplierpb.DhCodes {
	warehouseDhCodeMap := make(map[uint64]map[uint64]*supplierpb.DhCodes)
	var paymentAccountDetailWarehouseMappings []*models.PaymentAccountDetailWarehouseMapping
	database.DBAPM(ctx).Model(&models.PaymentAccountDetailWarehouseMapping{}).Where(
		"payment_account_detail_id IN (?)", paymentDetailIds,
	).Find(&paymentAccountDetailWarehouseMappings)
	for _, paymentDetailWarehouseMapping := range paymentAccountDetailWarehouseMappings {
		paymentAccountDetailID := paymentDetailWarehouseMapping.PaymentAccountDetailID
		if warehouseDhCodeMap[paymentAccountDetailID] == nil {
			warehouseDhCodeMap[paymentAccountDetailID] = map[uint64]*supplierpb.DhCodes{}
		}

		warehouseDhCodeMap[paymentAccountDetailID][paymentDetailWarehouseMapping.WarehouseID] = &supplierpb.DhCodes{}
		if strings.TrimSpace(paymentDetailWarehouseMapping.DhCode) != utils.EmptyString {
			warehouseDhCodeMap[paymentAccountDetailID][paymentDetailWarehouseMapping.WarehouseID] = &supplierpb.DhCodes{DhCode: strings.Split(paymentDetailWarehouseMapping.DhCode, ",")}
		}
	}
	return warehouseDhCodeMap
}

func UpdatePaymentAccountDetailWarehouseMapping(ctx context.Context, paymentAccountDetailId uint64, warehouseIds []uint64, warehouseDhCodeMap map[uint64]*paymentpb.DhCodes) error {
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
	warehousesToDelete, err := utils.SliceDifference(existingWarehouseIds, warehouseIds)
	if err == nil && warehousesToDelete != nil {
		for _, warehouseId := range warehousesToDelete.([]uint64) {
			database.DBAPM(ctx).Delete(&models.PaymentAccountDetailWarehouseMapping{}, existingMappings[warehouseId])
		}
	}

	if warehousesToDelete == nil {
		warehousesToDelete = []uint64{}
	}

	warehouseToUpdate, err := utils.SliceDifference(existingWarehouseIds, warehousesToDelete)
	if err == nil && warehouseToUpdate != nil && warehouseDhCodeMap != nil {
		for _, warehouseId := range warehouseToUpdate.([]uint64) {
			if warehouseDhCodeMap[warehouseId] == nil {
				continue
			}

			database.DBAPM(ctx).Model(&models.PaymentAccountDetailWarehouseMapping{}).
				Where("payment_account_detail_id = ? and warehouse_id = ?", paymentAccountDetailId, warehouseId).
				UpdateColumn("dh_code", strings.Join(warehouseDhCodeMap[warehouseId].GetDhCode(), ","))
		}
	}

	// insert new warehouse_ids
	if warehousesToMap, err := utils.SliceDifference(warehouseIds, existingWarehouseIds); err == nil && warehousesToMap != nil {
		for _, warehouseId := range warehousesToMap.([]uint64) {
			payload := &models.PaymentAccountDetailWarehouseMapping{
				WarehouseID: warehouseId,
			}

			if warehouseDhCodeMap[warehouseId].GetDhCode() != nil {
				payload.DhCode = strings.Join(warehouseDhCodeMap[warehouseId].GetDhCode(), ",")
			}

			database.DBAPM(ctx).Model(&paymentAccountDetail).Association("PaymentAccountDetailWarehouseMappings").Append(payload)
		}
	}

	return nil
}

func StoreEncryptCardInfo(ctx context.Context, extraDetails paymentpb.ExtraDetails, paymentAccountDetail *models.PaymentAccountDetail, accountNumber string) (bool, *models.PaymentAccountDetail) {
	uniqueId := utils.CreatePaywellUniqueKey(paymentAccountDetail.ID)
	expiryMonth, expiryYear := utils.FetchMonthAndYear(extraDetails.ExpiryDate)
	logger.FromContext(ctx).Info("Payload for CreatePaywellCard : unique id %v, card info %v, expiry month %v, expiry year %v", uniqueId, paymentAccountDetail.AccountNumber, expiryMonth, expiryYear)
	paywellResponse := getAPIHelperInstance().CreatePaywellCard(ctx, &paywellPb.CreateCardRequest{
		UniqueId:    uniqueId,
		CardInfo:    accountNumber,
		ExpiryMonth: expiryMonth,
		ExpiryYear:  expiryYear,
	})

	if paywellResponse.IsError {
		database.DBAPM(ctx).Delete(paymentAccountDetail)
		return false, nil
	}
	paymentAccountDetail.AccountNumber = paywellResponse.MaskedNumber
	paymentAccountDetail.SetExtraDetails(models.PaymentAccountDetailExtraDetails{
		UniqueId: uniqueId,
		Token:    paywellResponse.GetToken(),
	})
	database.DBAPM(ctx).Save(&paymentAccountDetail)
	return true, paymentAccountDetail
}

func HandleExtraDetailsValidation(ctx context.Context, extraDetails *paymentpb.ExtraDetails) (*paymentpb.BasicApiResponse, error) {
	resp := &paymentpb.BasicApiResponse{Success: true}

	if extraDetails != nil {
		if !utils.ValidDate(extraDetails.GetExpiryDate()) {
			resp.Message = "Invalid Date"
			resp.Success = false
			return resp, nil
		}
		if utils.CheckForOlderDate(extraDetails.GetExpiryDate()) {
			resp.Message = "Cannot set older date as expiry date"
			resp.Success = false
			return resp, nil
		}
	}
	return resp, nil
}
