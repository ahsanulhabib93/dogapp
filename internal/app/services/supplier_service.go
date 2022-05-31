package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	supplierpb "github.com/voonik/goConnect/api/go/ss2/supplier"
	aaaModels "github.com/voonik/goFramework/pkg/aaa/models"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/goFramework/pkg/file_helper"
	"github.com/voonik/ss2/internal/app/helpers"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/utils"
)

// SupplierService ...
type SupplierService struct{}

// Get Supplier Info
func (ss *SupplierService) Get(ctx context.Context, params *supplierpb.GetSupplierParam) (*supplierpb.SupplierResponse, error) {
	log.Printf("GetSupplierParam: %+v", params)
	supplier := models.Supplier{}
	result := database.DBAPM(ctx).Model(&models.Supplier{}).Preload("SupplierAddresses").First(&supplier, params.GetId())
	if result.RecordNotFound() {
		return &supplierpb.SupplierResponse{Success: false}, nil
	}

	paymentDetails := []*supplierpb.PaymentAccountDetailObject{}
	database.DBAPM(ctx).Model(&models.PaymentAccountDetail{}).Where("supplier_id = ?", params.GetId()).
		Joins(models.GetBankJoinStr()).Select("payment_account_details.*, banks.name bank_name").
		Scan(&paymentDetails)

	resp := helpers.SupplierDBResponse{}
	database.DBAPM(ctx).Model(&models.Supplier{}).
		Joins(models.GetCategoryMappingJoinStr()).Joins(models.GetOpcMappingJoinStr()).
		Where("suppliers.id = ?", supplier.ID).Group("suppliers.id").
		Select(ss.getResponseField()).Scan(&resp)

	resp.SupplierAddresses = supplier.SupplierAddresses
	supplierResp := helpers.PrepareSupplierResponse(resp)
	supplierResp.PaymentAccountDetails = paymentDetails
	log.Printf("GetSupplierResponse: %+v", resp)
	return &supplierpb.SupplierResponse{Success: true, Data: supplierResp}, nil
}

// List Suppliers
func (ss *SupplierService) List(ctx context.Context, params *supplierpb.ListParams) (*supplierpb.ListResponse, error) {
	log.Printf("ListSupplierParams: %+v", params)
	suppliers := []helpers.SupplierDBResponse{}
	query := database.DBAPM(ctx).Model(&models.Supplier{})
	query = helpers.PrepareFilter(ctx, query, params).
		Joins(models.GetCategoryMappingJoinStr()).Joins(models.GetOpcMappingJoinStr()).
		Group("suppliers.id")

	var total uint64
	query.Count(&total)
	helpers.SetPage(query, params)
	query.Select(ss.getResponseField()).Scan(&suppliers)
	resp := helpers.PrepareListResponse(suppliers, total)
	log.Printf("ListSupplierResponse: %+v", resp)
	return &resp, nil
}

// ListWithSupplierAddresses ...
func (ss *SupplierService) ListWithSupplierAddresses(ctx context.Context, params *supplierpb.ListParams) (*supplierpb.ListResponse, error) {
	log.Printf("ListwithAddressParams: %+v", params)
	resp := supplierpb.ListResponse{}

	query := database.DBAPM(ctx).Model(&models.Supplier{})
	query = helpers.PrepareFilter(ctx, query, params).
		Preload("SupplierAddresses").Joins("join supplier_addresses on supplier_addresses.supplier_id=suppliers.id").
		Group("suppliers.id")

	var total uint64
	query.Count(&total)
	helpers.SetPage(query, params)
	suppliersWithAddresses := []models.Supplier{{}}
	query.Select("suppliers.*").Find(&suppliersWithAddresses)

	temp, _ := json.Marshal(suppliersWithAddresses)
	json.Unmarshal(temp, &resp.Data)
	resp.TotalCount = total
	log.Printf("ListwithAddressResponse: %+v", resp)
	return &resp, nil
}

