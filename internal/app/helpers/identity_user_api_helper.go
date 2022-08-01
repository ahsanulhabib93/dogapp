package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	aaaModels "github.com/voonik/goFramework/pkg/aaa/models"
	"google.golang.org/grpc/metadata"
)

type IdentityUserObject struct {
	Id    uint64
	Name  string
	Email string
	Phone string
}

type IdentityBulkUserResponse struct {
	Data IdentityBulkUserData
}

type IdentityBulkUserData struct {
	Users []IdentityUserObject
}

func IdentityBulkUserApi(ctx context.Context, userIds []string) map[string]IdentityUserObject {
	return getIdentityUserApiHelperInstance().IdentityBulkUserDetailsApi(ctx, userIds)
}

type IdentityUserApiHelper struct{}

type IdentityUserApiHelperInterface interface {
	IdentityBulkUserDetailsApi(context.Context, []string) map[string]IdentityUserObject
}

var identityUserApiHelper IdentityUserApiHelperInterface

func InjectMockIdentityUserApiHelperInstance(mockObj IdentityUserApiHelperInterface) {
	identityUserApiHelper = mockObj
}

func getIdentityUserApiHelperInstance() IdentityUserApiHelperInterface {
	if identityUserApiHelper == nil {
		return new(IdentityUserApiHelper)
	}
	return identityUserApiHelper
}

func (apiHelper *IdentityUserApiHelper) IdentityBulkUserDetailsApi(ctx context.Context, userIds []string) map[string]IdentityUserObject {
	url := getIdentityServiceDomain(ctx) + getIdentityServicePrefix(ctx) + "v0/users/bulk"
	if len(userIds) > 0 {
		url += "?userIds=" + strings.Join(userIds, ",")
	}

	headers := make(map[string]string)
	reqHeaders, _ := metadata.FromIncomingContext(ctx)
	headers["Authorization"] = reqHeaders["authorization"][0]
	resp, _ := GetApiCallHelperInstance().Get(ctx, url, headers)

	var respData IdentityBulkUserResponse
	_ = json.Unmarshal([]byte(resp.Body), &respData)

	userDetails := make(map[string]IdentityUserObject)
	for _, identityUserObject := range respData.Data.Users {
		userDetails[fmt.Sprintf("%d", identityUserObject.Id)] = identityUserObject
	}

	return userDetails
}

func getIdentityServiceDomain(ctx context.Context) string {
	domain := aaaModels.AppPreference.GetValue(
		aaaModels.AppPreference{}, ctx, "identity_service_domain", "https://authfe.shopups2.xyz",
	)
	return domain.(string)
}

func getIdentityServicePrefix(ctx context.Context) string {
	prefix := aaaModels.AppPreference.GetValue(
		aaaModels.AppPreference{}, ctx, "identity_service_prefix", "/api/auth/",
	)
	return prefix.(string)
}
