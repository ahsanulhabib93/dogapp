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
	Address1   string `gorm:"not null"`
	Address2   string
	Landmark   string
	City       string
	State      string
	Country    string
	Zipcode    string `gorm:"not null"`
	Phone      string
	GstNumber  string `json:"gst_number"`
	IsDefault  bool   `json:"is_default"`
	Supplier   Supplier
}

func (supplierAddress SupplierAddress) Validate(db *gorm.DB) {
	if supplierAddress.Address1 == "" {
		db.AddError(errors.New("Address1 can't be blank"))
	}
	if supplierAddress.Zipcode == "" {
		db.AddError(errors.New("Zipcode can't be blank"))
	}
}
