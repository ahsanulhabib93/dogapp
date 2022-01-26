package migrations

import (
	"log"

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
			err := tx.Model(
				&models.PaymentAccountDetail{},
			).DropColumn("bank_name").Error
			if err != nil {
				log.Printf("ERROR: %v", err)
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})
	migrator.Register(&gormigrate.Migration{
		ID: "20220126055358",
		Migrate: func(tx *gorm.DB) error {
			err := tx.AutoMigrate(&models.Supplier{})
			if err != nil {
				log.Printf("ERROR: %v", err)
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			err := tx.Model(
				&models.Supplier{},
			).DropColumn("status").Error
			if err != nil {
				log.Printf("ERROR: %v", err)
			}
			return nil
		},
	})
}
