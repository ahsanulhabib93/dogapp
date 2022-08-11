package models

import (
	"github.com/voonik/goFramework/pkg/database"
)

type PaymentAccountDetailWarehouseMapping struct {
	database.VaccountGorm
	PaymentAccountDetailID uint64 `gorm:"index:idx_payment_detail_id"`
	WarehouseID            uint64
}
