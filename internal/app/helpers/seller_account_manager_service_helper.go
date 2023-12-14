package helpers

import (
	"context"

	sampb "github.com/voonik/goConnect/api/go/ss2/seller_account_manager"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/models"
)

func GetAndFormatSellerAccountManager(ctx context.Context, sellerID uint64) *sampb.AccountManagerObject {
	var Sam models.SellerAccountManager
	database.DBAPM(ctx).Model(&models.SellerAccountManager{}).Where(`seller_id =?`, sellerID).First(&Sam)

	accountManager := &sampb.AccountManagerObject{
		Id:       Sam.ID,
		Email:    Sam.Email,
		Phone:    uint64(Sam.Phone),
		Name:     Sam.Name,
		Priority: uint64(Sam.Priority),
		Role:     Sam.Role,
	}
	return accountManager
}
