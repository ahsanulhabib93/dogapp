package migrations

import (
	"github.com/jinzhu/gorm"
	migrator "github.com/voonik/goFramework/pkg/migrations"
	"github.com/voonik/ss2/internal/app/models"
	"gopkg.in/gormigrate.v1"
)

func init() {
	migrator.Register(&gormigrate.Migration{
		ID: "20220112044358",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(
				models.PaymentAccountDetail{},
				models.Bank{},
			).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})
	migrator.Register(&gormigrate.Migration{
		ID: "20220112055358",
		Migrate: func(tx *gorm.DB) error {
			return tx.Model(
				&models.PaymentAccountDetail{},
			).DropColumn("bank_name").Error
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})
}
