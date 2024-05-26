package services

import (
	"context"
	"fmt"
	"log"

	psmpb "github.com/voonik/goConnect/api/go/ss2/partner_service_mapping"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/helpers"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/utils"
)

type PartnerServiceMappingService struct{}

func (psm *PartnerServiceMappingService) Add(ctx context.Context, params *psmpb.PartnerServiceObject) (*psmpb.BasicApiResponse, error) {
	log.Printf("PartnerServiceMappingService Add params: %+v", params)
	response := psmpb.BasicApiResponse{Success: false, Message: "Failed to add partner service"}

	serviceType := utils.PartnerServiceTypeMapping[params.GetServiceType()]
	serviceLevel := helpers.GetServiceLevelByTypeAndName(ctx, serviceType, params.GetServiceLevel())

	if (serviceType == 0) || (serviceLevel.ID == 0) {
		response.Message = "Invalid Service Type and/or Service Level"
		return &response, nil
	}

	supplier := models.Supplier{}
	query := database.DBAPM(ctx).Model(&models.Supplier{}).First(&supplier, params.GetSupplierId())

	if query.RecordNotFound() {
		response.Message = "Partner Not Found"
	} else {
		partnerService := models.PartnerServiceMapping{
			SupplierId:            params.GetSupplierId(),
			ServiceType:           serviceType,
			PartnerServiceLevelID: serviceLevel.ID,
			Active:                true,
			TradeLicenseUrl:       params.GetTradeLicenseUrl(),
			AgreementUrl:          params.GetAgreementUrl(),
		}

		err := database.DBAPM(ctx).Save(&partnerService)
		if err != nil && err.Error != nil {
			response.Message = fmt.Sprintf("Error while creating Partner Service: %s", err.Error)
		} else {
			response.Message = "Partner Service Added Successfully"
			response.Success = true
			response.Id = partnerService.ID
		}
	}

	return &response, nil
}

func (psm *PartnerServiceMappingService) Edit(ctx context.Context, params *psmpb.PartnerServiceObject) (*psmpb.BasicApiResponse, error) {
	log.Printf("PartnerServiceMappingService Edit params: %+v", params)
	response := psmpb.BasicApiResponse{Message: "Partner Service Added Successfully"}

	var serviceType utils.ServiceType
	var serviceLevel models.PartnerServiceLevel

	if params.GetPartnerServiceId() == 0 || params.GetSupplierId() == 0 {
		response.Message = "Invalid Partner/Partner Service ID"
		return &response, nil
	}

	if params.GetServiceType() != "" && params.GetServiceLevel() != "" {
		serviceType = utils.PartnerServiceTypeMapping[params.GetServiceType()]
		serviceLevel = helpers.GetServiceLevelByTypeAndName(ctx, serviceType, params.GetServiceLevel())

		if serviceLevel.ServiceType != serviceType {
			response.Message = "Incompatible Service Type and Service Level"
			return &response, nil
		}
	}

	supplier := models.Supplier{}
	database.DBAPM(ctx).Model(&models.Supplier{}).First(&supplier, params.GetSupplierId())

	partnerService := models.PartnerServiceMapping{}
	partnerServiceQuery := database.DBAPM(ctx).Model(&models.PartnerServiceMapping{}).Where("id = ? and supplier_id = ?", params.GetPartnerServiceId(), params.GetSupplierId()).First(&partnerService)

	if partnerServiceQuery.RecordNotFound() {
		response.Message = "Partner/Partner Service Not Found"
	} else if !helpers.ValidatePartnerSericeEdit(ctx, helpers.PartnerServiceEditEntity{
		ServiceType:    serviceType,
		ServiceLevelId: serviceLevel.ID,
	}, helpers.PartnerServiceEditEntity{
		ServiceType:    partnerService.ServiceType,
		ServiceLevelId: partnerService.PartnerServiceLevelID,
	}) {
		response.Message = "Not allowed to edit Partner Service Info"
	} else {
		err := database.DBAPM(ctx).Model(&partnerService).Updates(models.PartnerServiceMapping{
			PartnerServiceLevelID: serviceLevel.ID,
			TradeLicenseUrl:       params.GetTradeLicenseUrl(),
			AgreementUrl:          params.GetAgreementUrl(),
		})

		errorMsg, _ := psm.updateUserStatus(ctx, supplier, params.GetSupplierId())

		if err != nil && err.Error != nil {
			response.Message = fmt.Sprintf("Error while updating Partner Service: %s", err.Error)
		} else if errorMsg != "" {
			response.Message = errorMsg
		} else {
			response.Message = "Partner Service Edited Successfully"
			response.Success = true
			response.Id = partnerService.ID
		}
	}

	log.Printf("PartnerServiceMappingService Edit response: Success: %+v, Message: %+v, Id: %+v", response.Success, response.Message, response.Id)
	return &response, nil
}

