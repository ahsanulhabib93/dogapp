package utils

type ServiceType uint16
type SupplierType uint16
type AccountType uint16
type AccountSubType uint16
type VerificationStatus string
type BusinessType string
type ColorCode string
type SellerPriceVerified string

const (
	Supplier ServiceType = 1 + iota
	Transporter
	RentVendor
	MwsOwner
	DoBuyer
	ProcurementVendor
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

const (
	Platinum ColorCode = "PLATINUM"
	Gold     ColorCode = "GOLD"
	Green    ColorCode = "GREEN"
	Brown    ColorCode = "BROWN"
	Black    ColorCode = "BLACK"
)

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