// Add Supplier
func (ss *SupplierService) Add(ctx context.Context, params *supplierpb.SupplierParam) (*supplierpb.BasicApiResponse, error) {
	log.Printf("AddSupplierParams: %+v", params)
	resp := supplierpb.BasicApiResponse{Success: false}
	if err := helpers.IsOpcListValid(ctx, params.GetOpcIds()); err != nil {
		resp.Message = err.Error()
		return &resp, nil
	}

	supplier := models.Supplier{
		Name:                     params.GetName(),
		Email:                    params.GetEmail(),
		UserID:                   utils.GetCurrentUserID(ctx),
		SupplierType:             utils.SupplierType(params.GetSupplierType()),
		BusinessName:             params.GetBusinessName(),
		Phone:                    params.GetPhone(),
		AlternatePhone:           params.GetAlternatePhone(),
		ShopImageURL:             params.GetShopImageUrl(),
		NidNumber:                params.GetNidNumber(),
		NidFrontImageUrl:         params.GetNidFrontImageUrl(),
		NidBackImageUrl:          params.GetNidBackImageUrl(),
		TradeLicenseUrl:          params.GetTradeLicenseUrl(),
		AgreementUrl:             params.GetAgreementUrl(),
		ShopOwnerImageUrl:        params.GetShopOwnerImageUrl(),
		GuarantorImageUrl:        params.GetGuarantorImageUrl(),
		ChequeImageUrl:           params.GetChequeImageUrl(),
		SupplierCategoryMappings: helpers.PrepareCategoreMapping(params.GetCategoryIds()),
		SupplierOpcMappings:      helpers.PrepareOpcMapping(ctx, params.GetOpcIds(), params.GetCreateWithOpcMapping()),
		SupplierAddresses:        helpers.PrepareSupplierAddress(params),
	}

	err := database.DBAPM(ctx).Save(&supplier)
	if err != nil && err.Error != nil {
		resp.Message = fmt.Sprintf("Error while creating Supplier: %s", err.Error)
	} else {
		resp.Message = "Supplier Added Successfully"
		resp.Success = true
		resp.Id = supplier.ID
	}
	log.Printf("AddSupplierResponse: %+v", resp)
	return &resp, nil
}

// Edit Supplier Details
func (ss *SupplierService) Edit(ctx context.Context, params *supplierpb.SupplierObject) (*supplierpb.BasicApiResponse, error) {
	log.Printf("EditSupplierParams: %+v", params)
	resp := supplierpb.BasicApiResponse{Success: false}

	supplier := models.Supplier{}
	query := database.DBAPM(ctx).Model(&models.Supplier{})
	if params.GetCategoryIds() != nil {
		query = query.Preload("SupplierCategoryMappings")
	}
	result := query.First(&supplier, params.GetId())
	if result.RecordNotFound() {
		resp.Message = "Supplier Not Found"
	} else if !supplier.IsChangeAllowed(ctx) {
		resp.Message = "Change Not Allowed"
	} else {
		var isPhoneVerified bool
		if params.GetPhone() != "" && params.GetPhone() != supplier.Phone {
			isPhoneVerified = false
		} else {
			isPhoneVerified = *supplier.IsPhoneVerified
		}

		var status models.SupplierStatus
		if supplier.Status == models.SupplierStatusVerified || supplier.Status == models.SupplierStatusFailed {
			status = models.SupplierStatusPending // Moving to Pending if any data is updated
		}

		err := database.DBAPM(ctx).Model(&supplier).Updates(models.Supplier{
			Status:                   status,
			Name:                     params.GetName(),
			Email:                    params.GetEmail(),
			SupplierType:             utils.SupplierType(params.GetSupplierType()),
			BusinessName:             params.GetBusinessName(),
			Phone:                    params.GetPhone(),
			AlternatePhone:           params.GetAlternatePhone(),
			IsPhoneVerified:          &isPhoneVerified,
			ShopImageURL:             params.GetShopImageUrl(),
			NidNumber:                params.GetNidNumber(),
			NidFrontImageUrl:         params.GetNidFrontImageUrl(),
			NidBackImageUrl:          params.GetNidBackImageUrl(),
			TradeLicenseUrl:          params.GetTradeLicenseUrl(),
			AgreementUrl:             params.GetAgreementUrl(),
			ShopOwnerImageUrl:        params.GetShopOwnerImageUrl(),
			GuarantorImageUrl:        params.GetGuarantorImageUrl(),
			ChequeImageUrl:           params.GetChequeImageUrl(),
			SupplierCategoryMappings: helpers.UpdateSupplierCategoryMapping(ctx, supplier.ID, params.GetCategoryIds()),
		})
		if err != nil && err.Error != nil {
			resp.Message = fmt.Sprintf("Error while updating Supplier: %s", err.Error)
		} else {
			resp.Message = "Supplier Edited Successfully"
			resp.Success = true
		}
	}
	log.Printf("EditSupplierResponse: %+v", resp)
	return &resp, nil
}

