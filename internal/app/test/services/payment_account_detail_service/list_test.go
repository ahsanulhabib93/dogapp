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

var _ = Describe("ListPaymentAccountDetail", func() {
	var ctx context.Context

	BeforeEach(func() {
		test_utils.GetContext(&ctx)
	})

	Context("List", func() {
		It("Should Respond with all the Payment Account Details", func() {
			supplier1 := test_helper.CreateSupplier(ctx, &models.Supplier{})
			bank := test_helper.CreateBank(ctx, &models.Bank{})
			accountDetail1 := test_helper.CreatePaymentAccountDetail(ctx, &models.PaymentAccountDetail{SupplierID: supplier1.ID, AccountType: utils.Mfs, BankID: bank.ID, IsDefault: true})
			accountDetail2 := test_helper.CreatePaymentAccountDetail(ctx, &models.PaymentAccountDetail{SupplierID: supplier1.ID, AccountType: utils.Bank, BankID: bank.ID})
			accountDetail3 := test_helper.CreatePaymentAccountDetail(ctx, &models.PaymentAccountDetail{SupplierID: supplier1.ID, AccountType: utils.PrepaidCard, AccountSubType: utils.EBL, BankID: bank.ID})
			extraDetails := models.PaymentAccountDetailExtraDetails{
				EmployeeId: uint64(12344),
				ClientId:   uint64(123),
				ExpiryDate: "2025-01-02",
				Token:      "sample_token_1",
				UniqueId:   "SS2-PAD-3",
			}
			accountDetail3.SetExtraDetails(extraDetails)
			database.DBAPM(ctx).Save(accountDetail3)

			res, err := new(services.PaymentAccountDetailService).List(ctx, &paymentpb.ListParams{SupplierId: supplier1.ID})
			Expect(err).To(BeNil())
			Expect(len(res.Data)).To(Equal(3))

			accountData1 := res.Data[0]
			Expect(accountData1.AccountType).To(Equal(uint64(utils.Mfs)))
			Expect(accountData1.AccountSubType).To(Equal(uint64(utils.Bkash)))
			Expect(accountData1.AccountName).To(Equal(accountDetail1.AccountName))
			Expect(accountData1.AccountNumber).To(Equal(accountDetail1.AccountNumber))
			Expect(accountData1.BankId).To(Equal(bank.ID))
			Expect(accountData1.BankName).To(Equal(bank.Name))
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

			accountData3 := res.Data[2]
			Expect(accountData3.AccountType).To(Equal(uint64(utils.PrepaidCard)))
			Expect(accountData3.AccountSubType).To(Equal(uint64(utils.EBL)))
			Expect(accountData3.AccountName).To(Equal(accountDetail3.AccountName))
			Expect(accountData3.AccountNumber).To(Equal(accountDetail3.AccountNumber))
			Expect(accountData3.BankId).To(Equal(bank.ID))
			Expect(accountData3.BankName).To(Equal(bank.Name))
			Expect(accountData3.BranchName).To(Equal(accountDetail3.BranchName))
			Expect(accountData3.RoutingNumber).To(Equal(accountDetail3.RoutingNumber))
			Expect(accountData3.IsDefault).To(Equal(false))

			testExtraDetails := models.PaymentAccountDetailExtraDetails{}
			utils.CopyStructAtoB(accountData3.ExtraDetails, &testExtraDetails)
			Expect(testExtraDetails.ClientId).To(Equal(uint64(123)))
			Expect(testExtraDetails.EmployeeId).To(Equal(uint64(12344)))
			Expect(testExtraDetails.ExpiryDate).To(Equal("2025-01-02"))
			Expect(testExtraDetails.Token).To(Equal("sample_token_1"))
			Expect(testExtraDetails.UniqueId).To(Equal("SS2-PAD-3"))
		})
	})
})
