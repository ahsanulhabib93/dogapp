package migrations

import (
	"github.com/jinzhu/gorm"
	migrator "github.com/voonik/goFramework/pkg/migrations"
	"gopkg.in/gormigrate.v1"
)

func init() {
	migrator.Register(&gormigrate.Migration{
		ID: "20220412000000",
		Migrate: func(tx *gorm.DB) error {
			return tx.Exec("ALTER TABLE supplier_addresses CHANGE COLUMN `zipcode` `zipcode` VARCHAR(255) NULL ;").Error
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})
}
