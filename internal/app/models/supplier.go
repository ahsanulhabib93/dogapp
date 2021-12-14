package models

import (
	"errors"

	"github.com/jinzhu/gorm"
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

func (supplier Supplier) Validate(db *gorm.DB) {
	result := db.Model(&supplier).First(&Supplier{}, "name = ?", supplier.Name)
	if !result.RecordNotFound() {
		db.AddError(errors.New("Name should be unique"))
	}
}
