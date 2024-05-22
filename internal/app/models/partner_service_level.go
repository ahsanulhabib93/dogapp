package models

import (
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/utils"
)

type PartnerServiceLevel struct {
	database.VModel
	ServiceLevel utils.SupplierType `valid:"required" json:"service_level"`
	Name         string             `gorm:"type:varchar(255)" json:"name"`
	Active       bool               `gorm:"default:false" json:"active"`
}
