package payment_account_detail_service_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	paymentpb "github.com/voonik/goConnect/api/go/ss2/payment_account_detail"
	"github.com/voonik/goFramework/pkg/database"
	test_utils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/services"
	"github.com/voonik/ss2/internal/app/test/test_helper"
	"github.com/voonik/ss2/internal/app/utils"
)

var _ = Describe("AddPaymentAccountDetail", func() {
	var ctx context.Context

	BeforeEach(func() {
		test_utils.GetContext(&ctx)
	})

	Context("Add", func() {
		It("Should create payment account detail and return success", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{SupplierType: utils.Hlc})
			bank := test_helper.CreateBank(ctx, &models.Bank{})
			param := paymentpb.PaymentAccountDetailParam{
				SupplierId:     supplier.ID,
				AccountType:    uint64(utils.Bank),
				AccountSubType: uint64(utils.Savings),
				AccountName:    "AccountName",
				AccountNumber:  "AccountNumber",
				BankId:         bank.ID,
				BranchName:     "BranchName",
				RoutingNumber:  "RoutingNumber",
				IsDefault:      true,
			}
			res, err := new(services.PaymentAccountDetailService).Add(ctx, &param)
			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("PaymentAccountDetail Added Successfully"))

			paymentAccounts := []*models.PaymentAccountDetail{{}}
			database.DBAPM(ctx).Model(supplier).Association("PaymentAccountDetails").Find(&paymentAccounts)
			Expect(len(paymentAccounts)).To(Equal(1))
			paymentAccount := paymentAccounts[0]

			Expect(paymentAccount.AccountType).To(Equal(utils.Bank))
			Expect(paymentAccount.AccountSubType).To(Equal(utils.Savings))
			Expect(paymentAccount.AccountName).To(Equal(param.AccountName))
			Expect(paymentAccount.AccountNumber).To(Equal(param.AccountNumber))
			Expect(paymentAccount.BankID).To(Equal(param.BankId))
			Expect(paymentAccount.BranchName).To(Equal(param.BranchName))
			Expect(paymentAccount.RoutingNumber).To(Equal(param.RoutingNumber))
			Expect(paymentAccount.IsDefault).To(Equal(true))
		})
	})

	Context("While adding default payment account", func() {
		It("Should return success response and other default payment should be updated as non-default", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{SupplierType: utils.Hlc})
			test_helper.CreatePaymentAccountDetail(ctx, &models.PaymentAccountDetail{SupplierID: supplier.ID, IsDefault: true})
			test_helper.CreatePaymentAccountDetail(ctx, &models.PaymentAccountDetail{SupplierID: supplier.ID, IsDefault: false})
			param := paymentpb.PaymentAccountDetailParam{
				SupplierId:     supplier.ID,
				AccountName:    "AccountName",
				AccountNumber:  "AccountNumber",
				AccountType:    uint64(utils.Bank),
				AccountSubType: uint64(utils.Savings),
				IsDefault:      true,
			}
			res, err := new(services.PaymentAccountDetailService).Add(ctx, &param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("PaymentAccountDetail Added Successfully"))

			paymentAccounts := []*models.PaymentAccountDetail{{}}
			database.DBAPM(ctx).Model(supplier).Association("PaymentAccountDetails").Find(&paymentAccounts)
			Expect(len(paymentAccounts)).To(Equal(3))

			paymentAccount1 := paymentAccounts[0]
			Expect(paymentAccount1.IsDefault).To(Equal(false))
			paymentAccount2 := paymentAccounts[1]
			Expect(paymentAccount2.IsDefault).To(Equal(false))
			paymentAccount3 := paymentAccounts[2]
			Expect(paymentAccount3.IsDefault).To(Equal(true))
		})
	})

	Context("While adding payment account detail for invalid Supplier ID", func() {
		It("Should return error response", func() {
			param := paymentpb.PaymentAccountDetailParam{
				SupplierId:  1000,
				AccountName: "AccountName",
			}
			res, err := new(services.PaymentAccountDetailService).Add(ctx, &param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Supplier Not Found"))
		})
	})

	Context("While adding payment account detail without account name", func() {
		It("Should return error response", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{SupplierType: utils.Hlc})
			param := paymentpb.PaymentAccountDetailParam{
				SupplierId:     supplier.ID,
				AccountNumber:  "AccountNumber",
				AccountType:    uint64(utils.Bank),
				AccountSubType: uint64(utils.Current),
				IsDefault:      true,
			}
			res, err := new(services.PaymentAccountDetailService).Add(ctx, &param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Error while creating PaymentAccountDetail: account_name can't be blank"))
		})
	})

	Context("While adding payment account detail without account number and account type", func() {
		It("Should return error response", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{SupplierType: utils.Hlc})
			param := paymentpb.PaymentAccountDetailParam{
				SupplierId:  supplier.ID,
				AccountName: "AccountName",
				IsDefault:   true,
			}
			res, err := new(services.PaymentAccountDetailService).Add(ctx, &param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Error while creating PaymentAccountDetail: Invalid Account SubType; account_type can't be blank; account_sub_type can't be blank; account_number can't be blank"))
		})
	})

	Context("While adding non-default payment account detail first time", func() {
		It("Should return error response", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{SupplierType: utils.Hlc})
			param := paymentpb.PaymentAccountDetailParam{
				SupplierId:     supplier.ID,
				AccountName:    "AccountName",
				AccountNumber:  "AccountNumber",
				AccountType:    uint64(utils.Bank),
				AccountSubType: uint64(utils.Current),
				IsDefault:      false,
			}
			res, err := new(services.PaymentAccountDetailService).Add(ctx, &param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Error while creating PaymentAccountDetail: Default Payment Account is required"))
		})
	})

	Context("While adding with invalid account subtype", func() {
		It("Should return error response", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{SupplierType: utils.Hlc})
			param := paymentpb.PaymentAccountDetailParam{
				SupplierId:     supplier.ID,
				AccountType:    uint64(utils.Bank),
				AccountSubType: uint64(utils.Bkash),
				AccountName:    "AccountName",
				AccountNumber:  "AccountNumber",
				IsDefault:      true,
			}
			res, err := new(services.PaymentAccountDetailService).Add(ctx, &param)
			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Error while creating PaymentAccountDetail: Invalid Account SubType"))
		})
	})

	Context("While adding with invalid bank ID", func() {
		It("Should return error response", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{SupplierType: utils.Hlc})
			param := paymentpb.PaymentAccountDetailParam{
				SupplierId:     supplier.ID,
				AccountType:    uint64(utils.Bank),
				AccountSubType: uint64(utils.Savings),
				AccountName:    "AccountName",
				AccountNumber:  "AccountNumber",
				BankId:         1000,
				IsDefault:      true,
			}
			res, err := new(services.PaymentAccountDetailService).Add(ctx, &param)
			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Error while creating PaymentAccountDetail: Invalid Bank Name"))
		})
	})
})
