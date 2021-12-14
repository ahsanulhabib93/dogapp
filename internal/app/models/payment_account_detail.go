package models

import (
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
