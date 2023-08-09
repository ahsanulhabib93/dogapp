package utils

type ServiceType uint16
type SupplierType uint16
type AccountType uint16
type AccountSubType uint16

const (
	Supplier ServiceType = 1 + iota
	Transporter
	RentVendor
	Mws
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
	MwsOwner
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
	MwsOwner:      "MwsOwner",
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
)

const ChangePendingSupplierStatus = "change_pending_supplier_status"
const ScheduleEveryDay = "0 0 * * *"

var PartnerServiceTypeMapping = map[string]ServiceType{
	"Supplier":    Supplier,
	"Transporter": Transporter,
	"RentVendor":  RentVendor,
	"Mws":         Mws,
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
	"MwsOwner":      MwsOwner,
}
