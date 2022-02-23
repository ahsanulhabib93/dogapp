package migrations

import (
	"github.com/jinzhu/gorm"
	migrator "github.com/voonik/goFramework/pkg/migrations"
	"github.com/voonik/ss2/internal/app/models"
	"gopkg.in/gormigrate.v1"
)

func init() {
	migrator.Register(&gormigrate.Migration{
		ID: "20220130045100",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(
				models.Supplier{},
				models.SupplierCategoryMapping{},
			).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})

	migrator.Register(&gormigrate.Migration{
		ID: "20220201045200",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(
				models.Supplier{},
				models.SupplierSaMapping{},
			).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})

	migrator.Register(&gormigrate.Migration{
		ID: "20220217230000",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(
				models.Supplier{},
				models.SupplierOpcMapping{},
			).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})
}
