package migrations

import (
	"github.com/jinzhu/gorm"
	migrator "github.com/voonik/goFramework/pkg/migrations"
	"github.com/voonik/ss2/internal/app/models"
	"gopkg.in/gormigrate.v1"
)

func init() {
	migrator.Register(&gormigrate.Migration{
		ID: "20231214000100",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(models.PaymentAccountDetail{}).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})
}
