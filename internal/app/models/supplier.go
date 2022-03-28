package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/utils"
)

// SupplierStatusConstants ...
const (
	SupplierStatusActive  = "Active"
	SupplierStatusPending = "Pending"
)

// Supplier ...
type Supplier struct {
	database.VaccountGorm
	Name                     string `gorm:"not null" valid:"required"`
	Status                   string `gorm:"default:'Pending'"`
	Email                    string
	UserID                   *uint64                `json:"user_id"`
	SupplierType             utils.SupplierType     `json:"supplier_type" valid:"required"`
	SupplierAddresses        []SupplierAddress      `json:"supplier_addresses"`
	PaymentAccountDetails    []PaymentAccountDetail `json:"payment_account_details"`
	KeyAccountManagers       []KeyAccountManager
	SupplierCategoryMappings []SupplierCategoryMapping
	SupplierOpcMappings      []SupplierOpcMapping
}

// Validate ...
func (supplier Supplier) Validate(db *gorm.DB) {
	s := &Supplier{}
	result := db.Model(&supplier).First(s, "name = ?", supplier.Name)
	if !result.RecordNotFound() && s.ID != supplier.ID {
		db.AddError(errors.New("Supplier Already Exists, please contact with the admin team to get access"))
	}

	if !(supplier.Status == SupplierStatusActive || supplier.Status == SupplierStatusPending) &&
		len(supplier.Status) > 0 {
		db.AddError(errors.New("Status should be Active/Pending"))
	}
}

// GetCategoryMappingJoinStr ...
func GetCategoryMappingJoinStr() string {
	return "LEFT JOIN supplier_category_mappings on supplier_category_mappings.supplier_id = suppliers.id and supplier_category_mappings.deleted_at IS NULL and supplier_category_mappings.vaccount_id = suppliers.vaccount_id"
}

// GetOpcMappingJoinStr ...
func GetOpcMappingJoinStr() string {
	return "LEFT JOIN supplier_opc_mappings on supplier_opc_mappings.supplier_id = suppliers.id and supplier_opc_mappings.deleted_at IS NULL and supplier_opc_mappings.vaccount_id = suppliers.vaccount_id"
}
