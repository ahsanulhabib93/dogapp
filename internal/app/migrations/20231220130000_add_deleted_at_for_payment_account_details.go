package migrations

import (
	"github.com/jinzhu/gorm"
	migrator "github.com/voonik/goFramework/pkg/migrations"
	"gopkg.in/gormigrate.v1"
	gormIO "gorm.io/gorm"
)

func init() {
	migrator.Register(&gormigrate.Migration{
		ID: "20231220130000",
		Migrate: func(tx *gorm.DB) error {
			type PaymentAccountDetail struct {
				DeletedAt gormIO.DeletedAt `json:"deleted_at,omitempty"`
			}
			return tx.AutoMigrate(PaymentAccountDetail{}).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})
}
