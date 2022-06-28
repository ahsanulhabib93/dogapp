package models

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"

	aaaModels "github.com/voonik/goFramework/pkg/aaa/models"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/utils"
)

// SupplierStatusConstants ...
type SupplierStatus string

const (
	SupplierStatusPending  SupplierStatus = "Pending"
	SupplierStatusVerified SupplierStatus = "Verified"
	SupplierStatusFailed   SupplierStatus = "Failed"
	SupplierStatusBlocked  SupplierStatus = "Blocked"
)

// Supplier ...
type Supplier struct {
	database.VaccountGorm
	Name                      string         `gorm:"not null" valid:"required"`
	Status                    SupplierStatus `gorm:"default:'Pending'"`
	Reason                    string
	Email                     string
	Phone                     string
	AlternatePhone            string                 `json:"alternate_phone"`
	BusinessName              string                 `json:"business_name"`
	IsPhoneVerified           *bool                  `gorm:"default:false" json:"is_phone_verified"` // using pointer to update false value in Edit API
	ShopImageURL              string                 `json:"shop_image_url"`
	UserID                    *uint64                `json:"user_id"`
	NidNumber                 string                 `json:"nid_number"`
	NidFrontImageUrl          string                 `gorm:"type:varchar(512)" json:"nid_front_image_url"`
	NidBackImageUrl           string                 `gorm:"type:varchar(512)" json:"nid_back_image_url"`
	TradeLicenseUrl           string                 `gorm:"type:varchar(512)" json:"trade_license_url"`
	AgreementUrl              string                 `gorm:"type:varchar(512)" json:"agreement_url"`
	ShopOwnerImageUrl         string                 `gorm:"type:varchar(512)" json:"shop_owner_image_url"`
	GuarantorImageUrl         string                 `gorm:"type:varchar(512)" json:"guarantor_image_url"`
	GuarantorNidNumber        string                 `json:"guarantor_nid_number"`
	GuarantorNidFrontImageUrl string                 `gorm:"type:varchar(512)" json:"guarantor_nid_front_image_url"`
	GuarantorNidBackImageUrl  string                 `gorm:"type:varchar(512)" json:"guarantor_nid_back_image_url"`
	ChequeImageUrl            string                 `gorm:"type:varchar(512)" json:"cheque_image_url"`
	SupplierType              utils.SupplierType     `json:"supplier_type" valid:"required"`
	SupplierAddresses         []SupplierAddress      `json:"supplier_addresses"`
	PaymentAccountDetails     []PaymentAccountDetail `json:"payment_account_details"`
	KeyAccountManagers        []KeyAccountManager
	SupplierCategoryMappings  []SupplierCategoryMapping
	SupplierOpcMappings       []SupplierOpcMapping
}

// Validate ...
func (supplier *Supplier) Validate(db *gorm.DB) {
	result := db.Model(&supplier).First(&Supplier{}, "id != ? and phone = ?", supplier.ID, supplier.Phone)
	if !result.RecordNotFound() {
		db.AddError(errors.New("Supplier Already Exists"))
	}

	isNIDInvalid := false
	for _, c := range supplier.NidNumber {
		isNIDInvalid = isNIDInvalid || !('0' <= c && c <= '9')
	}

	for _, c := range supplier.GuarantorNidNumber {
		isNIDInvalid = isNIDInvalid || !('0' <= c && c <= '9')
	}

	if isNIDInvalid {
		db.AddError(errors.New("NID number should only consist of digits"))
	}

	if phoneNumber := strings.TrimSpace(supplier.Phone); len(phoneNumber) == 0 {
		db.AddError(errors.New("Phone Number can't be blank"))
	} else if !((strings.HasPrefix(phoneNumber, "8801") && len(phoneNumber) == 13) ||
		(strings.HasPrefix(phoneNumber, "01") && len(phoneNumber) == 11) ||
		(strings.HasPrefix(phoneNumber, "1") && len(phoneNumber) == 10)) {
		db.AddError(errors.New("Invalid Phone Number"))
	}
}

func (supplier *Supplier) IsOTPVerified() bool {
	if supplier.IsPhoneVerified == nil {
		return false
	}

	return *supplier.IsPhoneVerified
}

func (supplier *Supplier) Verify(ctx context.Context) error {
	paymentAccountsCount := database.DBAPM(ctx).Model(supplier).Association("PaymentAccountDetails").Count()
	if paymentAccountsCount == 0 {
		return errors.New("At least one payment account details should be present")
	}

	addressesCount := database.DBAPM(ctx).Model(supplier).Association("SupplierAddresses").Count()
	if addressesCount == 0 {
		return errors.New("At least one supplier address should be present")
	}

	if !(supplier.IsOTPVerified() || supplier.IsAnyDocumentPresent()) {
		return errors.New("At least one primary document or OTP verification needed")
	}

	typeValue := utils.SupplierTypeValue[supplier.SupplierType]
	otpTypeVerificationList := aaaModels.GetAppPreferenceServiceInstance().GetValue(ctx, "enabled_otp_verification", []string{}).([]string)
	if utils.IsInclude(otpTypeVerificationList, typeValue) && !supplier.IsOTPVerified() {
		msg := fmt.Sprint("OTP verification required for supplier type: ", typeValue)
		return errors.New(msg)
	}

	docTypeVerificationList := aaaModels.GetAppPreferenceServiceInstance().GetValue(ctx, "enabled_primary_doc_verification", []string{}).([]string)
	if utils.IsInclude(docTypeVerificationList, typeValue) && !supplier.IsAnyDocumentPresent() {
		msg := fmt.Sprint("At least one primary document required for supplier type: ", typeValue)
		return errors.New(msg)
	}

	return nil
}

func (supplier *Supplier) IsAnyDocumentPresent() bool {
	return !(supplier.NidNumber == "" && supplier.NidFrontImageUrl == "" && supplier.NidBackImageUrl == "" &&
		supplier.TradeLicenseUrl == "" && supplier.AgreementUrl == "")
}

func (supplier *Supplier) IsChangeAllowed(ctx context.Context) bool {
	status := supplier.Status
	if !(status == SupplierStatusVerified || status == SupplierStatusBlocked) {
		return true
	}

	allowedPermission := aaaModels.GetAppPreferenceServiceInstance().GetValue(ctx, "supplier_update_allowed_permission", "supplierpanel:editverifiedblockedsupplieronly:admin").(string)
	permissions := utils.GetCurrentUserPermissions(ctx)
	for _, v := range permissions {
		if utils.IsInclude(strings.Split(v, " "), allowedPermission) {
			return true
		}
	}

	return false
}

// GetCategoryMappingJoinStr ...
func GetCategoryMappingJoinStr() string {
	return "LEFT JOIN supplier_category_mappings on supplier_category_mappings.supplier_id = suppliers.id and supplier_category_mappings.deleted_at IS NULL and supplier_category_mappings.vaccount_id = suppliers.vaccount_id"
}

// GetOpcMappingJoinStr ...
func GetOpcMappingJoinStr() string {
	return "LEFT JOIN supplier_opc_mappings on supplier_opc_mappings.supplier_id = suppliers.id and supplier_opc_mappings.deleted_at IS NULL and supplier_opc_mappings.vaccount_id = suppliers.vaccount_id"
}

// GetPaymentAccountDetailsJoinStr ...
func GetPaymentAccountDetailsJoinStr() string {
	return "LEFT JOIN payment_account_details on payment_account_details.supplier_id = suppliers.id and payment_account_details.vaccount_id = suppliers.vaccount_id"
}
