package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	aaaModels "github.com/voonik/goFramework/pkg/aaa/models"
	"github.com/voonik/ss2/internal/app/utils"
	"google.golang.org/grpc/metadata"
)

type IdentityUserObject struct {
	Id    uint64
	Name  string
	Email string
	Phone string
	Roles []string
}

type IdentityUserResponse struct {
	Data *IdentityUserObject
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

func GetIdentityUser(ctx context.Context, phone string) *IdentityUserObject {
	log.Printf("GetIdentityUser: phone = %s\n", phone)
	if utils.IsEmptyStr(phone) {
		return nil
	}

	return getIdentityUserApiHelperInstance().GetUserDetailsApiByPhone(ctx, phone)
}

type IdentityUserApiHelper struct{}

type IdentityUserApiHelperInterface interface {
	IdentityBulkUserDetailsApi(context.Context, []string) map[string]IdentityUserObject
	GetUserDetailsApiByPhone(context.Context, string) *IdentityUserObject
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
	url := getIdentityUrl(ctx, "v0/users/bulk")
	if len(userIds) > 0 {
		url += "?userIds=" + strings.Join(userIds, ",")
	}

	headers := getHeaders(ctx)
	resp, _ := GetApiCallHelperInstance().Get(ctx, url, headers)

	var respData IdentityBulkUserResponse
	_ = json.Unmarshal([]byte(resp.Body), &respData)

	userDetails := make(map[string]IdentityUserObject)
	for _, identityUserObject := range respData.Data.Users {
		userDetails[fmt.Sprintf("%d", identityUserObject.Id)] = identityUserObject
	}

	return userDetails
}

func (apiHelper *IdentityUserApiHelper) GetUserDetailsApiByPhone(ctx context.Context, phone string) *IdentityUserObject {
	url := getIdentityUrl(ctx, "v0/users/phone/"+phone+"?include=[roles]")
	headers := getHeaders(ctx)
	resp, _ := GetApiCallHelperInstance().Get(ctx, url, headers)

	var respData IdentityUserResponse
	_ = json.Unmarshal([]byte(resp.Body), &respData)

	return respData.Data
}

func getIdentityUrl(ctx context.Context, suffix string) string {
	return getIdentityServiceDomain(ctx) + getIdentityServicePrefix(ctx) + suffix
}

func getHeaders(ctx context.Context) map[string]string {
	headers := make(map[string]string)
	reqHeaders, _ := metadata.FromIncomingContext(ctx)
	headers["Authorization"] = reqHeaders["authorization"][0]
	return headers
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
