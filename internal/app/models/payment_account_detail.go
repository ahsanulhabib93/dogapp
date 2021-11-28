package models

import (
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/utils"
)

type PaymentAccountDetail struct {
	database.VaccountGorm
	SupplierID    uint64
	AccountType   utils.AccountType
	AccountName   string
	AccountNumber string
	BankName      string
	BranchName    string
	RoutingNumber string
	IsDefault     bool
	Supplier      Supplier
}
