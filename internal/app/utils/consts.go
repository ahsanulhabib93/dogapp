package utils

import (
	"database/sql/driver"
	"strings"
)

type ServiceType uint16
type AccountType uint16
type AccountSubType uint16
type VerificationStatus string
type BusinessType string
type ColorCode string
type SellerPriceVerified string
type StateReason int
type ActivationState int
type AttachableType uint64
type FileType uint64
type BusinessUnit uint64

const VoonikFulfilmentType = 2
const SellerFetchSuccessMessage = "fetched seller details successfully"
const SourcingAssociateRole = "sourcing_associate"

const (
	Supplier ServiceType = 1 + iota
	Transporter
	RentVendor
	MwsOwner
	DoBuyer
	ProcurementVendor
	Employee
	Vendor
	OfficeSuppliesAndConsumables
	LogisticsAndTransportation
	CleaningAndMaintenance
	ProfessionalServices
	UtilitiesAndCommunication
	EmployeeAndClientRelations
	Insurance
)

const (
	POLICY_VIOLATION StateReason = 1 + iota
	FULFILMENT_VIOLATION
	PERFORMANCE
	PRODUCT_QUALITY
	PENDING_RISK_REVIEW
	PENDING_CONTACT_WITH_SS
	VACATION_MODE
)

const (
	NOT_ACTIVATED ActivationState = 1 + iota
	DOWNLOADED_EXCEL
	ACTIVATED
	VERIFICATION_PENDING
	HOLD_OFF
	PRE_ACTIVATED
	PENDING_QUALITY
	POLICY_VIOLATED
	FULFILMENT_VIOLATED
	VACATION_PENDING
	SUSPENDED
	BLOCKED
	FRAUD
	ON_HOLD
	GST_PENDING
	UNDER_REVIEW
)

const (
	Verified    VerificationStatus = "VERIFIED"
	Rejected    VerificationStatus = "REJECTED"
	NotVerified VerificationStatus = "NOT_VERIFIED"
)

const (
	Success                      = "success"
	Failure                      = "failure"
	DefaultSellerMaxQuantity     = uint64(1000)
	DefaultSellerItemsPerPackage = uint64(1)
	DefaultSellerPickupType      = uint64(1)
	DefaultCountry               = "Bangladesh"
	DefaultState                 = "Dhaka"
	Required                     = "required"
)

const (
	Manufacturer BusinessType = "MANUFACTURER"
	Trader       BusinessType = "TRADER"
)

func (bt *BusinessType) Scan(value interface{}) error {
	*bt = BusinessType(value.([]byte))
	return nil
}

func (bt BusinessType) Value() (driver.Value, error) {
	return string(bt), nil
}

const (
	Platinum ColorCode = "PLATINUM"
	Gold     ColorCode = "GOLD"
	Green    ColorCode = "GREEN"
	Brown    ColorCode = "BROWN"
	Black    ColorCode = "BLACK"
)

func (cd *ColorCode) Scan(value interface{}) error {
	*cd = ColorCode(value.([]byte))
	return nil
}

func (cd ColorCode) Value() (driver.Value, error) {
	return string(cd), nil
}

const (
	PriceNotVerified  SellerPriceVerified = "NOT_VERIFIED"
	PriceVerified     SellerPriceVerified = "VERIFIED"
	PriceAutoVerified SellerPriceVerified = "AUTO_VERIFIED"
)

const (
	Bank AccountType = 1 + iota
	Mfs
	PrepaidCard
	Cheque
)

const (
	Current AccountSubType = 1 + iota
	Savings
	Bkash
	Nagada
	EBL
	UCBL
)

const (
	BucketFolder       = "ss2"
	SupplierAuditTopic = "cash_audit" // topic should be "audit_supplier". Renamed to incorporate SRE team
)

const (
	EmptyString       = ""
	Zero              = 0
	One               = 1
	Three             = 3
	Ten               = 10
	SixtyFour         = 64
	DefaultDateFormat = "2006-01-02"
	Params            = "params"
)

const ChangePendingSupplierStatus = "change_pending_supplier_status"
const CreateOMSSellerSync = "create_oms_seller_sync"
const ScheduleEveryDay = "0 0 * * *"
const SS2UinquePrefixKey = "SS2-PAD-"

var PartnerServiceTypeMapping = map[string]ServiceType{
	"Supplier":                     Supplier,
	"Transporter":                  Transporter,
	"RentVendor":                   RentVendor,
	"MwsOwner":                     MwsOwner,
	"DoBuyer":                      DoBuyer,
	"ProcurementVendor":            ProcurementVendor,
	"Employee":                     Employee,
	"Vendor":                       Vendor,
	"OfficeSuppliesAndConsumables": OfficeSuppliesAndConsumables,
	"LogisticsAndTransportation":   LogisticsAndTransportation,
	"CleaningAndMaintenance":       CleaningAndMaintenance,
	"ProfessionalServices":         ProfessionalServices,
	"UtilitiesAndCommunication":    UtilitiesAndCommunication,
	"EmployeeAndClientRelations":   EmployeeAndClientRelations,
	"Insurance":                    Insurance,
}

const (
	AttachableTypeSupplier AttachableType = iota + 1
	AttachableTypePartnerServiceMapping
)

