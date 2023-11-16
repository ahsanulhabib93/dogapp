package models

import (
	"github.com/voonik/goFramework/pkg/database"
	"gorm.io/datatypes"
)

type SellerActivityLog struct {
	database.VaccountGorm
	UserID            uint64
	SellerID          uint64
	Action            string
	SellerStateReason string
	Notes             datatypes.JSON `gorm:"type:json"`
}
