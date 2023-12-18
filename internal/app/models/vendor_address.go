package models

import (
	"time"

	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/utils"
)

type VendorAddress struct {
	database.VaccountGorm
	Firstname                    string `gorm:"column:firstname"`
	Lastname                     string `gorm:"column:lastname"`
	Address1                     string `gorm:"column:address1"`
	Address2                     string `gorm:"column:address2"`
	City                         string `gorm:"column:city"`
	Zipcode                      string `gorm:"column:zipcode"`
	AlternativePhone             string `gorm:"column:alternate_phone"`
	Company                      string `gorm:"column:company"`
	State                        string `gorm:"column:state"`
	Country                      string `gorm:"column:country"`
	AddressType                  int
	SellerID                     int `gorm:"column:seller_id"`
	LandMark                     string
	DefaultAddress               bool `gorm:"default:false"`
	DeletedAt                    *time.Time
	Phone                        int64
	AddressProofFileName         string
	AddressProofContentType      string
	AddressProofFileSize         int
	AddressProofUpdatedAt        string
	VerificationStatus           utils.VerificationStatus `gorm:"default:'NOT_VERIFIED'"`
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
	ExtraData                    string `gorm:"default:'{}';column:extra_detail"`
}
