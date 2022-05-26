package utils

var AllowedUploadType = map[string][]string{
	// key - name - extension
	"SupplierShopImage":     {"shop_images", "", "png"},
	"SupplierNIDFrontImage": {"nid_front_images", "", "jpg"},
	"SupplierNIDBackImage":  {"nid_back_images", "", "jpg"},
	"SupplierTradeLicense":  {"trade_licenses", "", "pdf"},
	"SupplierAgreement":     {"agreements", "", "pdf"},
}
