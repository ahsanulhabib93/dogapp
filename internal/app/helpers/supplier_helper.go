package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	supplierPb "github.com/voonik/goConnect/api/go/ss2/supplier"
	aaaModels "github.com/voonik/goFramework/pkg/aaa/models"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/utils"
)

type SupplierDBResponse struct {
	models.Supplier
	CategoryIds string `json:"category_ids,omitempty"`
	OpcIds      string `json:"opc_ids,omitempty"`
}

func PrepareFilter(ctx context.Context, query *gorm.DB, params *supplierPb.ListParams) *gorm.DB {
	if params.GetId() != 0 {
		query = query.Where("suppliers.id = ?", params.GetId())
	}
	if len(params.GetSupplierIds()) != 0 {
		query = query.Where("suppliers.id IN (?)", params.GetSupplierIds())
	}
	if params.GetName() != "" {
		searchValue := fmt.Sprintf("%%%s%%", params.GetName())
		query = query.Joins(models.GetPaymentAccountDetailsJoinStr()).Where("suppliers.name LIKE ? or suppliers.business_name LIKE ? or suppliers.phone LIKE ? or suppliers.id = ? or payment_account_details.account_number LIKE ?", searchValue, searchValue, searchValue, params.GetName(), searchValue)

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
	if len(params.GetTypes()) != 0 {
		// partner_service_mappings is already joined in places where PrepareFilter is called
		query = query.Where("partner_service_mappings.service_level IN (?)", params.GetTypes())
	}

	return query
}

func SetPage(ctx context.Context, query *gorm.DB, params *supplierPb.ListParams) {
	if params.GetPerPage() <= 0 || params.GetPerPage() > utils.DEFAULT_PER_PAGE {
		params.PerPage = utils.DEFAULT_PER_PAGE
	}

	params.Page = utils.Int64Max(utils.DEFAULT_PAGE, params.GetPage())
	offset := (params.GetPage() - 1) * params.GetPerPage()
	searchLimit := aaaModels.GetAppPreferenceServiceInstance().GetValue(ctx, "query_search_limit", int64(5)).(int64)
	if params.GetName() != "" {
		*query = *query.Offset(offset).Limit(searchLimit)
	} else {
		*query = *query.Offset(offset).Limit(params.GetPerPage())
	}

}

func PrepareCategoryMapping(ids []uint64) []models.SupplierCategoryMapping {
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

func PrepareListResponse(ctx context.Context, suppliersData []SupplierDBResponse) (data []*supplierPb.SupplierObject) {
	supplierIDs := make([]uint64, len(suppliersData))
	for i, supplierData := range suppliersData {
		supplierIDs[i] = supplierData.ID
	}

	suppliers := []models.Supplier{}
	database.DBAPM(ctx).Model(&models.Supplier{}).Preload("PartnerServiceMappings").
		Where("id IN (?)", supplierIDs).Find(&suppliers)

	supplierMap := make(map[uint64]models.Supplier)
	for _, supplier := range suppliers {
		supplierMap[supplier.ID] = supplier
	}

	for _, supplierData := range suppliersData {
		supplier := supplierMap[supplierData.ID]
		data = append(data, PrepareSupplierResponse(ctx, supplier, supplierData))
	}
	return data
}

func PrepareSupplierResponse(ctx context.Context, supplier models.Supplier, supplierData SupplierDBResponse) *supplierPb.SupplierObject {
	temp, _ := json.Marshal(supplierData)
	supplierObject := &supplierPb.SupplierObject{}
	json.Unmarshal(temp, supplierObject)

	supplierObject.CategoryIds = []uint64{}
	for _, cId := range strings.Split(supplierData.CategoryIds, ",") {
		if cId := strings.TrimSpace(cId); cId != "" {
			v, _ := strconv.Atoi(cId)
			supplierObject.CategoryIds = append(supplierObject.CategoryIds, uint64(v))
		}
	}

	supplierObject.OpcIds = []uint64{}
	for _, saId := range strings.Split(supplierData.OpcIds, ",") {
		if opcId := strings.TrimSpace(saId); opcId != "" {
			v, _ := strconv.Atoi(saId)
			supplierObject.OpcIds = append(supplierObject.OpcIds, uint64(v))
		}
	}

	supplierObject.PartnerServiceMappings = GetPartnerServiceMappings(ctx, supplier)
	return supplierObject
}

func GetPartnerServiceMappings(ctx context.Context, supplier models.Supplier) []*supplierPb.PartnerServiceObject {
	partnerServiceData := []*supplierPb.PartnerServiceObject{}

	partnerServices := supplier.PartnerServiceMappings // preloaded
	for _, partnerService := range partnerServices {
		partnerServiceData = append(partnerServiceData, &supplierPb.PartnerServiceObject{
			Id:              partnerService.ID,
			Active:          partnerService.Active,
			AgreementUrl:    partnerService.AgreementUrl,
			TradeLicenseUrl: partnerService.TradeLicenseUrl,
			ServiceType:     partnerService.ServiceType.String(),
			ServiceLevel:    partnerService.ServiceLevel.String(),
		})
	}

	return partnerServiceData
}

func PrepareSupplierAddress(params *supplierPb.SupplierParam) []models.SupplierAddress {
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

// IsValidStatusUpdate ...
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

func CheckSupplierExistWithDifferentRole(ctx context.Context, supplier models.Supplier) error {
	if user := FindCreUserByPhone(ctx, supplier.Phone); user != nil {
		log.Printf("getCreUserWithPhone: phone = %s response = %v\n", supplier.Phone, user)
		return fmt.Errorf("user(#%s) already exist as Retails/SalesRep", supplier.Phone)
	} else if user = FindCreUserByPhone(ctx, supplier.AlternatePhone); user != nil {
		log.Printf("getCreUserWithPhone: alternate_phone = %s response = %v\n", supplier.AlternatePhone, user)
		return fmt.Errorf("user(#%s) already exist as Retails/SalesRep", supplier.AlternatePhone)
	}

	if user := GetIdentityUser(ctx, supplier.Phone); user != nil {
		log.Printf("getCreUserWithPhone: phone = %s response = %v\n", supplier.Phone, user)
		return fmt.Errorf("user(#%s) already exist", supplier.Phone)
	}

	if users := GetTalentXUser(ctx, supplier.Phone); len(users) != utils.Zero {
		log.Printf("GetTalentXUser: phone = %s response = %v\n", supplier.Phone, users)
		return fmt.Errorf("user(#%s) already exist as shopup employee", supplier.Phone)
	}

	return nil
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

func GetDefaultServiceType(ctx context.Context) utils.ServiceType {
	return utils.ServiceType(aaaModels.GetAppPreferenceServiceInstance().GetValue(ctx, "default_service_type", int64(utils.Supplier)).(int64))
}
