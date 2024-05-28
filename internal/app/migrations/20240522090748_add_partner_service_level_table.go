package migrations

import (
	"github.com/jinzhu/gorm"
	migrator "github.com/voonik/goFramework/pkg/migrations"
	"github.com/voonik/ss2/internal/app/models"
	"gopkg.in/gormigrate.v1"
)

func init() {
	migrator.Register(&gormigrate.Migration{
		ID: "20240522090748",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(models.PartnerServiceLevel{}).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})

	migrator.Register(&gormigrate.Migration{
		ID: "20240522090749",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(models.PartnerServiceMapping{}).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})
}
