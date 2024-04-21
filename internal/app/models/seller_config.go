package models

import (
	"github.com/voonik/goFramework/pkg/database"
	"gorm.io/datatypes"
)

type SellerConfig struct {
	database.VaccountGorm
	AllowPriceUpdate      bool `gorm:"default:true"`
	CODConfirmationNeeded bool `gorm:"default:true"`
	SellerID              int
	RefundPolicy          datatypes.JSON `gorm:"type:json"`
	ServiceCheckerConfig  string
	MaxQuantity           int `gorm:"default:50"`
	DeliveryCheckType     int `gorm:"default:0"`
	PenaltyPolicy         string
	ItemsPerPackage       int  `gorm:"type:integer;size:2;default:1;not null"`
	PickupType            int  `gorm:"default:1"`
	QCFrequency           int  `gorm:"type:integer;size:1;default:0;not null"`
	AllowVendorCoupons    bool `gorm:"default:true"`
	TPEnabled             bool `gorm:"default:false"`
	Properties            string
	SellerStockEnabled    bool `gorm:"default:true"`
}
