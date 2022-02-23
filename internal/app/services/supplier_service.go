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

func (ss *SupplierService) Get(ctx context.Context, params *supplierpb.GetSupplierParam) (*supplierpb.SupplierResponse, error) {
	log.Printf("GetSupplierParam: %+v", params)
	supplier := models.Supplier{}
	result := database.DBAPM(ctx).Model(&models.Supplier{}).Preload("SupplierAddresses").
		Preload("PaymentAccountDetails").First(&supplier, params.GetId())
	if result.RecordNotFound() {
		return &supplierpb.SupplierResponse{Success: false}, nil
	}

	resp := helpers.SupplierDBResponse{}
	database.DBAPM(ctx).Model(&models.Supplier{}).Where("suppliers.id = ?", supplier.ID).
		Joins(" left join supplier_category_mappings on supplier_category_mappings.supplier_id=suppliers.id").
		Joins(" left join supplier_opc_mappings on supplier_opc_mappings.supplier_id=suppliers.id").Group("id").
		Where("supplier_category_mappings.deleted_at IS NULL").Where("supplier_opc_mappings.deleted_at IS NULL").
		Select(ss.getResponseField()).Scan(&resp)

	resp.SupplierAddresses = supplier.SupplierAddresses
	resp.PaymentAccountDetails = supplier.PaymentAccountDetails
	log.Printf("GetSupplierResponse: %+v", resp)
	return &supplierpb.SupplierResponse{
		Success: true,
		Data:    helpers.PrepareSupplierResponse(resp)}, nil
}

// List ...
func (ss *SupplierService) List(ctx context.Context, params *supplierpb.ListParams) (*supplierpb.ListResponse, error) {
	log.Printf("ListSupplierParams: %+v", params)
	suppliers := []helpers.SupplierDBResponse{}
	query := database.DBAPM(ctx).Model(&models.Supplier{})
	query = helpers.PrepareFilter(query, params)
	var total uint64
	query.Count(&total)

	helpers.SetPage(query, params)
	query.Joins(" left join supplier_category_mappings on supplier_category_mappings.supplier_id=suppliers.id").
		Joins(" left join supplier_opc_mappings on supplier_opc_mappings.supplier_id=suppliers.id").Group("id").
		Where("supplier_category_mappings.deleted_at IS NULL").Where("supplier_opc_mappings.deleted_at IS NULL").
		Select(ss.getResponseField()).Scan(&suppliers)

	resp := helpers.PrepareListResponse(suppliers, total)
	log.Printf("ListSupplierResponse: %+v", resp)
	return &resp, nil
}

// ListWithSupplierAddresses ...
func (ss *SupplierService) ListWithSupplierAddresses(ctx context.Context, params *supplierpb.ListParams) (*supplierpb.ListResponse, error) {
	log.Printf("ListwithAddressParams: %+v", params)
	resp := supplierpb.ListResponse{}

	query := database.DBAPM(ctx).Model(&models.Supplier{})
	query = helpers.PrepareFilter(query, params)

	var total uint64
	query.Count(&total)
	helpers.SetPage(query, params)
	suppliersWithAddresses := []models.Supplier{{}}
	query.Joins("join supplier_addresses on supplier_addresses.supplier_id=suppliers.id").
		Select("distinct suppliers.*").Preload("SupplierAddresses").Find(&suppliersWithAddresses)

	temp, _ := json.Marshal(suppliersWithAddresses)
	json.Unmarshal(temp, &resp.Data)
	resp.TotalCount = total
	log.Printf("ListwithAddressResponse: %+v", resp)
	return &resp, nil
}

// Add ...
func (ss *SupplierService) Add(ctx context.Context, params *supplierpb.SupplierParam) (*supplierpb.BasicApiResponse, error) {
	log.Printf("AddSupplierParams: %+v", params)
	resp := supplierpb.BasicApiResponse{Success: false}
	supplier := models.Supplier{
		Name:                     params.GetName(),
		Email:                    params.GetEmail(),
		Status:                   params.GetStatus(),
		SupplierType:             utils.SupplierType(params.GetSupplierType()),
		SupplierCategoryMappings: helpers.PrepareCategoreMapping(params.GetCategoryIds()),
		SupplierOpcMappings:      helpers.PrepareOpcMapping(params.GetOpcIds(), params.GetIsSaUser()),
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

// Edit ...
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
			Status:                   params.GetStatus(),
			SupplierType:             utils.SupplierType(params.GetSupplierType()),
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

func (ss *SupplierService) getResponseField() string {
	s := []string{
		"suppliers.id",
		"suppliers.status",
		"suppliers.supplier_type",
		"suppliers.name",
		"suppliers.email",
		"GROUP_CONCAT( DISTINCT supplier_category_mappings.category_id) as category_ids",
		"GROUP_CONCAT( DISTINCT supplier_opc_mappings.processing_center_id) as opc_ids",
	}

	return strings.Join(s, ",")
}
