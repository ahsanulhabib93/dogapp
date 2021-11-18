package models

import (
	"github.com/voonik/goFramework/pkg/database"
)

type KeyAccountManager struct {
	database.VaccountGorm
	Name     string
	Email    string
	Phone    string
	Supplier Supplier
}
