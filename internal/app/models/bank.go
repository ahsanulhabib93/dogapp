package models

import (
	"github.com/voonik/goFramework/pkg/database"
)

type Bank struct {
	database.VModel
	Name string `gorm:"not null" valid:"required"`
}

func GetBankJoinStr() string {
	return "left join banks on banks.id = payment_account_details.bank_id"
}
