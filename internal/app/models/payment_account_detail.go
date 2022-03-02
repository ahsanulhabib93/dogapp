package models

import (
	"errors"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/utils"
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
}

func (paymentAccount PaymentAccountDetail) Validate(db *gorm.DB) {
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
		utils.Bank: {utils.Current, utils.Savings},
		utils.Mfs:  {utils.Bkash, utils.Nagada},
	}
}
