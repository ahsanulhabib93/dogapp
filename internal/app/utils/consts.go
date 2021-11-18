package utils

type SupplierType uint16
type AccountType uint16

const (
	L0 SupplierType = 1 + iota
	L1
	L2
	L3
	Hlc
)

const (
	Bank AccountType = 1 + iota
	Mfs
)
