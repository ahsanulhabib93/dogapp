package models

import (
	"github.com/voonik/goFramework/pkg/database"
)

type SellerActivityLog struct {
	database.VaccountGorm
	UserID            uint64
	SellerID          uint64
	Action            string
	SellerStateReason string
	Notes             string `gorm:"type:text"`
}