func (psm *PartnerServiceMappingService) UpdateStatus(ctx context.Context, params *psmpb.PartnerServiceObject) (*psmpb.BasicApiResponse, error) {
	log.Printf("PartnerServiceMappingService UpdateStatus params: %+v", params)
	response := psmpb.BasicApiResponse{Message: "Partner Service updated Successfully"}

	supplier := models.Supplier{}
	supplierQuery := database.DBAPM(ctx).Model(&models.Supplier{}).First(&supplier, params.GetSupplierId())

	partnerService := models.PartnerServiceMapping{}
	partnerServiceQuery := database.DBAPM(ctx).Model(&models.PartnerServiceMapping{}).First(&partnerService, params.GetPartnerServiceId())

	if supplierQuery.RecordNotFound() {
		response.Message = "Partner Not Found"
	} else if partnerServiceQuery.RecordNotFound() {
		response.Message = "Partner Service Not Found"
	} else {
		err := database.DBAPM(ctx).Model(&partnerService).Updates(map[string]interface{}{"active": params.GetActive()})

		if err != nil && err.Error != nil {
			response.Message = fmt.Sprintf("Error while updating Partner Service: %s", err.Error)
		} else {
			response.Message = "Partner Service Updated Successfully"
			response.Success = true
			response.Id = partnerService.ID
		}
	}

	log.Printf("PartnerServiceMappingService UpdateStatus response: Success: %+v, Message: %+v, Id: %+v", response.Success, response.Message, response.Id)
	return &response, nil
}

func (psm *PartnerServiceMappingService) PartnerTypesList(ctx context.Context, params *psmpb.PartnerServiceObject) (*psmpb.PartnerTypeListResponse, error) {
	responseMappings := []*psmpb.PartnerServiceTypeMapping{}

	allowedServiceTypes := helpers.GetAllowedServiceTypes(ctx)
	serviceTypeLevelMappings := helpers.GetServiceTypeLevelMappings(ctx)
	for _, serviceType := range utils.PartnerServiceTypeMapping {
		if !utils.Includes(allowedServiceTypes, serviceType.String()) {
			continue // Skipping if not allowed
		}

		object := psmpb.PartnerServiceTypeMapping{}
		object.PartnerType = serviceType.String()
		for _, serviceLevel := range serviceTypeLevelMappings[serviceType] {
			object.ServiceTypes = append(object.ServiceTypes, serviceLevel.Name)
		}

		responseMappings = append(responseMappings, &object)
	}

	response := psmpb.PartnerTypeListResponse{
		PartnerServiceTypeMappings: responseMappings,
	}

	return &response, nil
}

func (psm *PartnerServiceMappingService) updateUserStatus(ctx context.Context, supplier models.Supplier, supplierId uint64) (string, error) {
	var status models.SupplierStatus
	if supplier.Status == models.SupplierStatusVerified || supplier.Status == models.SupplierStatusFailed {
		status = models.SupplierStatusPending // Moving to Pending if any data is updated
	}

	err1 := database.DBAPM(ctx).Model(&supplier).Updates(models.Supplier{
		Status: status,
	})

	if err1 != nil && err1.Error != nil {
		return fmt.Sprintf("Error while updating user status: %s", err1.Error), nil
	} else {
		associatedPartners := []models.PartnerServiceMapping{}
		database.DBAPM(ctx).Model(&[]models.PartnerServiceMapping{}).Where("supplier_id = ?", supplierId).Scan(&associatedPartners)

		err2 := database.DBAPM(ctx).Model(&associatedPartners).UpdateColumn("active", false)

		if err2 != nil && err2.Error != nil {
			return fmt.Sprintf("Error while updating service status: %s", err2.Error), nil
		}
	}

	return "", nil
}
