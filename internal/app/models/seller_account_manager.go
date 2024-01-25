package models

import (
	"github.com/voonik/goFramework/pkg/database"
	"gorm.io/gorm"
)

type SellerAccountManager struct {
	database.VaccountGorm
	SellerID  int
	Role      string
	DeletedAt gorm.DeletedAt
	Priority  int
	Phone     int64
	Name      string
	Email     string
	Seller    *Seller
}
