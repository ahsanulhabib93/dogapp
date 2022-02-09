package models

import (
	"errors"
	"strings"

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
	Phone      string `gorm:"not null" valid:"required"`
	GstNumber  string `json:"gst_number"`
	IsDefault  bool   `json:"is_default"`
}

func (supplierAddress SupplierAddress) Validate(db *gorm.DB) {
	if supplierAddress.SupplierID == 0 {
		db.AddError(errors.New("SupplierID can't be blank"))
	}

	if phoneNumber := strings.TrimSpace(supplierAddress.Phone); len(phoneNumber) == 0 {
		db.AddError(errors.New("Phone Number can't be blank"))
	} else if !((strings.HasPrefix(phoneNumber, "8801") && len(phoneNumber) == 13) ||
		(strings.HasPrefix(phoneNumber, "01") && len(phoneNumber) == 11) ||
		(strings.HasPrefix(phoneNumber, "1") && len(phoneNumber) == 10)) {
		db.AddError(errors.New("Invalid Phone Number"))
	}
}
