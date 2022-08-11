package models

import (
	"github.com/voonik/goFramework/pkg/database"
)

type Bank struct {
	database.VModel
	Name string `gorm:"not null" valid:"required"`
}

func GetBankJoinStr() string {
	return "LEFT JOIN banks ON banks.id = payment_account_details.bank_id"
}
