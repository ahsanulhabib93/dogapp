package helpers

import (
	"context"

	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/models"
)

func UpdateOtherAddress(ctx context.Context, address *models.SupplierAddress) error {
	if address.IsDefault {
		otherDefaultAddress := &models.SupplierAddress{}
		database.DBAPM(ctx).Model(address).Where("supplier_id = ? and is_default = ? and id != ?", address.SupplierID, true, address.ID).First(&otherDefaultAddress)
		if otherDefaultAddress != nil {
			database.DBAPM(ctx).Model(otherDefaultAddress).Update("is_default", false)
		}
	}
	return nil
}
