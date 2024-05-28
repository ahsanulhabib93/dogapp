package helpers

import (
	"context"

	aaaModels "github.com/voonik/goFramework/pkg/aaa/models"
	"github.com/voonik/ss2/internal/app/utils"
)

type PartnerServiceEditEntity struct {
	ServiceType  utils.ServiceType
	ServiceLevel utils.SupplierType
}

func ParseServiceLevels(serviceLevels []string) []utils.SupplierType {
	values := []utils.SupplierType{}
	for _, value := range serviceLevels {
		res := utils.PartnerServiceLevelMapping[value]
		values = append(values, res)
	}
	return values
}

func ValidatePartnerSericeEdit(
	ctx context.Context,
	updated PartnerServiceEditEntity,
	existing PartnerServiceEditEntity,
) bool {
	if updated.ServiceType != existing.ServiceType {
		return false
	}
	if updated.ServiceLevel == existing.ServiceLevel {
		return true
	}
	allowedServiceLevels := aaaModels.GetAppPreferenceServiceInstance().GetValue(ctx, "edit_allowed_service_levels", []string{}).([]string)
	parsedServiceLevels := ParseServiceLevels(allowedServiceLevels)
	return utils.Includes(parsedServiceLevels, existing.ServiceLevel) && utils.Includes(parsedServiceLevels, updated.ServiceLevel)
}
