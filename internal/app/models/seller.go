package models

import (
	"time"

	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/utils"
)

type Seller struct {
	database.VaccountGorm
	UserID                       uint64
	AffiliateURL                 string
	IsDirect                     bool `gorm:"default:false"`
	Configuration                string
	DataMapping                  string
	FullfillmentType             int
	DeletedAt                    *time.Time
	BrandName                    string
	CompanyName                  string
	PrimaryEmail                 string `json:"primary_email"`
	PrimaryPhone                 string
	SupportEmail                 string
	SupportPhone                 string
	ActivationState              int `gorm:"default:1"`
	Slug                         string
	ReturnExchangePolicy         string
	TinNumber                    string
	PanNumber                    string `gorm:"default:'AAAAA0000A'"`
	MouAgreed                    bool   `gorm:"default:true"`
	TinDeclaration               string
	AgentID                      int
	CompanyType                  int
	SellerInvoiceNumber          int     `gorm:"default:0"`
	SellerType                   int     `gorm:"type:int;size:1;default:3"`
	SellerRating                 float64 `gorm:"type:decimal(5,2);default:0.0"`
	TanNumber                    string
	AadharCard                   string
	AggregatorID                 int    `gorm:"default:0"`
	InternationalEnabled         string `gorm:"default:'0'"`
	BusinessType                 utils.BusinessType
	ColorCode                    utils.ColorCode
	EmailConfirmed               bool `gorm:"default:false"`
	StateReason                  int
	GSTNumber                    string
	GSTStatus                    string `gorm:"default:'NOT_VERIFIED'"`
	GSTRelatedPanNumber          string
	GSTRelatedPanCardFileName    string
	GSTRelatedPanCardContentType string
	GSTRelatedPanCardFileSize    int
	GSTRelatedPanCardUpdatedAt   string
	Hub                          string
	OfficeTime                   string
	Slot                         string
	SellerCloseDay               string
	AcceptedPaymentMethods       string
	DeliveryType                 int `gorm:"default:1"`
	ProcessingType               int `gorm:"default:1"`
	BusinessUnit                 int
	SellerConfig                 *SellerConfig
	VendorAddresses              []*VendorAddress
	SellerBankDetail             *SellerBankDetail
	SellerPricingDetails         []*SellerPricingDetail
	SellerAccountManagers        []*SellerAccountManager
}
