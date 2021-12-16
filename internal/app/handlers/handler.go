package handlers

import (
	"github.com/voonik/ss2/internal/app/services"
)

// GetSupplierInstance ...
func GetSupplierInstance() *services.SupplierService {
	return new(services.SupplierService)
}

// GetSupplierAddressInstance ...
func GetSupplierAddressInstance() *services.SupplierAddressService {
	return new(services.SupplierAddressService)
}

// GetPaymentAccountDetailInstance ...
func GetPaymentAccountDetailInstance() *services.PaymentAccountDetailService {
	return new(services.PaymentAccountDetailService)
}

// GetKeyAccountManagerInstance ...
func GetKeyAccountManagerInstance() *services.KeyAccountManagerService {
	return new(services.KeyAccountManagerService)
}