// RemoveDocument...
func (ss *SupplierService) RemoveDocument(ctx context.Context, params *supplierpb.RemoveDocumentParam) (*supplierpb.BasicApiResponse, error) {
	log.Printf("RemoveDocumentParam: %+v", params)
	resp := supplierpb.BasicApiResponse{Success: false}

	supplier := models.Supplier{}
	query := database.DBAPM(ctx).Model(&models.Supplier{})
	result := query.First(&supplier, params.GetId())
	if result.RecordNotFound() {
		resp.Message = "Supplier Not Found"
	} else if !supplier.IsChangeAllowed(ctx) {
		resp.Message = "Change Not Allowed"
	} else {
		isPrimaryDoc := utils.IsInclude(utils.SupplierPrimaryDocumentType, params.GetDocumentType())
		isSecondaryDoc := utils.IsInclude(utils.SupplierSecondaryDocumentType, params.GetDocumentType())
		if !(isPrimaryDoc || isSecondaryDoc) {
			resp.Message = "Invalid Document Type"
		} else {
			query := database.DBAPM(ctx).Model(&supplier).Where("suppliers.id = ?", params.GetId())
			query = query.Update(params.GetDocumentType(), "")
			if isPrimaryDoc &&
				(supplier.Status == models.SupplierStatusVerified || supplier.Status == models.SupplierStatusFailed) {
				query = query.Update("status", models.SupplierStatusPending) // Moving to Pending if any data is updated
			}

			if err := query.Error; err != nil {
				resp.Message = fmt.Sprintf("Error While Removing Supplier Document: %s", err.Error())
			} else {
				resp.Message = fmt.Sprintf("Supplier %s Removed Successfully", params.GetDocumentType())
				resp.Success = true
			}
		}
	}
	log.Printf("RemoveDocumentResponse: %+v", resp)
	return &resp, nil
}

// UpdateStatus of Supplier
func (ss *SupplierService) UpdateStatus(ctx context.Context, params *supplierpb.UpdateStatusParam) (*supplierpb.BasicApiResponse, error) {
	log.Printf("UpdateStatusParams: %+v", params)
	resp := supplierpb.BasicApiResponse{Success: false}

	supplier := models.Supplier{}
	newSupplierStatus := models.SupplierStatus(params.GetStatus())
	result := database.DBAPM(ctx).Model(&models.Supplier{}).First(&supplier, params.GetId())
	fmt.Printf("newSupplierStatus %+v", newSupplierStatus)
	if result.RecordNotFound() {
		resp.Message = "Supplier Not Found"
	} else if valid, message := helpers.IsValidStatusUpdate(ctx, supplier, newSupplierStatus); !valid {
		resp.Message = message
	} else {
		updateDetails := map[string]interface{}{"status": newSupplierStatus, "reason": params.GetReason()} //to allow empty string update for reason
		err := database.DBAPM(ctx).Model(&supplier).Updates(updateDetails)
		if err != nil && err.Error != nil {
			resp.Message = fmt.Sprintf("Error while updating Supplier: %s", err.Error)
		} else {
			resp.Message = "Supplier status updated successfully"
			resp.Success = true
		}
	}
	log.Printf("UpdateStatusResponse: %+v", resp)
	return &resp, nil
}

// SupplierMap - Supplier OPC mapping
func (ss *SupplierService) SupplierMap(ctx context.Context, params *supplierpb.SupplierMappingParams) (*supplierpb.BasicApiResponse, error) {
	supplier := &models.Supplier{}
	result := database.DBAPM(ctx).Model(&models.Supplier{}).First(supplier, params.GetSupplierId())
	if result.RecordNotFound() {
		return &supplierpb.BasicApiResponse{Message: "Supplier Not Found"}, nil
	}

	isDeleting := false
	if strings.ToLower(params.OperationType) == "delete" {
		isDeleting = true
	}

	switch strings.ToLower(params.MapWith) {
	case "opc":
		return helpers.UpdateSupplierOpcMapping(ctx, params.SupplierId, params.Id, isDeleting), nil
	}

	return &supplierpb.BasicApiResponse{Message: "Invalid mapping option"}, nil
}

// GetUploadURL for uploading supplier shop image
func (ss *SupplierService) GetUploadURL(ctx context.Context, params *supplierpb.GetUploadUrlParam) (*supplierpb.UrlResponse, error) {
	log.Printf("GetUploadUrlParams: %+v", params)
	resp := &supplierpb.UrlResponse{Success: false}

	if details, allowed := utils.AllowedUploadType[params.GetUploadType()]; !allowed {
		resp.Message = "Invalid File Type"
	} else {
		filePath := file_helper.GenerateFilePath(utils.BucketFolder, details[0], details[1], details[2])
		bucketName := utils.GetBucketName(ctx)
		fileURL, err := file_helper.GetUploadURL(ctx, bucketName, filePath)

		log.Printf("GetUploadUrl: %+v", fileURL)
		if err != nil {
			resp.Message = err.Error()
		} else {
			resp = &supplierpb.UrlResponse{
				Url:     fileURL,
				Path:    filePath,
				Success: true,
				Message: "Fetched upload url successfully",
			}
		}
	}

	log.Printf("GetUploadUrlResponse: %+v", resp)
	return resp, nil
}

