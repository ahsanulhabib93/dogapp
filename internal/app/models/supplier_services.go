package models

import (
	"github.com/voonik/goFramework/pkg/database"
)

type SupplierService struct {
	database.VaccountGorm
	SupplierId      uint64 `gorm:"not null" valid:"required" json:"supplier_id"`
	ServiceType     string `gorm:"type:varchar(512)" json:"service_type"`
	ServiceLevel    string `gorm:"type:varchar(512)" json:"service_level"`
	Active          bool   `gorm:"default:true" json:"active"`
	TradeLicenseUrl string `gorm:"type:varchar(512)" json:"trade_license_url"`
	AgreementUrl    string `gorm:"type:varchar(512)" json:"agreement_url"`

	Supplier *Supplier // belongs_to
}

// Unique index for supplier_id and service_type
