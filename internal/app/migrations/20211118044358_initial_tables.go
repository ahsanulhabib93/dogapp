package migrations

import (
	"github.com/jinzhu/gorm"
	aaaModels "github.com/voonik/goFramework/pkg/aaa/models"
	migrator "github.com/voonik/goFramework/pkg/migrations"
	"github.com/voonik/ss2/internal/app/models"
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
	migrator.Register(&gormigrate.Migration{
		ID: "20230625130730",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&aaaModels.AppPreference{}).AddUniqueIndex("pref_key", "preference_key", "vaccount_id", "portal_id").Error
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})
}
