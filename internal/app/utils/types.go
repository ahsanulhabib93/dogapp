package utils

var AllowedUploadType = map[string][]string{
	// key - name - extension
	"SupplierShopImage":     {"shop_images", "", "jpg"},
	"SupplierNIDFrontImage": {"nid_front_images", "", "jpg"},
	"SupplierNIDBackImage":  {"nid_back_images", "", "jpg"},
	"SupplierTradeLicense":  {"trade_licenses", "", "pdf"},
	"SupplierAgreement":     {"agreements", "", "pdf"},
	"ShopOwnerImage":        {"shop_owner_images", "", "jpg"},
	"GuarantorImage":        {"guarantor_images", "", "jpg"},
	"ChequeImage":           {"cheque_images", "", "jpg"},
}

var SupplierDocumentType = []string{
	"nid_number",
	"nid_front_image_url",
	"nid_back_image_url",
	"trade_license_url",
	"agreement_url",
	"shop_owner_images_url",
	"guarantor_images_url",
	"cheque_images_url",
}
