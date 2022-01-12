package models

import (
	"github.com/voonik/goFramework/pkg/database"
)

type Bank struct {
	database.VaccountGorm
	Name string `gorm:"not null" valid:"required"`
}
