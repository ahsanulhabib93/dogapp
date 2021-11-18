package models

import (
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/supplier_service/internal/app/utils"
)

type Supplier struct {
	database.VaccountGorm
	Name                  string
	Email                 string
	SupplierType          utils.SupplierType
	SupplierAddresses     []SupplierAddress
	PaymentAccountDetails []PaymentAccountDetail
	KeyAccountManagers    []KeyAccountManager
}