// GetDownloadURL to get supplier shop image
func (ss *SupplierService) GetDownloadURL(ctx context.Context, params *supplierpb.GetDownloadUrlParam) (*supplierpb.UrlResponse, error) {
	log.Printf("GetDownloadUrlParams: %+v", params)
	resp := &supplierpb.UrlResponse{Success: false}

	filePath := params.GetPath()
	bucketName := utils.GetBucketName(ctx)
	fileURL, err := file_helper.GetDownloadURL(ctx, bucketName, filePath)

	log.Printf("GetDownloadUrl: %+v", fileURL)
	if err != nil {
		resp.Message = err.Error()
	} else {
		resp = &supplierpb.UrlResponse{
			Url:     fileURL,
			Path:    filePath,
			Success: true,
			Message: "Fetched url successfully",
		}
	}

	log.Printf("GetDownloadUrlResponse: %+v", resp)
	return resp, nil
}

// SendVerificationOtp ...
func (ss *SupplierService) SendVerificationOtp(ctx context.Context, params *supplierpb.SendOtpParam) (*supplierpb.BasicApiResponse, error) {
	log.Printf("SendVerificationOtpParams: %+v", params)
	resp := &supplierpb.BasicApiResponse{Success: false}

	supplierID := params.GetSupplierId()
	supplier := models.Supplier{}
	result := database.DBAPM(ctx).Model(&models.Supplier{}).First(&supplier, supplierID)
	if result.RecordNotFound() {
		resp.Message = "Supplier Not Found"
	} else if *supplier.IsPhoneVerified {
		resp.Message = "Phone number is already verified"
	} else {
		content := aaaModels.GetAppPreferenceServiceInstance().GetValue(ctx, "supplier_phone_verification_otp_content", "OTP for supplier verification: $otp").(string)
		otpResponse := helpers.SendOtpAPI(ctx, supplierID, supplier.Phone, content, params.Resend)

		if !otpResponse.Success {
			resp.Message = otpResponse.Message
		} else {
			resp.Message = otpResponse.Message
			resp.Success = true
		}
	}

	log.Printf("SendVerificationOtpResponse: %+v", resp)
	return resp, nil
}

// VerifyOtp ...
func (ss *SupplierService) VerifyOtp(ctx context.Context, params *supplierpb.VerifyOtpParam) (*supplierpb.BasicApiResponse, error) {
	log.Printf("VerifyOtpParams: %+v", params)
	resp := &supplierpb.BasicApiResponse{Success: false}

	supplierID := params.GetSupplierId()
	supplier := models.Supplier{}
	result := database.DBAPM(ctx).Model(&models.Supplier{}).First(&supplier, supplierID)
	if result.RecordNotFound() {
		resp.Message = "Supplier Not Found"
	} else if *supplier.IsPhoneVerified {
		resp.Message = "Phone number is already verified"
	} else {
		otpResponse := helpers.VerifyOtpAPI(ctx, supplierID, params.OtpCode)

		if !otpResponse.Success {
			resp.Message = otpResponse.Message
		} else {
			database.DBAPM(ctx).Model(&supplier).Update("IsPhoneVerified", true)
			resp.Message = otpResponse.Message
			resp.Success = true
		}
	}

	log.Printf("VerifyOtpResponse: %+v", resp)
	return resp, nil
}

func (ss *SupplierService) getResponseField() string {
	s := []string{
		"suppliers.id",
		"suppliers.status",
		"suppliers.supplier_type",
		"suppliers.name",
		"suppliers.email",
		"suppliers.phone",
		"suppliers.alternate_phone",
		"suppliers.is_phone_verified",
		"suppliers.business_name",
		"suppliers.shop_image_url",
		"suppliers.nid_number",
		"suppliers.nid_front_image_url",
		"suppliers.nid_back_image_url",
		"suppliers.trade_license_url",
		"suppliers.agreement_url",
		"suppliers.reason",
		"suppliers.shop_owner_image_url",
		"suppliers.guarantor_image_url",
		"suppliers.cheque_image_url",
		"GROUP_CONCAT( DISTINCT supplier_category_mappings.category_id) as category_ids",
		"GROUP_CONCAT( DISTINCT supplier_opc_mappings.processing_center_id) as opc_ids",
	}

	return strings.Join(s, ",")
}
