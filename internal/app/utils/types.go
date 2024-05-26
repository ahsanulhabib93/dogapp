package utils

var AllowedUploadType = map[string][]string{
	// key - name - extension
	"SupplierShopImage":        {"shop_images", "", "jpg"},
	"SupplierNIDFrontImage":    {"nid_front_images", "", "jpg"},
	"SupplierNIDBackImage":     {"nid_back_images", "", "jpg"},
	"SupplierTradeLicense":     {"trade_licenses", "", "pdf"},
	"SupplierAgreement":        {"agreements", "", "pdf"},
	"ShopOwnerImage":           {"shop_owner_images", "", "jpg"},
	"GuarantorImage":           {"guarantor_images", "", "jpg"},
	"GuarantorNIDFrontImage":   {"guarantor_nid_front_images", "", "jpg"},
	"GuarantorNIDBackImage":    {"guarantor_nid_back_images", "", "jpg"},
	"ChequeImage":              {"cheque_images", "", "jpg"},
	"SecurityCheque":           {"security_cheque", "", "jpg"},
	"GuarantorNID":             {"guarantor_nid", "", "jpg"},
	"TIN":                      {"tin", "", "jpg"},
	"BIN":                      {"bin", "", "jpg"},
	"IncorporationCertificate": {"incorporation_certificate", "", "jpg"},
	"TradeLicense":             {"trade_license", "", "jpg"},
	"PartnershipDeed":          {"partnership_deed", "", "jpg"},
	"EngagementLetter":         {"engagement_letter", "", "jpg"},
	"ConfirmationLetter":       {"confirmation_letter", "", "jpg"},
	"AcknowledgementLetter":    {"acknowledgement_letter", "", "jpg"},
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

var AttachableFileTypeMapping = map[AttachableType][]FileType{
	AttachableTypeSupplier:              {SecurityCheque, GuarantorNID, TIN, BIN, IncorporationCertificate, TradeLicense, PartnershipDeed, EngagementLetter, ConfirmationLetter, AcknowledgementLetter},
	AttachableTypePartnerServiceMapping: {},
}
