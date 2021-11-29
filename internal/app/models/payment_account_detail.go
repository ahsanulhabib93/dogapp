package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/utils"
)

type PaymentAccountDetail struct {
	database.VaccountGorm
	SupplierID    uint64 `gorm:"index:idx_supplier_id"`
	AccountType   utils.AccountType
	AccountName   string `gorm:"not null"`
	AccountNumber string `gorm:"not null"`
	BankName      string
	BranchName    string
	RoutingNumber string
	IsDefault     bool
	Supplier      Supplier
}

func (paymentAccount PaymentAccountDetail) Validate(db *gorm.DB) {
	if paymentAccount.AccountName == "" {
		db.AddError(errors.New("AccountName can't be blank"))
	}
	if paymentAccount.AccountNumber == "" {
		db.AddError(errors.New("AccountNumber can't be blank"))
	}
}
