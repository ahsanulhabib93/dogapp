package models

import (
	"time"

	"github.com/voonik/goFramework/pkg/database"
)

type SellerAccountManager struct {
	database.VaccountGorm
	SellerID  int
	Role      string
	DeletedAt *time.Time
	Priority  int
	Phone     int64
	Name      string
	Email     string
}
