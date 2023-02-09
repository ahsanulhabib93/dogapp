package utils

type SupplierType uint16
type AccountType uint16
type AccountSubType uint16

const (
	L0 SupplierType = 1 + iota
	L1
	L2
	L3
	Hlc
)

var SupplierTypeValue = map[SupplierType]string{
	L0:  "L0",
	L1:  "L1",
	L2:  "L2",
	L3:  "L3",
	Hlc: "Hlc",
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
