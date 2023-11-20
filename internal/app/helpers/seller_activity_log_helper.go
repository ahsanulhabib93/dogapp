package helpers

import (
	"context"

	"github.com/shopuptech/go-libs/logger"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/utils"
)

func AddSellerActivityLog(ctx context.Context, log models.SellerActivityLog) error {
	if uid := utils.GetCurrentUserID(ctx); uid != nil {
		log.UserID = *uid
	}
	err := database.DBAPM(ctx).Model(&models.SellerActivityLog{}).Create(&log).Error
	if err != nil {
		logger.FromContext(ctx).Error("error while inserting into seller activity log: ", err.Error())
		return err
	}
	return nil
}
