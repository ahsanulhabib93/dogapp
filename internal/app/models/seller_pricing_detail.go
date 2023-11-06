package models

import (
	"time"

	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/utils"
)

type SellerPricingDetail struct {
	database.VaccountGorm
	IDType                    int
	FreeShippingLimit         int     `gorm:"default:0"`
	ShippingCharge            float64 `gorm:"default:28.0"`
	CODLimit                  int     `gorm:"default:0"`
	CODCharge                 float64 `gorm:"default:49.0"`
	DiscountPercent           float64
	CommissionPercent         float64   `gorm:"default:20.0;not null"`
	FlatShipping              bool      `gorm:"default:true"`
	StartDate                 time.Time `gorm:"default:'2015-05-28 01:43:53'"`
	EndDate                   time.Time
	DeletedAt                 *time.Time
	SellerID                  int
	LeadShippingDays          int `gorm:"default:2"`
	SellerPaymentType         int `gorm:"type:integer;size:2;default:1"`
	Config                    string
	SellerShippingCharge      float64                   `gorm:"default:67.0;not null"`
	SellerShippingPolicy      int                       `gorm:"type:integer;size:2"`
	Verified                  utils.SellerPriceVerified `gorm:"default:'NOT_VERIFIED'"`
	SellerShippingPercent     float64
	RTOCharges                float64 `gorm:"default:140.0;not null"`
	StockoutChargePercent     float64 `gorm:"default:2.0;not null"`
	CancellationChargePercent float64 `gorm:"default:4.0;not null"`
	TcsEnabled                bool    `gorm:"default:true"`
	TcsDeclarationFileName    string
	TcsDeclarationContentType string
	TcsDeclarationFileSize    int
	TcsDeclarationUpdatedAt   time.Time
}
