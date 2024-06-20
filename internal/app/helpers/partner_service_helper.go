package helpers

import (
	"context"
	"fmt"

	aaaModels "github.com/voonik/goFramework/pkg/aaa/models"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/utils"
)

type PartnerServiceEditEntity struct {
	ServiceType    utils.ServiceType
	ServiceLevelId uint64
}

func ValidatePartnerSericeEdit(
	ctx context.Context,
	updated PartnerServiceEditEntity,
	existing PartnerServiceEditEntity,
) bool {
	if updated.ServiceType != existing.ServiceType {
		return false
	}
	if updated.ServiceLevelId == existing.ServiceLevelId {
		return true
	}
	allowedServiceLevels := aaaModels.GetAppPreferenceServiceInstance().GetValue(ctx, "edit_allowed_service_levels", []string{}).([]string)
	return utils.Includes(allowedServiceLevels, fmt.Sprint(existing.ServiceLevelId)) && utils.Includes(allowedServiceLevels, fmt.Sprint(updated.ServiceLevelId))
}

func GetServiceTypeLevelMappings(ctx context.Context) map[utils.ServiceType][]models.PartnerServiceLevel {
	serviceTypeLevelMapping := make(map[utils.ServiceType][]models.PartnerServiceLevel)
	var serviceTypeLevels []models.PartnerServiceLevel
	database.DBAPM(ctx).Model(&models.PartnerServiceLevel{}).Find(&serviceTypeLevels)

	for _, serviceLevel := range serviceTypeLevels {
		serviceTypeLevelMapping[serviceLevel.ServiceType] = append(serviceTypeLevelMapping[serviceLevel.ServiceType], serviceLevel)
	}
	return serviceTypeLevelMapping
}

func GetServiceLevelByTypeAndName(ctx context.Context, serviceType utils.ServiceType, name string) models.PartnerServiceLevel {
	var serviceLevel models.PartnerServiceLevel
	database.DBAPM(ctx).Model(&models.PartnerServiceLevel{}).Where("service_type = ? and name = ?", serviceType, name).First(&serviceLevel)
	return serviceLevel
}

func GetServiceLevelById(ctx context.Context, id int) models.PartnerServiceLevel {
	var serviceLevel models.PartnerServiceLevel
	database.DBAPM(ctx).Model(&models.PartnerServiceLevel{}).First(&serviceLevel, id)
	return serviceLevel
}

func GetServiceLevelNameMapping(ctx context.Context) map[uint64]string {
	serviceLevelNameMapping := make(map[uint64]string)
	var serviceLevels []models.PartnerServiceLevel
	database.DBAPM(ctx).Model(&models.PartnerServiceLevel{}).Find(&serviceLevels)

	for _, serviceLevel := range serviceLevels {
		serviceLevelNameMapping[serviceLevel.ID] = serviceLevel.Name
	}
	return serviceLevelNameMapping
}

func GetServiceLevelIdMapping(ctx context.Context) map[string]uint64 {
	serviceLevelIDMapping := make(map[string]uint64)
	var serviceLevels []models.PartnerServiceLevel
	database.DBAPM(ctx).Model(&models.PartnerServiceLevel{}).Find(&serviceLevels)

	for _, serviceLevel := range serviceLevels {
		serviceLevelIDMapping[serviceLevel.Name] = serviceLevel.ID
	}
	return serviceLevelIDMapping
}
