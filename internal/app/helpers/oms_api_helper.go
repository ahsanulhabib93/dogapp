package helpers

import (
	"context"

	userMappingPb "github.com/voonik/goConnect/api/go/oms/user_mapping"
	userMappingService "github.com/voonik/goConnect/oms/user_mapping"
)

type UserMapping struct {
	BusinessUnits []uint64
	ZoneIDs       []uint64
	OpcIDs        []uint64
	UserData      ISUserData
}

type ISUserData struct {
	email string
	name  string
	roles []uint64
}

func FetchUserMappingData(ctx context.Context, userIDs []uint64) (map[uint64]UserMapping, error) {
	omsResp, err := userMappingService.OMSUserMapping().Index(ctx, &userMappingPb.UserMappingIndexParams{
		UserId: userIDs,
	})
	if err != nil {
		return nil, err
	}
	userMappingMap := make(map[uint64]UserMapping)
	for _, data := range omsResp.GetData() {
		mappingData := UserMapping{
			BusinessUnits: data.GetBusinessUnits(),
			OpcIDs:        data.GetOpcIds(),
			ZoneIDs:       data.GetZoneIds(),
		}

		if data.GetUserData() != nil {
			mappingData.UserData = ISUserData{
				email: data.GetUserData().GetEmail(),
				name:  data.GetUserData().GetName(),
				roles: data.GetUserData().GetRoles(),
			}
		}
		userMappingMap[data.GetUserId()] = mappingData
	}
	return userMappingMap, nil
}
