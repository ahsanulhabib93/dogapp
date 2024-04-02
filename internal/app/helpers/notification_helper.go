package helpers

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/voonik/goConnect/api/go/vigeon/notify"
	vigeon "github.com/voonik/goConnect/vigeon/notify"
	aaaModels "github.com/voonik/goFramework/pkg/aaa/models"

	"github.com/voonik/ss2/internal/app/models"
)

func SendStatusChangeEmailNotification(ctx context.Context, supplier models.Supplier, status, reason string) *notify.EmailResp {
	if supplier.UserID == nil {
		return nil
	}

	sUserId := strconv.Itoa(int(*supplier.UserID))
	user, found := IdentityBulkUserApi(ctx, []string{sUserId})[sUserId]
	if !found {
		return nil
	}

	fromEmail := aaaModels.AppPreference.GetValue(
		aaaModels.AppPreference{}, ctx, "update_status_from_email", "ss2@shopf.co").(string)
	subject := aaaModels.AppPreference.GetValue(
		aaaModels.AppPreference{}, ctx, "update_status_subject", "User Status Changed").(string)
	content := aaaModels.AppPreference.GetValue(
		aaaModels.AppPreference{}, ctx, "update_status_content", "Users($id) status has been updated to \"$status\". Reason: $reason").(string)

	subject = strings.ReplaceAll(subject, "$id", fmt.Sprint(supplier.ID))
	subject = strings.ReplaceAll(subject, "$status", strings.Title(status))

	content = strings.ReplaceAll(content, "$id", fmt.Sprint(supplier.ID))
	content = strings.ReplaceAll(content, "$status", strings.Title(status))
	content = strings.ReplaceAll(content, "$reason", reason)

	emailParam := notify.EmailParam{
		ToEmail:   user.Email,
		FromEmail: fromEmail,
		Subject:   subject,
		Content:   content,
	}

	return getVigeonAPIHelperInstance().SendEmailAPI(ctx, emailParam) //nolint:govet
}

type VigeonAPIHelper struct{}

type VigeonAPIHelperInterface interface {
	SendEmailAPI(ctx context.Context, emailParam notify.EmailParam) *notify.EmailResp
}

var vigeonApiHelper VigeonAPIHelperInterface

// InjectMockVigeonAPIHelperInstance ...
func InjectMockVigeonAPIHelperInstance(mockObj VigeonAPIHelperInterface) {
	vigeonApiHelper = mockObj
}

// getVigeonAPIHelperInstance ...
func getVigeonAPIHelperInstance() VigeonAPIHelperInterface {
	if vigeonApiHelper == nil {
		vigeonApiHelper = new(VigeonAPIHelper)
	}
	return vigeonApiHelper
}

func (apiHelper *VigeonAPIHelper) SendEmailAPI(ctx context.Context, emailParam notify.EmailParam) *notify.EmailResp { //nolint:govet
	resp, err := vigeon.Notify().EmailNotification(ctx, &emailParam)
	if err != nil {
		log.Println("SentEmailAPI: Failed to sent email. Error: ", err.Error())
	}

	return resp
}
