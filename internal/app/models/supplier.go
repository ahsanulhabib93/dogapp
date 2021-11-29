package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/utils"
)

type Supplier struct {
	database.VaccountGorm
	Name                  string `gorm:"not null"`
	Email                 string
	SupplierType          utils.SupplierType `json:"supplier_type"`
	SupplierAddresses     []SupplierAddress  `json:"supplier_addresses"`
	PaymentAccountDetails []PaymentAccountDetail
	KeyAccountManagers    []KeyAccountManager
}

func (supplier Supplier) Validate(db *gorm.DB) {
	if supplier.Name == "" {
		db.AddError(errors.New("Name can't be blank"))
	}
}
