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

// GetPartnerServiceMapping ...
func GetPartnerServiceMappingInstance() *services.PartnerServiceMappingService {
	return new(services.PartnerServiceMappingService)
}

func GetSellerInstance() *services.SellerService {
	return new(services.SellerService)
}

func GetSellerBankDetailInstance() *services.SellerBankDetailService {
	return new(services.SellerBankDetailService)
}

func GetSellerPricingDetailInstance() *services.SellerPricingDetailService {
	return new(services.SellerPricingDetailService)
}

func GetVendorAddressInstance() *services.VendorAddressService {
	return new(services.VendorAddressService)
}

func GetAttachmentInstance() *services.AttachmentService {
	return new(services.AttachmentService)
}
