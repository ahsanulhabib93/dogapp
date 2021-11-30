package test_helper

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/utils"
)

func CreateSupplier(ctx context.Context, supplier *models.Supplier) *models.Supplier {
	id := rand.Intn(100)
	supplier.Email = fmt.Sprintf("test-%v@shopup.org", id)
	if supplier.Name == "" {
		supplier.Name = fmt.Sprintf("Test-%v", id)
	}
	if supplier.SupplierType == 0 {
		supplier.SupplierType = utils.Hlc
	}
	database.DBAPM(ctx).Save(supplier)
	return supplier
}

func CreateSupplierWithAddress(ctx context.Context, supplier *models.Supplier) *models.Supplier {
	supplier = CreateSupplier(ctx, supplier)
	CreateSupplierAddress(ctx, &models.SupplierAddress{Supplier: *supplier, IsDefault: true})
	return supplier
}

func CreateSupplierAddress(ctx context.Context, supplierAddress *models.SupplierAddress) *models.SupplierAddress {
	id := rand.Intn(100)
	supplierAddress.Firstname = fmt.Sprintf("Firstname-%v", id)
	supplierAddress.Lastname = fmt.Sprintf("Lastname-%v", id)
	supplierAddress.Address1 = fmt.Sprintf("Address1-%v", id)
	supplierAddress.Address2 = fmt.Sprintf("Address2-%v", id)
	supplierAddress.Landmark = fmt.Sprintf("Landmark-%v", id)
	supplierAddress.City = fmt.Sprintf("City-%v", id)
	supplierAddress.State = fmt.Sprintf("State-%v", id)
	supplierAddress.Country = fmt.Sprintf("Country-%v", id)
	supplierAddress.Zipcode = fmt.Sprintf("Zipcode-%v", id)
	supplierAddress.Phone = fmt.Sprintf("Phone-%v", id)
	supplierAddress.GstNumber = fmt.Sprintf("GstNumber-%v", id)
	database.DBAPM(ctx).Save(supplierAddress)
	return supplierAddress
}

func CreatePaymentAccountDetail(ctx context.Context, paymentAccount *models.PaymentAccountDetail) *models.PaymentAccountDetail {
	id := rand.Intn(100)
	paymentAccount.AccountName = fmt.Sprintf("AccountName-%v", id)
	paymentAccount.AccountNumber = fmt.Sprintf("AccountNumber-%v", id)

	if paymentAccount.AccountType != utils.Mfs {
		paymentAccount.AccountType = utils.Bank
		paymentAccount.BankName = fmt.Sprintf("BankName-%v", id)
		paymentAccount.BranchName = fmt.Sprintf("BranchName-%v", id)
		paymentAccount.RoutingNumber = fmt.Sprintf("RoutingNumber-%v", id)
	}
	database.DBAPM(ctx).Save(paymentAccount)
	return paymentAccount
}

func CreateKeyAccountManager(ctx context.Context, accountManager *models.KeyAccountManager) *models.KeyAccountManager {
	id := rand.Intn(100)
	accountManager.Name = fmt.Sprintf("Test-%v", id)
	accountManager.Email = fmt.Sprintf("test-%v@shopup.org", id)
	accountManager.Email = fmt.Sprintf("Phone-%v", id)
	database.DBAPM(ctx).Save(accountManager)
	return accountManager
}