const (
	SecurityCheque FileType = iota + 1
	GuarantorNID
	TIN
	BIN
	IncorporationCertificate
	TradeLicense
	PartnershipDeed
	EngagementLetter
	ConfirmationLetter
	AcknowledgementLetter
	MerchantPhoto
	GuarantorPhoto
	SanctionLetter
	Undertaking
	BankStatement
	RentalDeed
	ForwardingLetter
	DealershipDeed
	CIBReport
	ExitLetter
	SettlementAgreement
	NOC
	TaxAcknowledgement
	SalesStatement
	CompanyProfile
	KYC
	Invoice
	UtilityBill
)

var FileTypeMapping = map[string]FileType{
	"SecurityCheque":           SecurityCheque,
	"GuarantorNID":             GuarantorNID,
	"TIN":                      TIN,
	"BIN":                      BIN,
	"IncorporationCertificate": IncorporationCertificate,
	"TradeLicense":             TradeLicense,
	"PartnershipDeed":          PartnershipDeed,
	"EngagementLetter":         EngagementLetter,
	"ConfirmationLetter":       ConfirmationLetter,
	"AcknowledgementLetter":    AcknowledgementLetter,
	"MerchantPhoto":            MerchantPhoto,
	"GuarantorPhoto":           GuarantorPhoto,
	"SanctionLetter":           SanctionLetter,
	"Undertaking":              Undertaking,
	"BankStatement":            BankStatement,
	"RentalDeed":               RentalDeed,
	"ForwardingLetter":         ForwardingLetter,
	"DealershipDeed":           DealershipDeed,
	"CIBReport":                CIBReport,
	"ExitLetter":               ExitLetter,
	"SettlementAgreement":      SettlementAgreement,
	"NOC":                      NOC,
	"TaxAcknowledgement":       TaxAcknowledgement,
	"SalesStatement":           SalesStatement,
	"CompanyProfile":           CompanyProfile,
	"KYC":                      KYC,
	"Invoice":                  Invoice,
	"UtilityBill":              UtilityBill,
}

var SellerDataMapping = map[string]interface{}{
	"title":            "{{item_name_product_name__title_}}",
	"description":      "{{product_description}}",
	"id":               "{{product_id}}",
	"update_delete":    "{{update_delete}}",
	"size":             "{{size}}",
	"size_sku":         "{{sku_id}}",
	"quantity":         "{{quantity}}",
	"original_price":   "{{standard_price__mrp_}}",
	"price":            "{{sale_price__sp_}}",
	"product_category": "{{product_sub_type__category_}}",
	"image":            "{{main_image_url}}",
	"image1":           "{{secondary_image_url1}}",
	"image2":           "{{secondary_image_url2}}",
	"image3":           "{{secondary_image_url3}}",
	"image4":           "{{secondary_image_url4}}",
	"image5":           "{{secondary_image_url5}}",
	"image6":           "{{secondary_image_url5}}",
	"image7":           "{{secondary_image_url5}}",
	"colour":           "{{colour}}",
	"brand":            "{{brand_name}}",
	"color":            "{{color}}",
	"fabric":           "{{fabric_material}}",
	"lead_time":        "{{lead_time}}",
	"piercing":         "{{piercing_required}}",
	"occassion":        "{{occassion}}",
	"type":             "{{type}}",
	"cut":              "{{cut}}",
	"neck_type":        "{{neck_type}}",
	"sleeve_type":      "{{sleeve_type}}",
	"quality":          "{{quality}}",
	"hem":              "{{hemlines}}",
	"weight":           "{{weight}}",
	"length":           "{{length}}",
	"width":            "{{width}}",
	"wash_care":        "{{wash_care}}",
	"mapping": map[string]interface{}{
		"category_mapping": map[string]interface{}{},
	},
}

const (
	UNICORN                   BusinessUnit = 1
	UNBRANDED                 BusinessUnit = 2
	BRANDED                   BusinessUnit = 3
	LIFESTYLE                 BusinessUnit = 4
	BLITZ                     BusinessUnit = 5
	AGRO                      BusinessUnit = 6
	WHOLESALE                 BusinessUnit = 7
	MWS                       BusinessUnit = 8
	FRESH                     BusinessUnit = 9
	POULTRY                   BusinessUnit = 10
	APPAREL                   BusinessUnit = 11
	INFRA                     BusinessUnit = 12
	ENERGY                    BusinessUnit = 13
	ELECTRONICS               BusinessUnit = 14
	MOKAM_X                   BusinessUnit = 20
	AGRO_FISH_PROJECT         BusinessUnit = 21
	POP_BOISHAKHI             BusinessUnit = 22
	REDX_FULFILLMENT_SERVICE  BusinessUnit = 100
	REDX_FULFILLMENT_SERVICE2 BusinessUnit = 101
	DAS                       BusinessUnit = 200
)

func IsValidBusinessUnit(bu BusinessUnit) bool {
	return !strings.Contains(bu.String(), "BusinessUnit")
}

func IsValidActivationState(as ActivationState) bool {
	return !strings.Contains(as.String(), "ActivationState")
}

func (pt BusinessUnit) ID() uint16 {
	return uint16(pt)
}

func IsValidColorCode(cc ColorCode) bool {
	switch cc {
	case Platinum, Gold, Green, Brown, Black:
		return true
	default:
		return false
	}
}
