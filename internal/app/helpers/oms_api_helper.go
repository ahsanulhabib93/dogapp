package helpers

import (
	"context"

	userMappingPb "github.com/voonik/goConnect/api/go/oms/user_mapping"
	userMappingService "github.com/voonik/goConnect/oms/user_mapping"
)

type OmsApiHelper struct{}

type OmsApiHelperInterface interface {
	FetchUserMappingData(ctx context.Context, userIDs []uint64) (*userMappingPb.UserMappingResponse, error)
}

var omsApiHelper OmsApiHelperInterface

func InjectMockOmsAPIHelperInstance(mockObj OmsApiHelperInterface) {
	omsApiHelper = mockObj
}

func GetOmsAPIHelperInstance() OmsApiHelperInterface {
	if omsApiHelper == nil {
		return new(OmsApiHelper)
	}
	return omsApiHelper
}

func (omsApiHelper *OmsApiHelper) FetchUserMappingData(ctx context.Context, userIDs []uint64) (*userMappingPb.UserMappingResponse, error) {
	return userMappingService.OMSUserMapping().Index(ctx, &userMappingPb.UserMappingIndexParams{
		UserId: userIDs,
	})
}
