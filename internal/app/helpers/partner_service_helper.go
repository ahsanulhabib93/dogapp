package helpers

import (
	"context"

	aaaModels "github.com/voonik/goFramework/pkg/aaa/models"
	"github.com/voonik/ss2/internal/app/utils"
)

func ParseServiceLevels(serviceLevels []string) []utils.SupplierType {
	values := []utils.SupplierType{}
	for _, value := range serviceLevels {
		res := utils.PartnerServiceLevelMapping[value]
		values = append(values, res)
	}
	return values
}

func ValidateSericeLevelEdit(
	ctx context.Context,
	updated utils.SupplierType,
	existing utils.SupplierType,
) bool {
	if updated == existing {
		return true
	}
	allowedServiceLevels := aaaModels.GetAppPreferenceServiceInstance().GetValue(ctx, "edit_allowed_service_levels", []string{}).([]string)
	parsedServiceLevels := ParseServiceLevels(allowedServiceLevels)
	return utils.Includes(parsedServiceLevels, existing) && utils.Includes(parsedServiceLevels, updated)
}
