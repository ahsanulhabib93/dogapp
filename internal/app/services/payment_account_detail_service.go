package services

import (
	"context"
	"fmt"
	"log"
	"strings"

	paymentpb "github.com/voonik/goConnect/api/go/ss2/payment_account_detail"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/helpers"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/utils"
)

// PaymentAccountDetailService ...
type PaymentAccountDetailService struct{}

// List ...
func (ps *PaymentAccountDetailService) List(ctx context.Context, params *paymentpb.ListParams) (*paymentpb.ListResponse, error) {
	log.Printf("ListPaymentAccountParams: %+v", params)
	resp := paymentpb.ListResponse{}
	database.DBAPM(ctx).Model(&models.PaymentAccountDetail{}).Joins(
		models.GetBankJoinStr()).Select("payment_account_details.*, banks.name bank_name").Where(
		"supplier_id = ?", params.GetSupplierId()).Scan(&resp.Data)
	return &resp, nil
}

// Add ...
func (ps *PaymentAccountDetailService) Add(ctx context.Context, params *paymentpb.PaymentAccountDetailParam) (*paymentpb.BasicApiResponse, error) {
	log.Printf("AddPaymentAccountParams: %+v", params)
	resp := paymentpb.BasicApiResponse{Success: false}

	supplier := &models.Supplier{}
	result := database.DBAPM(ctx).Model(&models.Supplier{}).First(supplier, params.GetSupplierId())
	if result.RecordNotFound() {
		resp.Message = "Supplier Not Found"
	} else if !supplier.IsChangeAllowed(ctx) {
		resp.Message = "Change Not Allowed"
	} else {
		paymentAccountDetail := models.PaymentAccountDetail{
			SupplierID:     supplier.ID,
			AccountType:    utils.AccountType(params.GetAccountType()),
			AccountSubType: utils.AccountSubType(params.GetAccountSubType()),
			AccountName:    params.GetAccountName(),
			AccountNumber:  params.GetAccountNumber(),
			BankID:         params.GetBankId(),
			BranchName:     params.GetBranchName(),
			RoutingNumber:  params.GetRoutingNumber(),
			IsDefault:      params.GetIsDefault(),
		}
		if params.GetExtraDetails() != nil {
			if !utils.ValidDate(params.GetExtraDetails().GetExpiryDate()) {
				resp.Message = "Invalid Date"
				return &resp, nil
			}
			if utils.CheckForOlderDate(params.GetExtraDetails().GetExpiryDate()) {
				resp.Message = "Cannot set older date as expiry date"
				return &resp, nil
			}
			extraDetails := models.PaymentAccountDetailExtraDetails{}
			utils.CopyStructAtoB(params.ExtraDetails, &extraDetails)
			paymentAccountDetail.SetExtraDetails(extraDetails)
		}
		err := database.DBAPM(ctx).Save(&paymentAccountDetail)

		if err != nil && err.Error != nil {
			resp.Message = fmt.Sprintf("Error while creating Payment Account Detail: %s", err.Error)
			return &resp, nil
		}
		if params.GetAccountType() == uint64(utils.PrepaidCard) {
			helpers.StoreEncryptCardInfo(ctx, *params.GetExtraDetails(), &paymentAccountDetail)
		}
		helpers.UpdateDefaultPaymentAccount(ctx, &paymentAccountDetail)
		resp.Message = "Payment Account Detail Added Successfully"
		resp.Success = true

	}
	log.Printf("AddPaymentAccountResponse: %+v", resp)
	return &resp, nil
}

// Edit ...
func (ps *PaymentAccountDetailService) Edit(ctx context.Context, params *paymentpb.PaymentAccountDetailObject) (*paymentpb.BasicApiResponse, error) {
	log.Printf("EditPaymentAccountParams: %+v", params)
	resp := paymentpb.BasicApiResponse{Success: false}

	paymentAccountDetail := models.PaymentAccountDetail{}
	result := database.DBAPM(ctx).Model(&models.PaymentAccountDetail{}).First(&paymentAccountDetail, params.GetId())
	if result.RecordNotFound() {
		resp.Message = "PaymentAccountDetail Not Found"
	} else {
		supplier := models.Supplier{}
		database.DBAPM(ctx).Model(&models.Supplier{}).First(&supplier, paymentAccountDetail.SupplierID)
		if !supplier.IsChangeAllowed(ctx) {
			resp.Message = "Change Not Allowed"
		} else {
			// extra details validation
			if params.GetExtraDetails() != nil {
				if !utils.ValidDate(params.GetExtraDetails().GetExpiryDate()) {
					resp.Message = "Invalid Date"
					return &resp, nil
				}
				if utils.CheckForOlderDate(params.GetExtraDetails().GetExpiryDate()) {
					resp.Message = "Cannot set older date as expiry date"
					return &resp, nil
				}
				extraDetails := models.PaymentAccountDetailExtraDetails{}
				utils.CopyStructAtoB(params.ExtraDetails, &extraDetails)
				paymentAccountDetail.SetExtraDetails(extraDetails)
			}
			err := database.DBAPM(ctx).Model(&paymentAccountDetail).Updates(models.PaymentAccountDetail{
				AccountType:    utils.AccountType(params.GetAccountType()),
				AccountSubType: utils.AccountSubType(params.GetAccountSubType()),
				AccountName:    params.GetAccountName(),
				AccountNumber:  params.GetAccountNumber(),
				BankID:         params.GetBankId(),
				BranchName:     params.GetBranchName(),
				RoutingNumber:  params.GetRoutingNumber(),
				IsDefault:      params.GetIsDefault(),
			})
			if err != nil && err.Error != nil {
				resp.Message = fmt.Sprintf("Error while updating PaymentAccountDetail: %s", err.Error)
				return &resp, nil
			}
			if params.GetAccountType() == uint64(utils.PrepaidCard) {
				helpers.StoreEncryptCardInfo(ctx, *params.GetExtraDetails(), &paymentAccountDetail)
			}
			helpers.UpdateDefaultPaymentAccount(ctx, &paymentAccountDetail)
			resp.Message = "PaymentAccountDetail Edited Successfully"
			resp.Success = true
		}
	}
	log.Printf("EditPaymentAccountResponse: %+v", resp)
	return &resp, nil
}

// ListBanks ...
func (ps *PaymentAccountDetailService) ListBanks(ctx context.Context, params *paymentpb.ListBankParams) (*paymentpb.ListBankResponse, error) {
	log.Printf("ListBanksParams: %+v", params)
	resp := paymentpb.ListBankResponse{}
	database.DBAPM(ctx).Model(&models.Bank{}).Scan(&resp.Data)
	return &resp, nil
}

// MapPaymentAccountDetail ...
func (ps *PaymentAccountDetailService) MapPaymentAccountDetail(ctx context.Context, params *paymentpb.MappingParam) (*paymentpb.BasicApiResponse, error) {
	resp := &paymentpb.BasicApiResponse{}
	switch strings.ToLower(params.MappableType) {
	case "warehouses":
		err := helpers.UpdatePaymentAccountDetailWarehouseMapping(ctx, params.GetId(), params.GetMappableIds(), params.GetWarehouseDhCodeMap())
		if err != nil {
			resp.Message = err.Error()
		} else {
			resp.Success = true
			resp.Message = "Mapping Updated Successfully"
		}
	default:
		resp.Message = "Invalid mapping_type"
	}
	return resp, nil
}
