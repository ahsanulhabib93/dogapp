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

func SendStatusChangeEmailNotification(ctx context.Context, supplier models.Supplier, status string) *notify.EmailResp {
	if supplier.UserID == nil {
		return nil
	}

	sUserId := strconv.Itoa(int(*supplier.UserID))
	user, found := IdentityBulkUserApi(ctx, []string{sUserId})[sUserId]
	if !found {
		return nil
	}

	fromEmail := aaaModels.AppPreference.GetValue(
		aaaModels.AppPreference{}, ctx, "update_status_from_email", "ss2@shopup.org").(string)
	subject := aaaModels.AppPreference.GetValue(
		aaaModels.AppPreference{}, ctx, "update_status_subject", "User Status Changed").(string)
	content := aaaModels.AppPreference.GetValue(
		aaaModels.AppPreference{}, ctx, "update_status_content",
		fmt.Sprintf("Users(#%v) status has been updated to \"%v\"", supplier.ID, strings.Title(status))).(string)

	emailParam := notify.EmailParam{
		ToEmail:   user.Email,
		FromEmail: fromEmail,
		Subject:   subject,
		Content:   content,
	}

	return getVigeonAPIHelperInstance().SendEmailAPI(ctx, emailParam)
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

func (apiHelper *VigeonAPIHelper) SendEmailAPI(ctx context.Context, emailParam notify.EmailParam) *notify.EmailResp {
	resp, err := vigeon.Notify().EmailNotification(ctx, &emailParam)
	if err != nil {
		log.Println("SentEmailAPI: Failed to sent email. Error: ", err.Error())
	}

	return resp
}
