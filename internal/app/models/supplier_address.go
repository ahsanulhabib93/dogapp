package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/voonik/goFramework/pkg/database"
)

type SupplierAddress struct {
	database.VaccountGorm
	SupplierID uint64 `gorm:"index:idx_supplier_id"`
	Firstname  string
	Lastname   string
	Address1   string `gorm:"not null" valid:"required"`
	Address2   string
	Landmark   string
	City       string
	State      string
	Country    string
	Zipcode    string `gorm:"not null" valid:"required"`
	Phone      string
	GstNumber  string `json:"gst_number"`
	IsDefault  bool   `json:"is_default"`
}

func (supplierAddress SupplierAddress) Validate(db *gorm.DB) {
	if supplierAddress.SupplierID == 0 {
		db.AddError(errors.New("SupplierID can't be blank"))
	}
}
