package models

import (
	"time"

	"github.com/voonik/goFramework/pkg/database"
)

type SellerAccountManager struct {
	database.VaccountGorm
	SellerID   int
	VaccountID int `gorm:"not null"`
	Role       string
	DeletedAt  time.Time
	Priority   int
	Phone      int64
	Name       string
	Email      string
}
