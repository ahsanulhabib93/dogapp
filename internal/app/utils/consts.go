package utils

type ServiceType uint16
type SupplierType uint16
type AccountType uint16
type AccountSubType uint16
type VerificationStatus string
type BusinessType uint8
type ColorCode uint8
type SellerPriceVerified string
type StateReason uint64
type ActivationState uint64

const (
	Supplier ServiceType = 1 + iota
	Transporter
	RentVendor
	MwsOwner
	DoBuyer
	ProcurementVendor
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
	NOT_ACTIVATED ActivationState = 1
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
	Manufacturer BusinessType = 1 + iota
	Trader
)

var SellerBusinessType = map[BusinessType]string{
	Manufacturer: "MANUFACTURER",
	Trader:       "TRADER",
}

const (
	Platinum ColorCode = 1 + iota
	Gold
	Green
	Brown
	Black
)

var SellerColorCode = map[ColorCode]string{
	Platinum: "PLATINUM",
	Gold:     "GOLD",
	Green:    "GREEN",
	Brown:    "BROWN",
	Black:    "BLACK",
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
)

var SupplierTypeValue = map[SupplierType]string{
	L0:            "L0",
	L1:            "L1",
	L2:            "L2",
	L3:            "L3",
	Hlc:           "Hlc",
	Captive:       "Captive",
	Driver:        "Driver",
	CashVendor:    "Cash Vendor",
	RedxHubVendor: "Redx Hub Vendor",
	CreditVendor:  "Credit Vendor",
	HubRent:       "Hub Rent",
	WarehouseRent: "Warehouse Rent",
	DBHouseRent:   "DBHouse Rent",
	OfficeRent:    "Office Rent",
	Mws:           "Mws",
	Buyer:         "Buyer",
	Procurement:   "Procurement",
}

const (
	Bank AccountType = 1 + iota
	Mfs
)

const (
	Current AccountSubType = 1 + iota
	Savings
	Bkash
	Nagada
)

const (
	BucketFolder       = "ss2"
	SupplierAuditTopic = "cash_audit" // topic should be "audit_supplier". Renamed to incorporate SRE team
)

const (
	EmptyString = ""
	Zero        = 0
	One         = 1
	Ten         = 10
)

const ChangePendingSupplierStatus = "change_pending_supplier_status"
const ScheduleEveryDay = "0 0 * * *"

var PartnerServiceTypeMapping = map[string]ServiceType{
	"Supplier":          Supplier,
	"Transporter":       Transporter,
	"RentVendor":        RentVendor,
	"MwsOwner":          MwsOwner,
	"DoBuyer":           DoBuyer,
	"ProcurementVendor": ProcurementVendor,
}

var PartnerServiceLevelMapping = map[string]SupplierType{
	"L0":            L0,
	"L1":            L1,
	"L2":            L2,
	"L3":            L3,
	"Hlc":           Hlc,
	"Captive":       Captive,
	"Driver":        Driver,
	"CashVendor":    CashVendor,
	"RedxHubVendor": RedxHubVendor,
	"CreditVendor":  CreditVendor,
	"HubRent":       HubRent,
	"WarehouseRent": WarehouseRent,
	"DBHouseRent":   DBHouseRent,
	"OfficeRent":    OfficeRent,
	"Mws":           Mws,
	"Buyer":         Buyer,
	"Procurement":   Procurement,
}
