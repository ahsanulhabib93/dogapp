package payment_account_detail_service_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	paymentpb "github.com/voonik/goConnect/api/go/ss2/payment_account_detail"
	aaaModels "github.com/voonik/goFramework/pkg/aaa/models"
	"github.com/voonik/goFramework/pkg/database"
	test_utils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/services"
	"github.com/voonik/ss2/internal/app/test/mocks"
	"github.com/voonik/ss2/internal/app/test/test_helper"
	"github.com/voonik/ss2/internal/app/utils"
)

var _ = Describe("EditPaymentAccountDetail", func() {
	var ctx context.Context

	BeforeEach(func() {
		test_utils.GetContext(&ctx)
		test_utils.SetPermission(&ctx, []string{"supplierpanel:editverifiedblockedsupplieronly:admin"})
		aaaModels.CreateAppPreferenceServiceInterface()
	})

	Context("Editing all attributes of existing PaymentAccount", func() {
		It("Should update and return success response", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			paymentAccount := test_helper.CreatePaymentAccountDetail(ctx, &models.PaymentAccountDetail{SupplierID: supplier.ID, AccountType: utils.Bank, IsDefault: true})
			bank := test_helper.CreateBank(ctx, &models.Bank{})
			param := &paymentpb.PaymentAccountDetailObject{
				Id:             paymentAccount.ID,
				AccountType:    uint64(utils.Bank),
				AccountSubType: uint64(utils.Savings),
				AccountName:    "AccountName",
				AccountNumber:  "AccountNumber",
				BankId:         bank.ID,
				BranchName:     "BranchName",
				RoutingNumber:  "RoutingNumber",
				IsDefault:      true,
			}
			res, err := new(services.PaymentAccountDetailService).Edit(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("PaymentAccountDetail Edited Successfully"))

			database.DBAPM(ctx).Model(&models.PaymentAccountDetail{}).First(&paymentAccount, paymentAccount.ID)
			Expect(paymentAccount.AccountType).To(Equal(utils.Bank))
			Expect(paymentAccount.AccountSubType).To(Equal(utils.Savings))
			Expect(paymentAccount.AccountName).To(Equal(param.AccountName))
			Expect(paymentAccount.AccountNumber).To(Equal(param.AccountNumber))
			Expect(paymentAccount.BankID).To(Equal(param.BankId))
			Expect(paymentAccount.BranchName).To(Equal(param.BranchName))
			Expect(paymentAccount.RoutingNumber).To(Equal(param.RoutingNumber))
			Expect(paymentAccount.IsDefault).To(Equal(true))

			database.DBAPM(ctx).Model(&models.Supplier{}).First(&supplier, supplier.ID)
			Expect(supplier.Status).To(Equal(models.SupplierStatusPending))
		})
	})

	Context("Updating non-default address as default", func() {
		It("Should update other default address as non-default and return success response", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{Status: models.SupplierStatusFailed})
			paymentAccount1 := test_helper.CreatePaymentAccountDetail(ctx, &models.PaymentAccountDetail{SupplierID: supplier.ID, IsDefault: true})
			paymentAccount2 := test_helper.CreatePaymentAccountDetail(ctx, &models.PaymentAccountDetail{SupplierID: supplier.ID, IsDefault: false})
			paymentAccount3 := test_helper.CreatePaymentAccountDetail(ctx, &models.PaymentAccountDetail{SupplierID: supplier.ID, IsDefault: false})
			param := &paymentpb.PaymentAccountDetailObject{
				Id:        paymentAccount3.ID,
				IsDefault: true,
			}
			res, err := new(services.PaymentAccountDetailService).Edit(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("PaymentAccountDetail Edited Successfully"))

			database.DBAPM(ctx).Model(&models.PaymentAccountDetail{}).First(&paymentAccount1, paymentAccount1.ID)
			Expect(paymentAccount1.IsDefault).To(Equal(false))
			database.DBAPM(ctx).Model(&models.PaymentAccountDetail{}).First(&paymentAccount2, paymentAccount2.ID)
			Expect(paymentAccount2.IsDefault).To(Equal(false))
			database.DBAPM(ctx).Model(&models.PaymentAccountDetail{}).First(&paymentAccount3, paymentAccount3.ID)
			Expect(paymentAccount3.IsDefault).To(Equal(true))

			database.DBAPM(ctx).Model(&models.Supplier{}).First(&supplier, supplier.ID)
			Expect(supplier.Status).To(Equal(models.SupplierStatusPending))
		})
	})

	Context("Editing only account number of existing record", func() {
		It("Should update and return success response", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			paymentAccount := test_helper.CreatePaymentAccountDetail(ctx, &models.PaymentAccountDetail{SupplierID: supplier.ID, AccountType: utils.Bank, IsDefault: true})
			param := &paymentpb.PaymentAccountDetailObject{
				Id:            paymentAccount.ID,
				AccountName:   "AccountName",
				AccountNumber: "AccountNumber",
			}
			res, err := new(services.PaymentAccountDetailService).Edit(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("PaymentAccountDetail Edited Successfully"))

			updatedPayment := &models.PaymentAccountDetail{}
			database.DBAPM(ctx).Model(&models.PaymentAccountDetail{}).First(&updatedPayment, paymentAccount.ID)
			Expect(paymentAccount.AccountType).To(Equal(utils.Bank))
			Expect(paymentAccount.AccountSubType).To(Equal(utils.Current))
			Expect(updatedPayment.AccountName).To(Equal(param.AccountName))
			Expect(updatedPayment.AccountNumber).To(Equal(param.AccountNumber))
			Expect(updatedPayment.BankID).To(Equal(paymentAccount.BankID))
			Expect(updatedPayment.BranchName).To(Equal(paymentAccount.BranchName))
			Expect(updatedPayment.RoutingNumber).To(Equal(paymentAccount.RoutingNumber))
			Expect(updatedPayment.IsDefault).To(Equal(paymentAccount.IsDefault))
		})
	})

	Context("Editing invalid payment account detail", func() {
		It("Should return error response", func() {
			param := &paymentpb.PaymentAccountDetailObject{Id: 1000}
			res, err := new(services.PaymentAccountDetailService).Edit(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("PaymentAccountDetail Not Found"))
		})
	})

	Context("Editing with invalid account sub_type", func() {
		It("Should return error response", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			paymentAccount := test_helper.CreatePaymentAccountDetail(ctx, &models.PaymentAccountDetail{SupplierID: supplier.ID, AccountType: utils.Bank, IsDefault: true})
			param := &paymentpb.PaymentAccountDetailObject{
				Id:             paymentAccount.ID,
				AccountSubType: uint64(utils.Bkash),
			}
			res, err := new(services.PaymentAccountDetailService).Edit(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Error while updating PaymentAccountDetail: Invalid Account SubType"))
		})
	})

	Context("Editing with invalid bank name", func() {
		It("Should return error response", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			paymentAccount := test_helper.CreatePaymentAccountDetail(ctx, &models.PaymentAccountDetail{SupplierID: supplier.ID, AccountType: utils.Bank, IsDefault: true})
			param := &paymentpb.PaymentAccountDetailObject{
				Id:     paymentAccount.ID,
				BankId: 100,
			}
			res, err := new(services.PaymentAccountDetailService).Edit(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Error while updating PaymentAccountDetail: Invalid Bank Name"))
		})
	})

	Context("Editing all attributes of existing PaymentAccount when supplier in verified state", func() {
		It("Should update and return success response", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{Status: models.SupplierStatusVerified})
			paymentAccount := test_helper.CreatePaymentAccountDetail(ctx, &models.PaymentAccountDetail{SupplierID: supplier.ID, AccountType: utils.Bank, IsDefault: true})
			bank := test_helper.CreateBank(ctx, &models.Bank{})
			database.DBAPM(ctx).Model(&supplier).Updates(&models.Supplier{Status: models.SupplierStatusVerified})
			param := &paymentpb.PaymentAccountDetailObject{
				Id:             paymentAccount.ID,
				AccountType:    uint64(utils.Bank),
				AccountSubType: uint64(utils.Savings),
				AccountName:    "AccountName",
				AccountNumber:  "AccountNumber",
				BankId:         bank.ID,
				BranchName:     "BranchName",
				RoutingNumber:  "RoutingNumber",
				IsDefault:      true,
			}
			res, err := new(services.PaymentAccountDetailService).Edit(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("PaymentAccountDetail Edited Successfully"))

			database.DBAPM(ctx).Model(&models.PaymentAccountDetail{}).First(&paymentAccount, paymentAccount.ID)
			Expect(paymentAccount.AccountType).To(Equal(utils.Bank))
			Expect(paymentAccount.AccountSubType).To(Equal(utils.Savings))
			Expect(paymentAccount.AccountName).To(Equal(param.AccountName))
			Expect(paymentAccount.AccountNumber).To(Equal(param.AccountNumber))
			Expect(paymentAccount.BankID).To(Equal(param.BankId))
			Expect(paymentAccount.BranchName).To(Equal(param.BranchName))
			Expect(paymentAccount.RoutingNumber).To(Equal(param.RoutingNumber))
			Expect(paymentAccount.IsDefault).To(Equal(true))

			database.DBAPM(ctx).Model(&models.Supplier{}).First(&supplier, supplier.ID)
			Expect(supplier.Status).To(Equal(models.SupplierStatusPending))
		})

		It("Should return error", func() {
			test_utils.SetPermission(&ctx, []string{})

			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{Status: models.SupplierStatusVerified})
			paymentAccount := test_helper.CreatePaymentAccountDetail(ctx, &models.PaymentAccountDetail{SupplierID: supplier.ID, AccountType: utils.Bank, IsDefault: true})
			bank := test_helper.CreateBank(ctx, &models.Bank{})
			database.DBAPM(ctx).Model(&supplier).Updates(&models.Supplier{Status: models.SupplierStatusVerified})
			param := &paymentpb.PaymentAccountDetailObject{
				Id:             paymentAccount.ID,
				AccountType:    uint64(utils.Bank),
				AccountSubType: uint64(utils.Savings),
				AccountName:    "AccountName",
				AccountNumber:  "AccountNumber",
				BankId:         bank.ID,
				BranchName:     "BranchName",
				RoutingNumber:  "RoutingNumber",
				IsDefault:      true,
			}
			res, err := new(services.PaymentAccountDetailService).Edit(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Change Not Allowed"))
		})
	})

	Context("Editing with existing account number with AppPreference", func() {
		It("Should return error response", func() {
			supplier1 := test_helper.CreateSupplier(ctx, &models.Supplier{})
			supplier2 := test_helper.CreateSupplier(ctx, &models.Supplier{})
			_ = test_helper.CreatePaymentAccountDetail(ctx, &models.PaymentAccountDetail{SupplierID: supplier1.ID, AccountType: utils.Bank, AccountNumber: "AccountNum", IsDefault: true})
			paymentAccount2 := test_helper.CreatePaymentAccountDetail(ctx, &models.PaymentAccountDetail{SupplierID: supplier2.ID, AccountType: utils.Bank, IsDefault: true})
			aaaModels.InjectMockAppPreferenceServiceInstance(mocks.GetAppPreferenceMock(map[string]interface{}{
				"enabled_account_number_validation": true,
			}))
			param := &paymentpb.PaymentAccountDetailObject{
				Id:            paymentAccount2.ID,
				AccountNumber: "AccountNum",
			}
			res, err := new(services.PaymentAccountDetailService).Edit(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Error while updating PaymentAccountDetail: Provided bank account number already exists"))
		})
	})

	Context("Editing with existing account number without AppPreference", func() {
		It("Should return error response", func() {
			supplier1 := test_helper.CreateSupplier(ctx, &models.Supplier{})
			supplier2 := test_helper.CreateSupplier(ctx, &models.Supplier{})
			_ = test_helper.CreatePaymentAccountDetail(ctx, &models.PaymentAccountDetail{SupplierID: supplier1.ID, AccountType: utils.Bank, AccountNumber: "AccountNum", IsDefault: true})
			paymentAccount2 := test_helper.CreatePaymentAccountDetail(ctx, &models.PaymentAccountDetail{SupplierID: supplier2.ID, AccountType: utils.Bank, IsDefault: true})
			aaaModels.InjectMockAppPreferenceServiceInstance(mocks.GetAppPreferenceMock(map[string]interface{}{
				"enabled_account_number_validation": false,
			}))
			param := &paymentpb.PaymentAccountDetailObject{
				Id:            paymentAccount2.ID,
				AccountNumber: "AccountNum",
			}
			res, err := new(services.PaymentAccountDetailService).Edit(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("PaymentAccountDetail Edited Successfully"))
		})
	})
})
