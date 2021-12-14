package models

import (
	"github.com/voonik/goFramework/pkg/database"
)

type KeyAccountManager struct {
	database.VaccountGorm
	SupplierID uint64 `gorm:"index:idx_supplier_id" valid:"required"`
	Name       string `gorm:"not null" valid:"required"`
	Email      string
	Phone      string
}
