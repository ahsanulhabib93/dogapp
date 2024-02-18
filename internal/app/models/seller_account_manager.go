package models

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/voonik/goFramework/pkg/database"
	gormio "gorm.io/gorm"
)

type SellerAccountManager struct {
	database.VaccountGorm
	SellerID  uint64
	Role      string
	DeletedAt gormio.DeletedAt
	Priority  int
	Phone     int64
	Name      string
	Email     string
	Seller    *Seller
}

func (sam *SellerAccountManager) Validate(db *gorm.DB) {
	if sam.SellerID == 0 {
		db.AddError(errors.New("SellerID can't be blank"))
	}
	if sam.Role == "" {
		db.AddError(errors.New("Role can't be blank"))
	}
	if sam.Name == "" {
		db.AddError(errors.New("Name can't be blank"))
	}
	if phoneNumber := fmt.Sprint(sam.Phone); !(strings.HasPrefix(phoneNumber, "8801") && len(phoneNumber) == 13) {
		db.AddError(errors.New("Phone Number should have 13 digits"))
	}
}
