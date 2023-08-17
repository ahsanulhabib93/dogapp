package payment_account_detail_service_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	paymentpb "github.com/voonik/goConnect/api/go/ss2/payment_account_detail"
	test_utils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/services"
	"github.com/voonik/ss2/internal/app/test/test_helper"
	"github.com/voonik/ss2/internal/app/utils"
)

var _ = Describe("ListPaymentAccountDetail", func() {
	var ctx context.Context

	BeforeEach(func() {
		test_utils.GetContext(&ctx)
	})

	Context("List", func() {
		It("Should Respond with all the Payment Account Details", func() {
			supplier1 := test_helper.CreateSupplier(ctx, &models.Supplier{})
			accountDetail1 := test_helper.CreatePaymentAccountDetail(ctx, &models.PaymentAccountDetail{SupplierID: supplier1.ID, AccountType: utils.Mfs, IsDefault: true})
			bank := test_helper.CreateBank(ctx, &models.Bank{})
			accountDetail2 := test_helper.CreatePaymentAccountDetail(ctx, &models.PaymentAccountDetail{SupplierID: supplier1.ID, AccountType: utils.Bank, BankID: bank.ID})

			res, err := new(services.PaymentAccountDetailService).List(ctx, &paymentpb.ListParams{SupplierId: supplier1.ID})
			Expect(err).To(BeNil())
			Expect(len(res.Data)).To(Equal(2))

			accountData1 := res.Data[0]
			Expect(accountData1.AccountType).To(Equal(uint64(utils.Mfs)))
			Expect(accountData1.AccountSubType).To(Equal(uint64(utils.Bkash)))
			Expect(accountData1.AccountName).To(Equal(accountDetail1.AccountName))
			Expect(accountData1.AccountNumber).To(Equal(accountDetail1.AccountNumber))
			Expect(accountData1.IsDefault).To(Equal(true))

			accountData2 := res.Data[1]
			Expect(accountData2.AccountType).To(Equal(uint64(utils.Bank)))
			Expect(accountData2.AccountSubType).To(Equal(uint64(utils.Current)))
			Expect(accountData2.AccountName).To(Equal(accountDetail2.AccountName))
			Expect(accountData2.AccountNumber).To(Equal(accountDetail2.AccountNumber))
			Expect(accountData2.BankId).To(Equal(bank.ID))
			Expect(accountData2.BankName).To(Equal(bank.Name))
			Expect(accountData2.BranchName).To(Equal(accountDetail2.BranchName))
			Expect(accountData2.RoutingNumber).To(Equal(accountDetail2.RoutingNumber))
			Expect(accountData2.IsDefault).To(Equal(false))
		})
	})
})
