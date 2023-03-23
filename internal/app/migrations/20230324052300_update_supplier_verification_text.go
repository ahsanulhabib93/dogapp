package migrations

import (
	"time"

	"github.com/jinzhu/gorm"
	aaaModels "github.com/voonik/goFramework/pkg/aaa/models"
	migrator "github.com/voonik/goFramework/pkg/migrations"
	"gopkg.in/gormigrate.v1"
)

func init() {
	migrator.Register(&gormigrate.Migration{
		ID: "20230324052300",
		Migrate: func(tx *gorm.DB) error {
			tableName := tx.NewScope(aaaModels.AppPreference{}).TableName()
			updatedText := "প্রিয় ব্যবসায়িক পার্টনার, আপনার সাপ্লাইয়ার ভেরিফিকেশন কোডটি হলো - $otp। ভেরিফিকেশন কোডটি মোকাম পার্টনারকে প্রদান করে রেজিস্ট্রেশন সম্পূর্ন করুন" //nolint:revive
			appPref := aaaModels.AppPreference{}
			resp := tx.Raw("SELECT * FROM "+tableName+" WHERE preference_key = ?", "supplier_phone_verification_otp_content").Scan(&appPref)
			if resp.RecordNotFound() {
				query := `
					INSERT INTO ` + tableName + ` (preference_key,value,value_type,vaccount_id,portal_id,created_at,updated_at)
					VALUES ('supplier_phone_verification_otp_content', ?, 'string', 2, 2, ?, ?)
				`
				if err := tx.Exec(query, updatedText, time.Now(), time.Now()).Error; err != nil {
					return err
				}
			}

			return tx.Exec("update "+tableName+" set value = ? where preference_key='supplier_phone_verification_otp_content'", updatedText).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return nil
		},
	})
}
