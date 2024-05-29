package migrations

import (
	"github.com/jinzhu/gorm"
	migrator "github.com/voonik/goFramework/pkg/migrations"
	"github.com/voonik/ss2/internal/app/models"
	"gopkg.in/gormigrate.v1"
)

// keep only one migration
func init() {
	migrator.Register(&gormigrate.Migration{
		ID: "20230523113708",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(
				models.PartnerServiceMapping{},
			).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})
	migrator.Register(&gormigrate.Migration{
		ID: "20230523113716",
		Migrate: func(tx *gorm.DB) error {
			return tx.Model(
				models.PartnerServiceMapping{},
			).AddUniqueIndex("idx_partner_service", "supplier_id", "service_type", "vaccount_id").Error
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})
}
