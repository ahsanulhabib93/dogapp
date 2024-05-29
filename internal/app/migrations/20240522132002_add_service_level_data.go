package migrations

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
	"runtime"

	"github.com/jinzhu/gorm"
	migrator "github.com/voonik/goFramework/pkg/migrations"
	"github.com/voonik/ss2/internal/app/models"
	"gopkg.in/gormigrate.v1"
)

func init() {

	migrator.Register(&gormigrate.Migration{
		ID: "20240522132002",
		Migrate: func(tx *gorm.DB) (err error) {
			_, path, _, _ := runtime.Caller(0)
			filename := filepath.Join(filepath.Dir(path), "/20240522132002_partner_service_level.json")

			serviceLevels := []models.PartnerServiceLevel{}
			serviceLevelsJson, _ := ioutil.ReadFile(filename)
			err = json.Unmarshal(serviceLevelsJson, &serviceLevels)

			if err != nil {
				log.Println("Error while Unmarshal", err.Error())
				return err
			}

			for _, serviceLevel := range serviceLevels {
				err = tx.Create(&serviceLevel).Error
				if err != nil {
					log.Println("Error while creating service level", err.Error())
					break
				}
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})

	migrator.Register(&gormigrate.Migration{
		ID: "20240527045457",
		Migrate: func(tx *gorm.DB) (err error) {
			tx.Exec("update partner_service_mappings set partner_service_level_id = service_level;")
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})

}
