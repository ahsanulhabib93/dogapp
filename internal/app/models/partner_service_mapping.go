package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/utils"
)

type PartnerServiceMapping struct {
	database.VaccountGorm
	SupplierId            uint64            `gorm:"not null" json:"supplier_id"`
	ServiceType           utils.ServiceType `valid:"required" json:"service_type"`
	Active                bool              `gorm:"default:false"`
	TradeLicenseUrl       string            `gorm:"type:varchar(512)" json:"trade_license_url"`
	AgreementUrl          string            `gorm:"type:varchar(512)" json:"agreement_url"`
	PartnerServiceLevelID uint64            `valid:"required" json:"partner_service_level_id"`
}

func (partnerService *PartnerServiceMapping) Validate(db *gorm.DB) {
	result := db.Model(&partnerService).First(&PartnerServiceMapping{}, "id != ? and supplier_id = ? and service_type = ?", partnerService.ID, partnerService.SupplierId, partnerService.ServiceType)
	if !result.RecordNotFound() {
		db.AddError(errors.New("Partner Service Already Exists")) //nolint:errcheck
	}
}
