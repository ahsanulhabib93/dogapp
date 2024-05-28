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
		ID: "20240117104215",
		Migrate: func(tx *gorm.DB) error {
			if !tx.HasTable(models.SellerAccountManager{}) {
				if err := tx.AutoMigrate(models.SellerAccountManager{}).Error; err != nil {
					return err
				}
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})

	migrator.Register(&gormigrate.Migration{
		ID: "20240227152449",
		Migrate: func(tx *gorm.DB) error {
			tx.Exec("ALTER TABLE sellers MODIFY COLUMN data_mapping JSON, MODIFY COLUMN return_exchange_policy JSON")
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})
	migrator.Register(&gormigrate.Migration{
		ID: "20240227152450",
		Migrate: func(tx *gorm.DB) error {
			tx.Exec("ALTER TABLE seller_configs MODIFY COLUMN refund_policy JSON")
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})
}
