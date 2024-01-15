package models

import (
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/utils"
	"time"
)

type Attachment struct {
	database.VaccountGorm
	AttachableID    uint64 `gorm:"not null"`
	FileURL         string `gorm:"type:varchar(512); not null"`
	Parent          utils.AttachableType
	ReferenceNumber string `gorm:"type:varchar(255); not null"`
	Type            utils.FileType
	DeletedAt       *time.Time `gorm:"index"`
}
