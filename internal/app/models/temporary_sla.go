package models

import (
	"time"

	"github.com/voonik/goFramework/pkg/database"
)

type TemporarySla struct {
	database.VaccountGorm
	SellerID     int
	ImpactSla    int
	SlaStartDate *time.Time
	SlaEndDate   *time.Time
}
