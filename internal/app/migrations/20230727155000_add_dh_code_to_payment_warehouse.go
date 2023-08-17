package migrations

import (
	"github.com/jinzhu/gorm"
	migrator "github.com/voonik/goFramework/pkg/migrations"
	"gopkg.in/gormigrate.v1"
)

func init() {
	migrator.Register(&gormigrate.Migration{
		ID: "20230727155000",
		Migrate: func(tx *gorm.DB) error {
			type PaymentAccountDetailWarehouseMapping struct {
				DhCode string `json:"dh_code,omitempty"`
			}
			return tx.AutoMigrate(PaymentAccountDetailWarehouseMapping{}).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})
}
