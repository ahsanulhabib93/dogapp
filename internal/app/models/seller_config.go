package models

import "github.com/voonik/goFramework/pkg/database"

type SellerConfig struct {
	database.VaccountGorm
	AllowPriceUpdate      bool `gorm:"default:true"`
	CODConfirmationNeeded bool `gorm:"default:true"`
	SellerID              int
	RefundPolicy          string
	ServiceCheckerConfig  string
	MaxQuantity           int `gorm:"default:50"`
	DeliveryCheckType     int `gorm:"default:0"`
	PenaltyPolicy         string
	ItemsPerPackage       int  `gorm:"type:integer;size:2;default:1;not null"`
	PickupType            int  `gorm:"default:1"`
	QCFrequency           int  `gorm:"type:integer;size:1;default:0;not null"`
	AllowVendorCoupons    bool `gorm:"default:true"`
	VaccountID            int
	TPEnabled             bool `gorm:"default:false"`
	Properties            string
	SellerStockEnabled    bool `gorm:"default:true"`
}
