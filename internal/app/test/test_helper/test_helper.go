package test_helper

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/goFramework/pkg/misc"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/utils"
)

func getUniqueID() string {
	id := uuid.New()
	return id.String()
}

func CreateSupplier(ctx context.Context, supplier *models.Supplier) *models.Supplier {
	id := getUniqueID()

	partnerServiceMapping := models.PartnerServiceMapping{}
	if len(supplier.PartnerServiceMappings) != 0 {
		partnerServiceMapping = supplier.PartnerServiceMappings[0]
	}
	partnerServiceMapping.ServiceType = utils.Supplier
	partnerServiceMapping.Active = true

	supplier.Email = fmt.Sprintf("test-%v@shopup.org", id)
	supplier.AlternatePhone = fmt.Sprintf("8801234567890%v", id)
	supplier.BusinessName = fmt.Sprintf("Test Business %v", id)
	supplier.Reason = fmt.Sprintf("Test reason %v", id)
	supplier.ShopImageURL = fmt.Sprintf("/ss2/test_shop_image_url/%v", id)

	if supplier.Name == "" {
		supplier.Name = fmt.Sprintf("Test-%v", id)
	}
	if supplier.Status == "" {
		supplier.Status = models.SupplierStatusPending
	}
	if supplier.Phone == "" {
		supplier.Phone = fmt.Sprintf("8801%v", id[:9])
	}

	if supplier.SupplierType != 0 {
		partnerServiceMapping.ServiceLevel = supplier.SupplierType
	} else if partnerServiceMapping.ServiceLevel == 0 {
		partnerServiceMapping.ServiceLevel = utils.Hlc
	}

	supplier.PartnerServiceMappings = []models.PartnerServiceMapping{partnerServiceMapping}
	database.DBAPM(ctx).Save(supplier)
	return supplier
}

func CreatePartnerServiceMapping(ctx context.Context, partnerServiceMapping *models.PartnerServiceMapping) *models.PartnerServiceMapping {
	id := getUniqueID()
	partnerServiceMapping.TradeLicenseUrl = fmt.Sprintf("trade_license_url_%v", id)
	partnerServiceMapping.AgreementUrl = fmt.Sprintf("agreement_url_%v", id)

	database.DBAPM(ctx).Save(partnerServiceMapping)
	return partnerServiceMapping
}

func CreateSupplierWithDateTime(ctx context.Context, supplier *models.Supplier, createAt time.Time) *models.Supplier {
	supplier.CreatedAt = createAt
	supplier.UpdatedAt = createAt
	return CreateSupplier(ctx, supplier)
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
	supplierAddress.Phone = fmt.Sprintf("880%v", rand.Intn(999999999)+1000000000)
	supplierAddress.GstNumber = fmt.Sprintf("GstNumber-%v", id)
	database.DBAPM(ctx).Save(supplierAddress)
	return supplierAddress
}

func CreatePaymentAccountDetail(ctx context.Context, paymentAccount *models.PaymentAccountDetail) *models.PaymentAccountDetail {
	id := getUniqueID()
	paymentAccount.AccountName = fmt.Sprintf("AccountName-%v", id)
	number := paymentAccount.AccountNumber
	if number == "" {
		paymentAccount.AccountNumber = fmt.Sprintf("AccountNumber-%v", id)
	} else {
		paymentAccount.AccountNumber = number
	}

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

func CreatePaymentAccountDetailWarehouseMappings(ctx context.Context, paymentAccountDetailWarehouseMapping *models.PaymentAccountDetailWarehouseMapping) *models.PaymentAccountDetailWarehouseMapping {
	if paymentAccountDetailWarehouseMapping.WarehouseID == 0 {
		paymentAccountDetailWarehouseMapping.WarehouseID = rand.Uint64()
	}
	if paymentAccountDetailWarehouseMapping.PaymentAccountDetailID == 0 {
		paymentAccountDetailWarehouseMapping.PaymentAccountDetailID = rand.Uint64()
	}
	database.DBAPM(ctx).Save(paymentAccountDetailWarehouseMapping)
	return paymentAccountDetailWarehouseMapping
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

func SetContextUser(ctx context.Context, userId uint64, permissions []string) context.Context {
	threadObject := &misc.ThreadObject{
		VaccountId:    1,
		PortalId:      1,
		CurrentActId:  1,
		XForwardedFor: "5079327",
		UserData: &misc.UserData{
			UserId:      userId,
			Name:        "John",
			Email:       "john@gmail.com",
			Phone:       "8801855533367",
			Permissions: permissions,
		},
	}
	return misc.SetInContextThreadObject(ctx, threadObject)
}
