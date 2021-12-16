package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/utils"
)

type PaymentAccountDetail struct {
	database.VaccountGorm
	SupplierID    uint64            `gorm:"index:idx_supplier_id" valid:"required"`
	AccountType   utils.AccountType `valid:"required"`
	AccountName   string            `gorm:"not null" valid:"required"`
	AccountNumber string            `gorm:"not null" valid:"required"`
	BankName      string
	BranchName    string
	RoutingNumber string
	IsDefault     bool
}

func (paymentAccount PaymentAccountDetail) Validate(db *gorm.DB) {
	if !paymentAccount.IsDefault {
		result := db.Model(&paymentAccount).Where("supplier_id = ? and is_default = ? and id != ?", paymentAccount.SupplierID, true, paymentAccount.ID).First(&PaymentAccountDetail{})
		if result.RecordNotFound() {
			db.AddError(errors.New("Default Payment Account is required"))
		}
	}
}
