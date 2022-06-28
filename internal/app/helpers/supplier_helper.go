package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"

	supplierpb "github.com/voonik/goConnect/api/go/ss2/supplier"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/utils"
)

type SupplierDBResponse struct {
	models.Supplier
	CategoryIds string `json:"category_ids,omitempty"`
	OpcIds      string `json:"opc_ids,omitempty"`
}

func PrepareFilter(ctx context.Context, query *gorm.DB, params *supplierpb.ListParams) *gorm.DB {
	if params.GetId() != 0 {
		query = query.Where("suppliers.id = ?", params.GetId())
	}
	if len(params.GetSupplierIds()) != 0 {
		query = query.Where("suppliers.id IN (?)", params.GetSupplierIds())
	}
	if params.GetName() != "" {
		searchValue := fmt.Sprintf("%%%s%%", params.GetName())
		query = query.Joins("left join payment_account_details on suppliers.id = payment_account_details.supplier_id").Where("suppliers.name LIKE ? or suppliers.business_name LIKE ? or suppliers.phone LIKE ? or suppliers.id LIKE ? or payment_account_details.account_number LIKE ?", searchValue, searchValue, searchValue, searchValue, searchValue)

	}
	if params.GetEmail() != "" {
		query = query.Where("suppliers.email = ?", params.GetEmail())
	}
	if params.GetPhone() != "" {
		query = query.Where("suppliers.phone = ?", params.GetPhone())
	}
	if params.GetCreatedAtGte() != "" {
		query = query.Where("date(suppliers.created_at) >= ?", params.GetCreatedAtGte())
	}
	if params.GetCreatedAtLte() != "" {
		query = query.Where("date(suppliers.created_at) <= ?", params.GetCreatedAtLte())
	}
	if params.GetStatus() != "" {
		status := strings.Split(params.GetStatus(), ",")
		query = query.Where("suppliers.status IN (?)", status)
	}

	if params.GetCity() != "" {
		query = query.Where("supplier_addresses.city = ?", params.GetCity())
	}

	OpcIds := params.GetOpcIds()
	if params.AssociatedWithCurrentUser {
		OpcIds = append(OpcIds, GetOPCListForCurrentUser(ctx)...)
	}
	if len(OpcIds) > 0 || params.AssociatedWithCurrentUser {
		query = query.Where("supplier_opc_mappings.processing_center_id IN (?)", OpcIds)
	}
	if params.GetOpcId() != 0 {
		query = query.Where("supplier_opc_mappings.processing_center_id = ?", params.GetOpcId())
	}

	return query
}

func SetPage(query *gorm.DB, params *supplierpb.ListParams) {
	if params.GetPerPage() <= 0 || params.GetPerPage() > utils.DEFAULT_PER_PAGE {
		params.PerPage = utils.DEFAULT_PER_PAGE
	}

	params.Page = utils.Int64Max(utils.DEFAULT_PAGE, params.GetPage())
	offset := (params.GetPage() - 1) * params.GetPerPage()
	*query = *query.Offset(offset).Limit(params.GetPerPage())
}

func PrepareCategoreMapping(ids []uint64) []models.SupplierCategoryMapping {
	categories := []models.SupplierCategoryMapping{}
	for _, id := range ids {
		categories = append(categories, models.SupplierCategoryMapping{
			CategoryID: id,
		})
	}

	return categories
}

func PrepareOpcMapping(ctx context.Context, ids []uint64, fetchOpc bool) []models.SupplierOpcMapping {
	if fetchOpc {
		ids = append(ids, GetOPCListForCurrentUser(ctx)...)
	}

	processCenters := []models.SupplierOpcMapping{}
	for _, id := range ids {
		processCenters = append(processCenters, models.SupplierOpcMapping{
			ProcessingCenterID: id,
		})
	}
	return processCenters
}

func PrepareListResponse(suppliers []SupplierDBResponse, total uint64) supplierpb.ListResponse {
	data := []*supplierpb.SupplierObject{}
	for _, supplier := range suppliers {
		data = append(data, PrepareSupplierResponse(supplier))
	}

	return supplierpb.ListResponse{Data: data, TotalCount: total}
}

func PrepareSupplierResponse(supplier SupplierDBResponse) *supplierpb.SupplierObject {
	temp, _ := json.Marshal(supplier)
	so := &supplierpb.SupplierObject{}
	json.Unmarshal(temp, so)

	so.CategoryIds = []uint64{}
	for _, cId := range strings.Split(supplier.CategoryIds, ",") {
		if cId := strings.TrimSpace(cId); cId != "" {
			v, _ := strconv.Atoi(cId)
			so.CategoryIds = append(so.CategoryIds, uint64(v))
		}
	}

	so.OpcIds = []uint64{}
	for _, saId := range strings.Split(supplier.OpcIds, ",") {
		if opcId := strings.TrimSpace(saId); opcId != "" {
			v, _ := strconv.Atoi(saId)
			so.OpcIds = append(so.OpcIds, uint64(v))
		}
	}

	return so
}

func PrepareSupplierAddress(params *supplierpb.SupplierParam) []models.SupplierAddress {
	if params.GetFirstname() == "" && params.GetLastname() == "" && params.GetAddress1() == "" && params.GetAddress2() == "" &&
		params.GetLandmark() == "" && params.GetCity() == "" && params.GetState() == "" && params.GetCountry() == "" &&
		params.GetZipcode() == "" && params.GetGstNumber() == "" {
		return []models.SupplierAddress{}
	}

	return []models.SupplierAddress{{
		Firstname: params.GetFirstname(),
		Lastname:  params.GetLastname(),
		Address1:  params.GetAddress1(),
		Address2:  params.GetAddress2(),
		Landmark:  params.GetLandmark(),
		City:      params.GetCity(),
		State:     params.GetState(),
		Country:   params.GetCountry(),
		Zipcode:   params.GetZipcode(),
		Phone:     params.GetPhone(),
		GstNumber: params.GetGstNumber(),
		IsDefault: true,
	}}
}

//IsValidStatusUpdate ...
func IsValidStatusUpdate(ctx context.Context, supplier models.Supplier, newStatus models.SupplierStatus) (valid bool, message string) {
	if !isValidStatus(newStatus) {
		return false, "Invalid Status"
	} else if !isValidStatusTransition(supplier.Status, newStatus) {
		return false, "Status transition not allowed"
	} else if newStatus == models.SupplierStatusVerified {
		err := supplier.Verify(ctx)
		if err != nil {
			return false, err.Error()
		}
	}
	return true, ""
}

func isValidStatus(newStatus models.SupplierStatus) (valid bool) {
	validStates := []models.SupplierStatus{models.SupplierStatusPending, models.SupplierStatusVerified, models.SupplierStatusFailed, models.SupplierStatusBlocked}
	for _, status := range validStates {
		if status == newStatus {
			valid = true
			break
		}
	}
	return
}

func isValidStatusTransition(oldStatus, newStatus models.SupplierStatus) (valid bool) {
	validStateTransitions := map[models.SupplierStatus][]models.SupplierStatus{
		models.SupplierStatusPending:  {models.SupplierStatusVerified, models.SupplierStatusFailed, models.SupplierStatusBlocked},
		models.SupplierStatusVerified: {models.SupplierStatusBlocked},
		models.SupplierStatusBlocked:  {models.SupplierStatusVerified},
	}
	for fromStatus, toStates := range validStateTransitions {
		if oldStatus == fromStatus {
			for _, toStatus := range toStates {
				if newStatus == toStatus {
					valid = true
				}
			}
		}
	}
	return
}
