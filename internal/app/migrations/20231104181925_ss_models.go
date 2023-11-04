package migrations

import (
	"github.com/jinzhu/gorm"
	migrator "github.com/voonik/goFramework/pkg/migrations"
	"github.com/voonik/ss2/internal/app/models"
	"gopkg.in/gormigrate.v1"
)

func init() {
	migrator.Register(&gormigrate.Migration{
		ID: "20231104181925",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(
				models.Seller{},
				models.VendorAddress{},
				models.SellerBankDetail{},
				models.SellerPricingDetail{},
				models.SellerConfig{},
			).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})
}
