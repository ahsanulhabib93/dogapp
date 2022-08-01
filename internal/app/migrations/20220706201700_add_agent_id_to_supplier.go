package migrations

import (
	"github.com/jinzhu/gorm"
	migrator "github.com/voonik/goFramework/pkg/migrations"
	"gopkg.in/gormigrate.v1"
)

func init() {
	migrator.Register(&gormigrate.Migration{
		ID: "20220706201700",
		Migrate: func(tx *gorm.DB) error {
			type Supplier struct {
				AgentID *uint64 `json:"agent_id"`
			}

			return tx.AutoMigrate(Supplier{}).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})
}
