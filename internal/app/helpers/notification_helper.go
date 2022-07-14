package helpers

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/voonik/goConnect/api/go/vigeon/notify"
	vigeon "github.com/voonik/goConnect/vigeon/notify"

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

	emailParam := notify.EmailParam{
		ToEmail:   user.Email,
		FromEmail: "ss2@shopup.org",
		Subject:   "User Status Changed",
		Content:   fmt.Sprintf("User(#%v) status has been updated to \"%v\"", supplier.ID, strings.Title(status)),
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
