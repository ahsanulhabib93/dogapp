package models

import (
	"time"

	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/utils"
)

type SellerBankDetail struct {
	database.VaccountGorm
	AccountNumber              string
	AccountHolderName          string
	IFSCCode                   string
	AccountType                int
	BankName                   string
	BankBranch                 string
	DeletedAt                  *time.Time
	PanCardFileName            string
	PanCardContentType         string
	PanCardFileSize            int
	PanCardUpdatedAt           *time.Time
	CancelledChequeFileName    string
	CancelledChequeContentType string
	CancelledChequeFileSize    int
	CancelledChequeUpdatedAt   *time.Time
	AgreementCopyFileName      string
	AgreementCopyContentType   string
	AgreementCopyFileSize      int
	AgreementCopyUpdatedAt     *time.Time
	SellerID                   int
	TinCardFileName            string
	TinCardContentType         string
	TinCardFileSize            int
	TinCardUpdatedAt           *time.Time
	TanCardFileName            string
	TanCardContentType         string
	TanCardFileSize            int
	TanCardUpdatedAt           *time.Time
	AadharCardFileName         string
	AadharCardContentType      string
	AadharCardFileSize         int
	AadharCardUpdatedAt        *time.Time
	VerificationStatus         utils.VerificationStatus `gorm:"default:'NOT_VERIFIED'"`
	GSTCardFileName            string
	GSTCardContentType         string
	GSTCardFileSize            int
	GSTCardUpdatedAt           *time.Time
}
