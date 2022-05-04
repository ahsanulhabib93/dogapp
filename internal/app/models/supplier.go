package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/utils"
)

// SupplierStatusConstants ...
type SupplierStatus string

const (
	SupplierStatusPending  SupplierStatus = "Pending"
	SupplierStatusVerified SupplierStatus = "Verified"
	SupplierStatusFailed   SupplierStatus = "Failed"
	SupplierStatusBlocked  SupplierStatus = "Blocked"
)

// Supplier ...
type Supplier struct {
	database.VaccountGorm
	Name                     string         `gorm:"not null" valid:"required"`
	Status                   SupplierStatus `gorm:"default:'Pending'"`
	Reason                   string
	Email                    string
	Phone                    string
	AlternatePhone           string                 `json:"alternate_phone"`
	BusinessName             string                 `json:"business_name"`
	IsPhoneVerified          bool                   `gorm:"default:false" json:"is_phone_verified"`
	ShopImageURL             string                 `json:"shop_image_url"`
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

	// if !supplier.Status.IsValid() {
	// 	db.AddError(errors.New("Status should be Active/Pending/Deactivated"))
	// }
}

// GetCategoryMappingJoinStr ...
func GetCategoryMappingJoinStr() string {
	return "LEFT JOIN supplier_category_mappings on supplier_category_mappings.supplier_id = suppliers.id and supplier_category_mappings.deleted_at IS NULL and supplier_category_mappings.vaccount_id = suppliers.vaccount_id"
}

// GetOpcMappingJoinStr ...
func GetOpcMappingJoinStr() string {
	return "LEFT JOIN supplier_opc_mappings on supplier_opc_mappings.supplier_id = suppliers.id and supplier_opc_mappings.deleted_at IS NULL and supplier_opc_mappings.vaccount_id = suppliers.vaccount_id"
}

// var supplierStatusTransitionState = map[SupplierStatus][]SupplierStatus{
// 	SupplierStatusPending:    {SupplierStatusActive, SupplierStatusDeactivate},
// 	SupplierStatusActive:     {SupplierStatusDeactivate},
// 	SupplierStatusDeactivate: {SupplierStatusActive},
// }

// func (s SupplierStatus) IsValid() bool {
// 	_, found := supplierStatusTransitionState[s]
// 	return len(s) == 0 || found
// }

// func (s SupplierStatus) IsTransitionAllowed(status SupplierStatus) bool {
// 	if s == status {
// 		return true
// 	}

// 	for _, next := range supplierStatusTransitionState[s] {
// 		if next == status {
// 			return true
// 		}
// 	}

// 	return false
// }
