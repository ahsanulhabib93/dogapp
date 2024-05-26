package helpers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/shopuptech/go-libs/logger"
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

	allowedServiceTypes := GetServiceTypesForFiltering(ctx, params.GetServiceTypes())
	if len(allowedServiceTypes) == 0 {
		log.Println("User does not have permission to view any service type")
		return query.Where("1=0")
	} else {
		serviceTypes := ParseServiceTypes(ctx, allowedServiceTypes)
		query = query.Where("partner_service_mappings.service_type IN (?)", serviceTypes)
	}

	if len(params.GetTypes()) != 0 {
		query = query.Where("partner_service_mappings.partner_service_level_id IN (?)", params.GetTypes())
	}

	if len(params.GetServiceLevels()) != 0 {
		serviceLevelIds := getServiceLevelForFiltering(ctx, params.GetServiceLevels())
		query = query.Where("partner_service_mappings.partner_service_level_id IN (?)", serviceLevelIds)
	}
	return query
}

func getServiceLevelForFiltering(ctx context.Context, names []string) []uint64 {
	var serviceLevels []uint64
	serviceLevelNameMapping := GetServiceLevelIdMapping(ctx)
	for _, name := range names {
		if val, ok := serviceLevelNameMapping[name]; ok {
			serviceLevels = append(serviceLevels, val)
		}
	}
	return serviceLevels
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
	allowedServiceTypes := GetAllowedServiceTypes(ctx)
	serviceTypes := ParseServiceTypes(ctx, allowedServiceTypes)

	supplierIDs := make([]uint64, len(suppliersData))
	for i, supplierData := range suppliersData {
		supplierIDs[i] = supplierData.ID
	}

	suppliers := []models.Supplier{}
	database.DBAPM(ctx).Model(&models.Supplier{}).
		Preload("PartnerServiceMappings", "partner_service_mappings.service_type IN (?)", serviceTypes).
		Where("id IN (?)", supplierIDs).
		Find(&suppliers)

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
	err := json.Unmarshal(temp, supplierObject)
	if err != nil {
		logger.Log().Errorf("Unmarshal Error: %+v", err)
	}

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

	supplierObject.PartnerServices = GetPartnerServiceMappings(ctx, supplier)
	return supplierObject
}

func GetPartnerServiceMappings(ctx context.Context, supplier models.Supplier) []*supplierPb.PartnerServiceObject {
	partnerServiceData := []*supplierPb.PartnerServiceObject{}
	serviceLevelNameMapping := GetServiceLevelNameMapping(ctx)
	partnerServices := supplier.PartnerServiceMappings // preloaded
	for _, partnerService := range partnerServices {
		partnerServiceData = append(partnerServiceData, &supplierPb.PartnerServiceObject{
			Id:              partnerService.ID,
			Active:          partnerService.Active,
			AgreementUrl:    partnerService.AgreementUrl,
			TradeLicenseUrl: partnerService.TradeLicenseUrl,
			ServiceType:     partnerService.ServiceType.String(),
			ServiceLevel:    serviceLevelNameMapping[partnerService.PartnerServiceLevelID],
		})
	}

	return partnerServiceData
}

