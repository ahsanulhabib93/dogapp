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
			models := []interface{}{
				models.Seller{},
				models.VendorAddress{},
				models.SellerBankDetail{},
				models.SellerPricingDetail{},
				models.SellerConfig{},
				models.SellerActivityLog{},
			}

			for _, model := range models {
				if !tx.HasTable(model) {
					if err := tx.AutoMigrate(model).Error; err != nil {
						return err
					}
				}
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})

	migrator.Register(&gormigrate.Migration{
		ID: "20231116113708",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(
				models.SellerActivityLog{},
			).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})
	migrator.Register(&gormigrate.Migration{
		ID: "20231214175442",
		Migrate: func(tx *gorm.DB) error {
			models := []interface{}{
				models.SellerAccountManager{},
			}

			for _, model := range models {
				if !tx.HasTable(model) {
					if err := tx.AutoMigrate(model).Error; err != nil {
						return err
					}
				}
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})
}
