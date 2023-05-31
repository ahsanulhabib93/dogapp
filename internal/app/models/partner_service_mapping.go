package models

import (
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/utils"
)

type PartnerServiceMapping struct {
	database.VaccountGorm
	SupplierId      uint64 `gorm:"not null" json:"supplier_id"`
	ServiceType     utils.ServiceType
	ServiceLevel    utils.SupplierType
	Active          bool   `gorm:"default:false"`
	TradeLicenseUrl string `gorm:"type:varchar(512)" json:"trade_license_url"`
	AgreementUrl    string `gorm:"type:varchar(512)" json:"agreement_url"`
}
