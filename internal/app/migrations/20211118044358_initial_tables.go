package migrations

import (
	"github.com/jinzhu/gorm"
	aaaModels "github.com/voonik/goFramework/pkg/aaa/models"
	migrator "github.com/voonik/goFramework/pkg/migrations"
	"github.com/voonik/supplier_service/internal/app/models"
	"gopkg.in/gormigrate.v1"
)

func init() {
	migrator.Register(&gormigrate.Migration{
		ID: "20211118044358",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(
				aaaModels.Audit{},
				aaaModels.AppPreference{},
				aaaModels.CommonPreference{},
			).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})
	migrator.Register(&gormigrate.Migration{
		ID: "20211118054358",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(
				models.Supplier{},
				models.SupplierAddress{},
				models.KeyAccountManager{},
				models.PaymentAccountDetail{},
			).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})
}
