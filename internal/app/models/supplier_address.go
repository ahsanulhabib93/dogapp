package models

import (
	"errors"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/voonik/goFramework/pkg/database"
)

// SupplierAddress ...
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
	Zipcode    string
	Phone      string
	GstNumber  string `json:"gst_number"`
	IsDefault  bool   `json:"is_default"`
}

// Validate ...
func (supplierAddress *SupplierAddress) Validate(db *gorm.DB) {
	if supplierAddress.SupplierID == 0 {
		db.AddError(errors.New("SupplierID can't be blank"))
	}

	if phoneNumber := strings.TrimSpace(supplierAddress.Phone); len(phoneNumber) == 0 {
		db.AddError(errors.New("Phone Number can't be blank"))
	} else if !(strings.HasPrefix(phoneNumber, "8801") && len(phoneNumber) == 13) {
		db.AddError(errors.New("Phone Number should have 13 digits"))
	}
}

// AfterSave ...
func (supplierAddress *SupplierAddress) AfterSave(db *gorm.DB) error {
	supplier := Supplier{}
	db.Model(&supplier).First(&supplier, "id = ? ", supplierAddress.SupplierID)
	if supplier.Status == SupplierStatusVerified || supplier.Status == SupplierStatusFailed {
		db.Model(&supplier).Update("status", SupplierStatusPending)
	}
	return nil
}
