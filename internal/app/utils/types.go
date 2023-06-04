package utils

var AllowedUploadType = map[string][]string{
	// key - name - extension
	"SupplierShopImage":      {"shop_images", "", "jpg"},
	"SupplierNIDFrontImage":  {"nid_front_images", "", "jpg"},
	"SupplierNIDBackImage":   {"nid_back_images", "", "jpg"},
	"SupplierTradeLicense":   {"trade_licenses", "", "pdf"},
	"SupplierAgreement":      {"agreements", "", "pdf"},
	"ShopOwnerImage":         {"shop_owner_images", "", "jpg"},
	"GuarantorImage":         {"guarantor_images", "", "jpg"},
	"GuarantorNIDFrontImage": {"guarantor_nid_front_images", "", "jpg"},
	"GuarantorNIDBackImage":  {"guarantor_nid_back_images", "", "jpg"},
	"ChequeImage":            {"cheque_images", "", "jpg"},
}

var SupplierPrimaryDocumentType = []string{
	"nid_number",
	"nid_front_image_url",
	"nid_back_image_url",
	"trade_license_url",
	"agreement_url",
}

var SupplierSecondaryDocumentType = []string{
	"shop_image_url",
	"shop_owner_image_url",
	"guarantor_image_url",
	"guarantor_nid_number",
	"guarantor_nid_front_image_url",
	"guarantor_nid_back_image_url",
	"cheque_image_url",
}

var PartnerServiceTypeLevelMapping = map[ServiceType][]SupplierType{
	// service type - service level
	Supplier:    {L0, L1, L2, L3, Hlc},
	Transporter: {Captive, Driver, CashVendor},
}
