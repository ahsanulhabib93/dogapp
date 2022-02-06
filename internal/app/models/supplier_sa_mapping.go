package models

import (
	"time"

	"github.com/voonik/goFramework/pkg/database"
)

type SupplierSaMapping struct {
	database.VaccountGorm
	SupplierID          uint64
	SourcingAssociateId uint64
	DeletedAt           *time.Time `sql:"index"`
}
