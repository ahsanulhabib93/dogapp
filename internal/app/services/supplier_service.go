package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	supplierpb "github.com/voonik/goConnect/api/go/ss2/supplier"
	"github.com/voonik/goFramework/pkg/database"
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
	} else {
		err := database.DBAPM(ctx).Model(&supplier).Updates(models.Supplier{
			Name:                     params.GetName(),
			Email:                    params.GetEmail(),
			SupplierType:             utils.SupplierType(params.GetSupplierType()),
			BusinessName:             params.GetBusinessName(),
			Phone:                    params.GetPhone(),
			AlternatePhone:           params.GetAlternatePhone(),
			ShopImageURL:             params.GetShopImageUrl(),
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

// UpdateStatus of Supplier
func (ss *SupplierService) UpdateStatus(ctx context.Context, params *supplierpb.UpdateStatusParam) (*supplierpb.BasicApiResponse, error) {
	log.Printf("UpdateStatusParams: %+v", params)
	resp := supplierpb.BasicApiResponse{Success: false}

	supplier := models.Supplier{}
	result := database.DBAPM(ctx).Model(&models.Supplier{}).First(&supplier, params.GetId())
	if result.RecordNotFound() {
		resp.Message = "Supplier Not Found"
	} else {
		err := database.DBAPM(ctx).Model(&supplier).Updates(models.Supplier{
			Status: models.SupplierStatus(params.GetStatus()),
		})
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

	if params.GetUploadType() != "SupplierShopImage" {
		resp.Message = "Invalid Upload Type"
	} else {
		object := utils.GetObjectName("shop_images", "", "")
		bucketName := utils.GetBucketName(ctx)
		fileURL, err := utils.GetUploadURL(ctx, bucketName, object)

		log.Printf("GetUploadUrl: %+v", fileURL)
		if err != nil {
			resp.Message = err.Error()
		} else {
			resp = &supplierpb.UrlResponse{
				Url:     fileURL,
				Path:    object,
				Success: true,
				Message: "Fetched upload url successfully",
			}
		}
	}

	log.Printf("GetUploadUrlResponse: %+v", resp)
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
		"suppliers.reason",
		"GROUP_CONCAT( DISTINCT supplier_category_mappings.category_id) as category_ids",
		"GROUP_CONCAT( DISTINCT supplier_opc_mappings.processing_center_id) as opc_ids",
	}

	return strings.Join(s, ",")
}
