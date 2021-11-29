package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/voonik/goFramework/pkg/database"
)

type KeyAccountManager struct {
	database.VaccountGorm
	SupplierID uint64 `gorm:"index:idx_supplier_id"`
	Name       string `gorm:"not null"`
	Email      string
	Phone      string
	Supplier   Supplier
}

func (kam KeyAccountManager) Validate(db *gorm.DB) {
	if kam.Name == "" {
		db.AddError(errors.New("Name can't be blank"))
	}
}
