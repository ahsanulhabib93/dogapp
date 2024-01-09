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
	paymentAccountDetails := []*models.PaymentAccountDetail{}
	if warehouseID != 0 {
		database.DBAPM(ctx).Model(&models.PaymentAccountDetail{}).Preload("PaymentAccountDetailWarehouseMappings").Joins(models.JoinPaymentAccountDetailWarehouseMappings()).Where("warehouse_id = ?", warehouseID).Where("supplier_id = ?", supplier.ID).Scan(&paymentAccountDetails)
	} else {
		database.DBAPM(ctx).Model(&models.PaymentAccountDetail{}).Where("supplier_id = ?", supplier.ID).Scan(&paymentAccountDetails)
	}

	var paymentDetailIds []uint64
	for _, paymentDetail := range paymentAccountDetails {
		paymentDetailIds = append(paymentDetailIds, paymentDetail.ID)
	}

	paymentResponse := []*supplierpb.PaymentAccountDetailObject{}

	warehouseDhCodeMap := GetWarehouseDhCodeForPaymentAccountDetails(ctx, paymentDetailIds)
	for _, paymentDetail := range paymentAccountDetails {

		warehouses := []uint64{}
		for whId := range warehouseDhCodeMap[paymentDetail.ID] {
			warehouses = append(warehouses, whId)
		}

		bank := models.Bank{}
		database.DBAPM(ctx).Model(&models.Bank{}).Where("banks.id = ?", paymentDetail.BankID).Scan(&bank)

		extraDetails := &supplierpb.ExtraDetails{}
		utils.CopyStructAtoB(paymentDetail.ExtraDetails, extraDetails)
		paymentResponse = append(paymentResponse, &supplierpb.PaymentAccountDetailObject{
			Id:                 paymentDetail.ID,
			SupplierId:         paymentDetail.SupplierID,
			AccountType:        uint64(paymentDetail.AccountType),
			AccountSubType:     uint64(paymentDetail.AccountSubType),
			AccountNumber:      paymentDetail.AccountNumber,
			AccountName:        paymentDetail.AccountName,
			BankName:           bank.Name,
			BranchName:         paymentDetail.BranchName,
			RoutingNumber:      paymentDetail.RoutingNumber,
			IsDefault:          paymentDetail.IsDefault,
			BankId:             bank.ID,
			Warehouses:         warehouses,
			ExtraDetails:       extraDetails,
			WarehouseDhCodeMap: warehouseDhCodeMap[paymentDetail.ID],
			DhCode:             warehouseDhCodeMap[paymentDetail.ID][warehouseID].GetDhCode(),
		})
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

func PrepaidCardValidations(ctx context.Context, extraDetails paymentpb.ExtraDetails, paymentAccountDetails *models.PaymentAccountDetail, accountNumber string) (bool, string) {
	extraDetailStruct := models.PaymentAccountDetailExtraDetails{}
	isError, errMsg := HandleExtraDetailsValidation(&extraDetails)

	if isError {
		database.DBAPM(ctx).Delete(paymentAccountDetails)
		return true, errMsg
	}

	extraDetailStruct.EmployeeId = extraDetails.EmployeeId
	extraDetailStruct.ExpiryDate = extraDetails.ExpiryDate
	extraDetailStruct.ClientId = extraDetails.ClientId

	success := StoreEncryptCardInfo(ctx, &extraDetailStruct, paymentAccountDetails, accountNumber)
	if !success {
		return true, "Failed to create Paywell Card"
	}
	paymentAccountDetails.ExtraDetails = extraDetailStruct
	database.DBAPM(ctx).Save(&paymentAccountDetails)

	return false, utils.EmptyString
}

func StoreEncryptCardInfo(ctx context.Context, extraDetails *models.PaymentAccountDetailExtraDetails, paymentAccountDetail *models.PaymentAccountDetail, accountNumber string) bool {
	uniqueId := utils.CreatePaywellUniqueKey(paymentAccountDetail.ID)
	expiryMonth, expiryYear := utils.FetchMonthAndYear(extraDetails.ExpiryDate)
	logger.FromContext(ctx).Infof("Payload for CreatePaywellCard : unique id %v, card info %v, expiry month %v, expiry year %v", uniqueId, paymentAccountDetail.AccountNumber, expiryMonth, expiryYear)
	paywellResponse := getAPIHelperInstance().CreatePaywellCard(ctx, &paywellPb.CreateCardRequest{
		UniqueId:    uniqueId,
		CardInfo:    accountNumber,
		ExpiryMonth: expiryMonth,
		ExpiryYear:  expiryYear,
	})

	if paywellResponse.IsError {
		database.DBAPM(ctx).Delete(paymentAccountDetail)
		return false
	}
	paymentAccountDetail.AccountNumber = paywellResponse.MaskedNumber
	extraDetails.UniqueId = uniqueId
	extraDetails.Token = paywellResponse.GetToken()
	return true
}

func HandleExtraDetailsValidation(extraDetails *paymentpb.ExtraDetails) (bool, string) {
	if extraDetails != nil {
		if extraDetails.EmployeeId == utils.Zero {
			return true, "Employee ID is mandatory"
		}
		if !utils.ValidDate(extraDetails.GetExpiryDate()) {
			return true, "Invalid Date"
		}
		if utils.CheckForOlderDate(extraDetails.GetExpiryDate()) {
			return true, "Cannot set older date as expiry date"
		}
	}
	return false, utils.EmptyString
}
