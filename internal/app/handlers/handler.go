package handlers

import (
	"github.com/voonik/ss2/internal/app/services"
)

func GetSupplierInstance() *services.SupplierService {
	return new(services.SupplierService)
}

func GetSupplierAddressInstance() *services.SupplierAddressService {
	return new(services.SupplierAddressService)
}

func GetPaymentAccountDetailInstance() *services.PaymentAccountDetailService {
	return new(services.PaymentAccountDetailService)
}

func GetKeyAccountManagerInstance() *services.KeyAccountManagerService {
	return new(services.KeyAccountManagerService)
}
