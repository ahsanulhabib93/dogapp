package models

import (
	"time"

	"github.com/voonik/goFramework/pkg/database"
)

type SupplierOpcMapping struct {
	database.VaccountGorm
	SupplierID         uint64 `gorm:"index:idx_supplier_id"`
	ProcessingCenterID uint64
	DeletedAt          *time.Time `sql:"index"`
}
