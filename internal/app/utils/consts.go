package utils

import "database/sql/driver"

type ServiceType uint16
type SupplierType uint16
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

const (
	Supplier ServiceType = 1 + iota
	Transporter
	RentVendor
	MwsOwner
	DoBuyer
	ProcurementVendor
	Employee
	Vendor
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
	Success = "success"
	Failure = "failure"
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
	L0 SupplierType = 1 + iota
	L1
	L2
	L3
	Hlc
	Captive
	Driver
	CashVendor
	RedxHubVendor
	CreditVendor
	HubRent
	WarehouseRent
	DBHouseRent
	OfficeRent
	Mws
	Buyer
	Procurement
	InternalEmployee
	ExternalVendor
)

var SupplierTypeValue = map[SupplierType]string{
	L0:               "L0",
	L1:               "L1",
	L2:               "L2",
	L3:               "L3",
	Hlc:              "Hlc",
	Captive:          "Captive",
	Driver:           "Driver",
	CashVendor:       "Cash Vendor",
	RedxHubVendor:    "Redx Hub Vendor",
	CreditVendor:     "Credit Vendor",
	HubRent:          "Hub Rent",
	WarehouseRent:    "Warehouse Rent",
	DBHouseRent:      "DBHouse Rent",
	OfficeRent:       "Office Rent",
	Mws:              "Mws",
	Buyer:            "Buyer",
	Procurement:      "Procurement",
	InternalEmployee: "Internal Employee",
	ExternalVendor:   "External Vendor",
}

const (
	Bank AccountType = 1 + iota
	Mfs
	PrepaidCard
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
	Ten               = 10
	SixtyFour         = 64
	DefaultDateFormat = "2006-01-02"
)

const ChangePendingSupplierStatus = "change_pending_supplier_status"
const ScheduleEveryDay = "0 0 * * *"
const SS2UinquePrefixKey = "SS2-PAD-"

var PartnerServiceTypeMapping = map[string]ServiceType{
	"Supplier":          Supplier,
	"Transporter":       Transporter,
	"RentVendor":        RentVendor,
	"MwsOwner":          MwsOwner,
	"DoBuyer":           DoBuyer,
	"ProcurementVendor": ProcurementVendor,
	"Employee":          Employee,
	"Vendor":            Vendor,
}

var PartnerServiceLevelMapping = map[string]SupplierType{
	"L0":               L0,
	"L1":               L1,
	"L2":               L2,
	"L3":               L3,
	"Hlc":              Hlc,
	"Captive":          Captive,
	"Driver":           Driver,
	"CashVendor":       CashVendor,
	"RedxHubVendor":    RedxHubVendor,
	"CreditVendor":     CreditVendor,
	"HubRent":          HubRent,
	"WarehouseRent":    WarehouseRent,
	"DBHouseRent":      DBHouseRent,
	"OfficeRent":       OfficeRent,
	"Mws":              Mws,
	"Buyer":            Buyer,
	"Procurement":      Procurement,
	"InternalEmployee": InternalEmployee,
	"ExternalVendor":   ExternalVendor,
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
}
