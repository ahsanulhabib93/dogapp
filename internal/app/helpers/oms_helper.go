package helpers

import (
	"context"
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

func FetchFormattedUserMappingData(ctx context.Context, userIDs []uint64) (map[uint64]UserMapping, error) {
	omsResp, err := GetOmsAPIHelperInstance().FetchUserMappingData(ctx, userIDs)
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
