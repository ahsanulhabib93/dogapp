package models

import (
	"time"

	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/utils"
)

type VendorAddress struct {
	database.VaccountGorm
	FirstName                    string
	LastName                     string
	Address1                     string
	Address2                     string
	City                         string
	ZipCode                      string
	AlternativePhone             string
	Company                      string
	State                        string
	Country                      string
	AddressType                  int
	SellerID                     int
	LandMark                     string
	DefaultAddress               bool `gorm:"default:false"`
	DeletedAt                    time.Time
	Phone                        int64
	AddressProofFileName         string
	AddressProofContentType      string
	AddressProofFileSize         int
	AddressProofUpdatedAt        string
	VerificationStatus           utils.VerificationStatus `gorm:"default:'NOT_VERIFIED'"`
	VaccountID                   int                      `gorm:"not null"`
	UUID                         string                   `gorm:"type:varchar(50)"`
	GSTNumber                    string
	GSTStatus                    string `gorm:"default:'NOT_VERIFIED'"`
	GSTCardFileName              string
	GSTCardContentType           string
	GSTCardFileSize              int
	GSTCardUpdatedAt             string
	GSTRelatedPanNumber          string
	GSTRelatedPanCardFileName    string
	GSTRelatedPanCardContentType string
	GSTRelatedPanCardFileSize    int
	GSTRelatedPanCardUpdatedAt   string
}
