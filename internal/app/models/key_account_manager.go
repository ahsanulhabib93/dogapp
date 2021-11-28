package models

import (
	"github.com/voonik/goFramework/pkg/database"
)

type KeyAccountManager struct {
	database.VaccountGorm
	SupplierID uint64
	Name       string
	Email      string
	Phone      string
	Supplier   Supplier
}
