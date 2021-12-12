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

var _ = Describe("EditPaymentAccountDetail", func() {
	var ctx context.Context

	BeforeEach(func() {
		test_utils.GetContext(&ctx)
	})

	Context("Editing all attributes of existing PaymentAccount", func() {
		It("Should update and return success response", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			paymentAccount := test_helper.CreatePaymentAccountDetail(ctx, &models.PaymentAccountDetail{SupplierID: supplier.id, AccountType: utils.Bank})
			param := &paymentpb.PaymentAccountDetailObject{
				Id:            paymentAccount.ID,
				AccountType:   uint64(utils.Bank),
				AccountName:   "AccountName",
				AccountNumber: "AccountNumber",
				BankName:      "BankName",
				BranchName:    "BranchName",
				RoutingNumber: "RoutingNumber",
				IsDefault:     false,
			}
			res, err := new(services.PaymentAccountDetailService).Edit(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("PaymentAccountDetail Edited Successfully"))

			database.DBAPM(ctx).Model(&models.PaymentAccountDetail{}).First(&paymentAccount, address.ID)
			Expect(paymentAccount.AccountType).To(Equal(utils.Bank))
			Expect(paymentAccount.AccountName).To(Equal(param.AccountName))
			Expect(paymentAccount.AccountNumber).To(Equal(param.AccountNumber))
			Expect(paymentAccount.BankName).To(Equal(param.BankName))
			Expect(paymentAccount.BranchName).To(Equal(param.BranchName))
			Expect(paymentAccount.RoutingNumber).To(Equal(param.RoutingNumber))
			Expect(paymentAccount.IsDefault).To(Equal(false))
		})
	})

	Context("Editing only account number of existing record", func() {
		It("Should update and return success response", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			paymentAccount := test_helper.CreatePaymentAccountDetail(ctx, &models.PaymentAccountDetail{SupplierID: supplier.id, AccountType: utils.Bank})
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
			Expect(updatedPayment.AccountName).To(Equal(param.AccountName))
			Expect(updatedPayment.AccountNumber).To(Equal(param.AccountNumber))
			Expect(updatedPayment.BankName).To(Equal(paymentAccount.BankName))
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
})
