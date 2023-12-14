package helpers

import (
	"context"

	sampb "github.com/voonik/goConnect/api/go/ss2/seller_account_manager"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/models"
)

func GetAndFormatSellerAccountManager(ctx context.Context, sellerID uint64) []*sampb.AccountManagerObject {
	var samList []models.SellerAccountManager
	var accountManagers []*sampb.AccountManagerObject
	database.DBAPM(ctx).Model(&models.SellerAccountManager{}).Where(`seller_id =?`, sellerID).Order("role, priority").Scan(&samList)

	for _, sam := range samList {
		accountManagers = append(accountManagers, &sampb.AccountManagerObject{
			Id:       sam.ID,
			Email:    sam.Email,
			Phone:    uint64(sam.Phone),
			Name:     sam.Name,
			Priority: uint64(sam.Priority),
			Role:     sam.Role,
		})
	}
	return accountManagers
}
