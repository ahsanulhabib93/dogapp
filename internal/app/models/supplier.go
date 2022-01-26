package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/utils"
)

const (
	SupplierStatusActive  = "Active"
	SupplierStatusPending = "Pending"
)

type Supplier struct {
	database.VaccountGorm
	Name                  string `gorm:"not null" valid:"required"`
	Email                 string
	Status                string
	SupplierType          utils.SupplierType `json:"supplier_type" valid:"required"`
	SupplierAddresses     []SupplierAddress  `json:"supplier_addresses"`
	PaymentAccountDetails []PaymentAccountDetail
	KeyAccountManagers    []KeyAccountManager
}

func (supplier Supplier) Validate(db *gorm.DB) {
	s := &Supplier{}
	result := db.Model(&supplier).First(s, "name = ?", supplier.Name)
	if !result.RecordNotFound() && s.ID != supplier.ID {
		db.AddError(errors.New("Name should be unique"))
	}

	if !(supplier.Status == SupplierStatusActive || supplier.Status == SupplierStatusPending) &&
		len(supplier.Status) > 0 {
		db.AddError(errors.New("Status should be Active/Pending"))
	}
}
