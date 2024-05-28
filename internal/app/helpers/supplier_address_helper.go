package helpers

import (
	"context"

	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/models"
)

func UpdateDefaultAddress(ctx context.Context, address *models.SupplierAddress) {
	otherDefaultAddress := &models.SupplierAddress{}
	err := database.DBAPM(ctx).Model(address).Where("supplier_id = ? and is_default = ? and id != ?", address.SupplierID, true, address.ID).First(&otherDefaultAddress).Error
	if err != nil && err.Error() == "record not found" {
		database.DBAPM(ctx).Model(address).Update("is_default", true)
	}

	if address.IsDefault && otherDefaultAddress != nil {
		database.DBAPM(ctx).Model(otherDefaultAddress).Update("is_default", false)
	}
}

func UpdateDefaultPaymentAccount(ctx context.Context, paymentAccount *models.PaymentAccountDetail) {
	if paymentAccount.IsDefault {
		otherDefaultPayment := &models.PaymentAccountDetail{}
		database.DBAPM(ctx).Model(paymentAccount).Where("supplier_id = ? and is_default = ? and id != ?", paymentAccount.SupplierID, true, paymentAccount.ID).First(&otherDefaultPayment)
		if otherDefaultPayment != nil {
			database.DBAPM(ctx).Model(otherDefaultPayment).Update("is_default", false)
		}
	}
}
