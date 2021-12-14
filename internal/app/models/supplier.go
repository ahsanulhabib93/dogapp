package models

import (
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/utils"
)

type Supplier struct {
	database.VaccountGorm
	Name                  string `gorm:"not null" valid:"required"`
	Email                 string
	SupplierType          utils.SupplierType `json:"supplier_type" valid:"required"`
	SupplierAddresses     []SupplierAddress  `json:"supplier_addresses"`
	PaymentAccountDetails []PaymentAccountDetail
	KeyAccountManagers    []KeyAccountManager
}
