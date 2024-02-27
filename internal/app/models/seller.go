package models

import (
	"time"

	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/utils"
	"gorm.io/datatypes"
)

type Seller struct {
	database.VaccountGorm
	UserID                       uint64 `json:"user_id"`
	AffiliateURL                 string
	IsDirect                     bool `gorm:"default:false"`
	Configuration                string
	DataMapping                  datatypes.JSON `gorm:"type:json"`
	FullfillmentType             int
	DeletedAt                    *time.Time
	BrandName                    string `json:"brand_name"`
	CompanyName                  string `json:"company_name"`
	PrimaryEmail                 string `json:"primary_email"`
	PrimaryPhone                 string `json:"primary_phone"`
	SupportEmail                 string
	SupportPhone                 string
	ActivationState              utils.ActivationState `gorm:"default:1" json:"activation_state"`
	Slug                         string                `json:"slug"`
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
	AggregatorID                 int                `gorm:"default:0"`
	InternationalEnabled         string             `gorm:"default:'0'"`
	BusinessType                 utils.BusinessType `gorm:"column:business_type;type:enum('MANUFACTURER', 'TRADER');default:NULL"`
	ColorCode                    utils.ColorCode    `gorm:"column:color_code;type:enum('PLATINUM','GOLD','GREEN','BROWN','BLACK');default:NULL"`
	EmailConfirmed               bool               `gorm:"default:false"`
	StateReason                  utils.StateReason
	GSTNumber                    string
	GSTStatus                    string `gorm:"default:'NOT_VERIFIED'"`
	GSTRelatedPanNumber          string
	GSTRelatedPanCardFileName    string
	GSTRelatedPanCardContentType string
	GSTRelatedPanCardFileSize    int
	GSTRelatedPanCardUpdatedAt   string
	Hub                          string `json:"hub"`
	OfficeTime                   string
	Slot                         string `json:"slot"`
	SellerCloseDay               string
	AcceptedPaymentMethods       string
	DeliveryType                 int            `gorm:"default:1" json:"delivery_type"`
	ProcessingType               int            `gorm:"default:1" json:"processing_type"`
	BusinessUnit                 int            `json:"business_unit"`
	ReturnExchangePolicy         datatypes.JSON `gorm:"type:json"`
	SellerConfig                 *SellerConfig
	VendorAddresses              []*VendorAddress
	SellerBankDetail             *SellerBankDetail
	SellerPricingDetails         []*SellerPricingDetail
	SellerAccountManagers        []*SellerAccountManager
}