func GetAttachments(ctx context.Context, supplierId uint64, partnerServices []*supplierPb.PartnerServiceObject) []*supplierPb.AttachmentObject {
	var attachmentData []*supplierPb.AttachmentObject
	var attachments []models.Attachment
	var partnerServiceIds []uint64

	for _, partnerService := range partnerServices {
		partnerServiceIds = append(partnerServiceIds, partnerService.Id)
	}

	err := database.DBAPM(ctx).Model(&models.Attachment{}).
		Where("(attachable_id = ? AND attachable_type = ?) OR (attachable_id IN (?) AND attachable_type = ?)",
			supplierId, utils.AttachableTypeSupplier, partnerServiceIds, utils.AttachableTypePartnerServiceMapping).
		Find(&attachments).Error
	if err != nil {
		return attachmentData
	}

	for _, attachment := range attachments {
		attachmentData = append(attachmentData, &supplierPb.AttachmentObject{
			Id:              attachment.ID,
			AttachableType:  uint64(attachment.AttachableType),
			AttachableId:    attachment.AttachableID,
			FileType:        attachment.FileType.String(),
			FileUrl:         attachment.FileURL,
			ReferenceNumber: attachment.ReferenceNumber,
		})
	}

	return attachmentData
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

func verify(ctx context.Context, supplier *models.Supplier) error {
	paymentAccountsCount := database.DBAPM(ctx).Model(supplier).Association("PaymentAccountDetails").Count()
	if paymentAccountsCount == 0 {
		return errors.New("At least one payment account details should be present")
	}

	addressesCount := database.DBAPM(ctx).Model(supplier).Association("SupplierAddresses").Count()
	if addressesCount == 0 {
		return errors.New("At least one supplier address should be present")
	}

	if !(supplier.IsOTPVerified() || supplier.IsAnyDocumentPresent(ctx)) {
		return errors.New("At least one primary document or OTP verification needed")
	}

	// TBD: How to handle if multiple service mappings present
	partnerService := models.PartnerServiceMapping{}
	database.DBAPM(ctx).Model(models.PartnerServiceMapping{}).Where("supplier_id = ?", supplier.ID).First(&partnerService)
	serviceLevelNameMapping := GetServiceLevelNameMapping(ctx)
	typeValue := serviceLevelNameMapping[partnerService.PartnerServiceLevelID]

	otpTypeVerificationList := aaaModels.GetAppPreferenceServiceInstance().GetValue(ctx, "enabled_otp_verification", []string{}).([]string)
	if utils.IsInclude(otpTypeVerificationList, typeValue) && !supplier.IsOTPVerified() {
		msg := fmt.Sprint("OTP verification required for supplier type: ", typeValue)
		return errors.New(msg)
	}

	docTypeVerificationList := aaaModels.GetAppPreferenceServiceInstance().GetValue(ctx, "enabled_primary_doc_verification", []string{}).([]string)
	if utils.IsInclude(docTypeVerificationList, typeValue) && !(supplier.IsAnyDocumentPresent(ctx)) {
		msg := fmt.Sprint("At least one primary document required for supplier type: ", typeValue)
		return errors.New(msg)
	}

	return nil
}

func IsValidStatusUpdate(ctx context.Context, supplier models.Supplier, newStatus models.SupplierStatus) (valid bool, message string) {
	if !isValidStatus(newStatus) {
		return false, "Invalid Status"
	} else if !isValidStatusTransition(supplier.Status, newStatus) {
		return false, "Status transition not allowed"
	} else if newStatus == models.SupplierStatusVerified {
		err := verify(ctx, &supplier)
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

func GetServiceTypesForFiltering(ctx context.Context, serviceTypes []string) []string {
	allowedServiceTypes := GetAllowedServiceTypes(ctx)

	// if no service type is passed in filter, then return allowed service types
	if len(serviceTypes) == 0 {
		return allowedServiceTypes
	}

	// use common service types of allowed types and api filter types
	allowedServiceTypes = utils.GetCommonElements(allowedServiceTypes, serviceTypes)
	return allowedServiceTypes
}

func GetAllowedServiceTypes(ctx context.Context) []string {
	var allowedServiceTypes []string
	permissions := utils.GetCurrentUserPermissions(ctx)

	frameworkUser := utils.GetCurrentUserID(ctx) == nil
	globalPermission := utils.IsInclude(permissions, "supplierpanel:allservices:view")
	allServiceAccess := frameworkUser || globalPermission

	for _, serviceType := range utils.PartnerServiceTypeMapping {
		requiedPermission := fmt.Sprintf("supplierpanel:%sservice:view", strings.ToLower(serviceType.String()))
		if allServiceAccess || utils.IsInclude(permissions, requiedPermission) {
			allowedServiceTypes = append(allowedServiceTypes, serviceType.String())
		}
	}

	return allowedServiceTypes
}

func ParseServiceTypes(ctx context.Context, allowedServiceTypes []string) []utils.ServiceType {
	var serviceTypes []utils.ServiceType
	for _, serviceType := range allowedServiceTypes {
		if val, ok := utils.PartnerServiceTypeMapping[serviceType]; ok {
			serviceTypes = append(serviceTypes, val)
		}
	}
	return serviceTypes
}
