package models

import (
	"github.com/voonik/goFramework/pkg/database"
)

type SupplierAddress struct {
	database.VaccountGorm
	SupplierID uint64
	Firstname  string
	Lastname   string
	Address1   string
	Address2   string
	Landmark   string
	City       string
	State      string
	Country    string
	Zipcode    string
	Phone      string
	GstNumber  string `json:"gst_number"`
	IsDefault  bool
	Supplier   Supplier
}
