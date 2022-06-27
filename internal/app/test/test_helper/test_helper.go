package test_helper

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/google/uuid"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/utils"
)

func getUniqueID() string {
	id := uuid.New()
	return id.String()
}

func CreateSupplier(ctx context.Context, supplier *models.Supplier) *models.Supplier {
	id := getUniqueID()

	supplier.Email = fmt.Sprintf("test-%v@shopup.org", id)
	supplier.AlternatePhone = fmt.Sprintf("8801234567890%v", id)
	supplier.BusinessName = fmt.Sprintf("Test Business %v", id)
	supplier.Reason = fmt.Sprintf("Test reason %v", id)
	supplier.ShopImageURL = fmt.Sprintf("/ss2/test_shop_image_url/%v", id)

	if supplier.Name == "" {
		supplier.Name = fmt.Sprintf("Test-%v", id)
	}
	if supplier.SupplierType == 0 {
		supplier.SupplierType = utils.Hlc
	}
	if supplier.Status == "" {
		supplier.Status = models.SupplierStatusPending
	}
	if supplier.Phone == "" {
		supplier.Phone = fmt.Sprintf("8801%v", id[:9])
	}

	database.DBAPM(ctx).Save(supplier)
	return supplier
}

func CreateSupplierWithAddress(ctx context.Context, supplier *models.Supplier) *models.Supplier {
	supplier = CreateSupplier(ctx, supplier)
	CreateSupplierAddress(ctx, &models.SupplierAddress{SupplierID: supplier.ID, IsDefault: true})
	return supplier
}

func CreateSupplierAddress(ctx context.Context, supplierAddress *models.SupplierAddress) *models.SupplierAddress {
	id := getUniqueID()
	supplierAddress.Firstname = fmt.Sprintf("Firstname-%v", id)
	supplierAddress.Lastname = fmt.Sprintf("Lastname-%v", id)
	supplierAddress.Address1 = fmt.Sprintf("Address1-%v", id)
	supplierAddress.Address2 = fmt.Sprintf("Address2-%v", id)
	supplierAddress.Landmark = fmt.Sprintf("Landmark-%v", id)
	supplierAddress.City = fmt.Sprintf("City-%v", id)
	supplierAddress.State = fmt.Sprintf("State-%v", id)
	supplierAddress.Country = fmt.Sprintf("Country-%v", id)
	supplierAddress.Zipcode = fmt.Sprintf("Zipcode-%v", id)
	supplierAddress.Phone = fmt.Sprintf("Phone-%v", rand.Intn(13))
	supplierAddress.GstNumber = fmt.Sprintf("GstNumber-%v", id)
	database.DBAPM(ctx).Save(supplierAddress)
	return supplierAddress
}

func CreatePaymentAccountDetail(ctx context.Context, paymentAccount *models.PaymentAccountDetail) *models.PaymentAccountDetail {
	id := getUniqueID()
	paymentAccount.AccountName = fmt.Sprintf("AccountName-%v", id)
	paymentAccount.AccountNumber = fmt.Sprintf("AccountNumber-%v", id)

	if paymentAccount.AccountType == utils.Mfs {
		paymentAccount.AccountSubType = utils.Bkash
	} else {
		paymentAccount.AccountType = utils.Bank
		paymentAccount.AccountSubType = utils.Current
		if paymentAccount.BankID == 0 {
			bank := CreateBank(ctx, &models.Bank{})
			paymentAccount.BankID = bank.ID
		}
		paymentAccount.BranchName = fmt.Sprintf("BranchName-%v", id)
		paymentAccount.RoutingNumber = fmt.Sprintf("RoutingNumber-%v", id)
	}

	database.DBAPM(ctx).Save(paymentAccount)
	return paymentAccount
}

func CreateKeyAccountManager(ctx context.Context, accountManager *models.KeyAccountManager) *models.KeyAccountManager {
	id := getUniqueID()
	accountManager.Name = fmt.Sprintf("Test-%v", id)
	accountManager.Email = fmt.Sprintf("test-%v@shopup.org", id)
	accountManager.Email = fmt.Sprintf("Phone-%v", id)
	database.DBAPM(ctx).Save(accountManager)
	return accountManager
}

func CreateBank(ctx context.Context, bank *models.Bank) *models.Bank {
	id := getUniqueID()
	bank.Name = fmt.Sprintf("TestBank-%v", id)
	database.DBAPM(ctx).Save(bank)
	return bank
}
