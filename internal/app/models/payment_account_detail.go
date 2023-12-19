package models

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/shopuptech/go-libs/logger"
	paymentpb "github.com/voonik/goConnect/api/go/ss2/payment_account_detail"
	aaaModels "github.com/voonik/goFramework/pkg/aaa/models"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/utils"
	"gorm.io/datatypes"
)

type PaymentAccountDetail struct {
	database.VaccountGorm
	SupplierID     uint64               `gorm:"index:idx_supplier_id" valid:"required"`
	AccountType    utils.AccountType    `valid:"required" json:"account_type,omitempty"`
	AccountSubType utils.AccountSubType `valid:"required" json:"account_sub_type,omitempty"`
	AccountName    string               `gorm:"not null" valid:"required" json:"account_name,omitempty"`
	AccountNumber  string               `gorm:"not null" valid:"required" json:"account_number,omitempty"`
	BankID         uint64               `json:"bank_id,omitempty"`
	BranchName     string               `json:"branch_name,omitempty"`
	RoutingNumber  string               `json:"routing_number,omitempty"`
	IsDefault      bool                 `json:"is_default,omitempty"`
	ExtraDetails   datatypes.JSON       `gorm:"type:json"`

	PaymentAccountDetailWarehouseMappings []*PaymentAccountDetailWarehouseMapping
}

// Validate ...
func (paymentAccount PaymentAccountDetail) Validate(db *gorm.DB) {
	if ctxx, ok := db.Get("context"); ok {
		if aaaModels.GetAppPreferenceServiceInstance().GetValue(ctxx.(context.Context), "enabled_account_number_validation", false).(bool) {
			res := db.Model(&paymentAccount).First(&PaymentAccountDetail{}, "supplier_id!= ? and account_number = ?", paymentAccount.SupplierID, paymentAccount.AccountNumber)
			if !res.RecordNotFound() {
				db.AddError(errors.New("Provided bank account number already exists"))
			}
		}
	}

	if !paymentAccount.IsDefault {
		result := db.Model(&paymentAccount).Where("supplier_id = ? and is_default = ? and id != ?", paymentAccount.SupplierID, true, paymentAccount.ID).First(&PaymentAccountDetail{})
		if result.RecordNotFound() {
			db.AddError(errors.New("Default Payment Account is required"))
		}
	}

	if !paymentAccount.validAccountSubType() {
		db.AddError(errors.New("Invalid Account SubType"))
	}

	if paymentAccount.BankID != 0 {
		result := db.Model(&Bank{}).First(&Bank{}, paymentAccount.BankID)
		if result.RecordNotFound() {
			db.AddError(errors.New("Invalid Bank Name"))
		}
	}

	if paymentAccount.AccountType == utils.Bank &&
		(paymentAccount.BankID == 0 || len(strings.TrimSpace(paymentAccount.BranchName)) == 0) {
		db.AddError(errors.New("For Bank account type BankID and BranchName needed"))
	}
}

// AfterSave ...
func (paymentAccount *PaymentAccountDetail) AfterSave(db *gorm.DB) error {
	supplier := Supplier{}
	db.Model(&supplier).First(&supplier, "id = ? ", paymentAccount.SupplierID)
	if supplier.Status == SupplierStatusVerified || supplier.Status == SupplierStatusFailed {
		db.Model(&supplier).Update("status", SupplierStatusPending)
	}
	return nil
}

func (paymentAccount PaymentAccountDetail) validAccountSubType() bool {
	mapping := paymentAccount.accountTypeMapping()
	for _, accountSubType := range mapping[paymentAccount.AccountType] {
		if accountSubType == paymentAccount.AccountSubType {
			return true
		}
	}
	return false
}

func (paymentAccount PaymentAccountDetail) accountTypeMapping() map[utils.AccountType][]utils.AccountSubType {
	return map[utils.AccountType][]utils.AccountSubType{
		utils.Bank:        {utils.Current, utils.Savings},
		utils.Mfs:         {utils.Bkash, utils.Nagada},
		utils.PrepaidCard: {utils.UCBL, utils.EBL},
	}
}

func (paymentAccount *PaymentAccountDetail) GetExtraDetails() *paymentpb.ExtraDetails {
	if paymentAccount.ExtraDetails == nil {
		return &paymentpb.ExtraDetails{}
	}

	ExtraDetails := &paymentpb.ExtraDetails{}
	err := json.Unmarshal(paymentAccount.ExtraDetails, ExtraDetails)
	if err != nil {
		logger.Log().Errorf("GetExtraDetails error = ", err.Error())
	}
	return ExtraDetails
}

func (paymentAccount *PaymentAccountDetail) SetExtraDetails(updatedExtraDetails paymentpb.ExtraDetails) *PaymentAccountDetail {
	var existingExtraDetails paymentpb.ExtraDetails
	if paymentAccount.ExtraDetails != nil {
		if err := json.Unmarshal(paymentAccount.ExtraDetails, &existingExtraDetails); err != nil {
			fmt.Printf("Error decoding existing ExtraDetails: %v\n", err)
			return paymentAccount
		}
	}
	mergeExtraDetails(&existingExtraDetails, updatedExtraDetails)

	updatedJSON, err := json.Marshal(existingExtraDetails)
	if err != nil {
		fmt.Printf("Error marshaling updated ExtraDetails: %v\n", err)
		return paymentAccount
	}

	paymentAccount.ExtraDetails = updatedJSON
	return paymentAccount
}

// Function to merge two ExtraDetails structs
func mergeExtraDetails(existing *paymentpb.ExtraDetails, updated paymentpb.ExtraDetails) {
	if updated.EmployeeId != 0 {
		existing.EmployeeId = updated.EmployeeId
	}
	if updated.ClientId != 0 {
		existing.ClientId = updated.ClientId
	}
	if updated.ExpiryDate != "" {
		existing.ExpiryDate = updated.ExpiryDate
	}
	if updated.UniqueId != "" {
		existing.UniqueId = updated.UniqueId
	}
	if updated.Token != "" {
		existing.Token = updated.Token
	}
}

func JoinPaymentAccountDetailWarehouseMappings() string {
	return "JOIN payment_account_detail_warehouse_mappings ON payment_account_detail_warehouse_mappings.payment_account_detail_id = payment_account_details.id AND payment_account_detail_warehouse_mappings.vaccount_id = payment_account_details.vaccount_id"
}
