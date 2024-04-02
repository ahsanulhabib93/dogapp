package models

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"

	"github.com/jinzhu/gorm"
	aaaModels "github.com/voonik/goFramework/pkg/aaa/models"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/utils"
	gormIO "gorm.io/gorm"
)

type PaymentAccountDetail struct {
	database.VaccountGorm
	SupplierID                            uint64                           `gorm:"index:idx_supplier_id" valid:"required"`
	AccountType                           utils.AccountType                `valid:"required" json:"account_type,omitempty"`
	AccountSubType                        utils.AccountSubType             `json:"account_sub_type,omitempty"`
	AccountName                           string                           `gorm:"not null" valid:"required" json:"account_name,omitempty"`
	AccountNumber                         string                           `gorm:"not null" json:"account_number,omitempty"`
	BankID                                uint64                           `json:"bank_id,omitempty"`
	BranchName                            string                           `json:"branch_name,omitempty"`
	RoutingNumber                         string                           `json:"routing_number,omitempty"`
	IsDefault                             bool                             `json:"is_default,omitempty"`
	ExtraDetails                          PaymentAccountDetailExtraDetails `gorm:"type:json"`
	DeletedAt                             gormIO.DeletedAt                 `json:"deleted_at,omitempty"`
	PaymentAccountDetailWarehouseMappings []*PaymentAccountDetailWarehouseMapping
}

type PaymentAccountDetailExtraDetails struct {
	EmployeeId uint64 `json:"employee_id,omitempty"`
	ClientId   uint64 `json:"client_id,omitempty"`
	ExpiryDate string `json:"expiry_date,omitempty"`
	Token      string `json:"token,omitempty"`
	UniqueId   string `json:"unique_id,omitempty"`
}

func (extraDetails PaymentAccountDetailExtraDetails) Value() (driver.Value, error) {
	jsonBytes, err := json.Marshal(extraDetails)
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}

func (extraDetails *PaymentAccountDetailExtraDetails) Scan(value interface{}) error {
	jsonBytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(jsonBytes, &extraDetails)
}

// Validate ...
func (paymentAccount PaymentAccountDetail) Validate(db *gorm.DB) {
	if ctxx, ok := db.Get("context"); ok {
		if aaaModels.GetAppPreferenceServiceInstance().GetValue(ctxx.(context.Context), "enabled_account_number_validation", false).(bool) {
			res := db.Model(&paymentAccount).First(&PaymentAccountDetail{}, "supplier_id!= ? and account_number = ?", paymentAccount.SupplierID, paymentAccount.AccountNumber)
			if !res.RecordNotFound() {
				db.AddError(errors.New("Provided bank account number already exists")) //nolint:errcheck
			}
		}
	}

	if !paymentAccount.IsDefault {
		result := db.Model(&paymentAccount).Where("supplier_id = ? and is_default = ? and id != ?", paymentAccount.SupplierID, true, paymentAccount.ID).First(&PaymentAccountDetail{})
		if result.RecordNotFound() {
			db.AddError(errors.New("Default Payment Account is required")) //nolint:errcheck
		}
	}

	if !paymentAccount.validAccountSubType() {
		db.AddError(errors.New("Invalid Account SubType")) //nolint:errcheck
	}

	if paymentAccount.BankID != 0 {
		result := db.Model(&Bank{}).First(&Bank{}, paymentAccount.BankID)
		if result.RecordNotFound() {
			db.AddError(errors.New("Invalid Bank Name")) //nolint:errcheck
		}
	}

	if paymentAccount.AccountType == utils.Bank &&
		(paymentAccount.BankID == 0 || len(strings.TrimSpace(paymentAccount.BranchName)) == 0) {
		db.AddError(errors.New("For Bank account type BankID and BranchName needed")) //nolint:errcheck
	}

	if !paymentAccount.validAccountNumber() {
		db.AddError(errors.New("account_number is required"))
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

func (paymentAccount *PaymentAccountDetail) validAccountSubType() bool {
	mapping := paymentAccount.accountTypeMapping()
	subTypes := mapping[paymentAccount.AccountType]
	if len(subTypes) == utils.Zero {
		return true
	}
	for _, accountSubType := range subTypes {
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

func (paymentAccount PaymentAccountDetail) validAccountNumber() bool {
	if paymentAccount.AccountType != utils.Cheque {
		return len(strings.TrimSpace(paymentAccount.AccountNumber)) != utils.Zero
	}
	return true
}

func JoinPaymentAccountDetailWarehouseMappings() string {
	return "JOIN payment_account_detail_warehouse_mappings ON payment_account_detail_warehouse_mappings.payment_account_detail_id = payment_account_details.id AND payment_account_detail_warehouse_mappings.vaccount_id = payment_account_details.vaccount_id"
}
