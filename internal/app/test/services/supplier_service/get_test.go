package supplier_service_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	supplierpb "github.com/voonik/goConnect/api/go/ss2/supplier"
	test_utils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/services"
	"github.com/voonik/ss2/internal/app/test/test_helper"
	"github.com/voonik/ss2/internal/app/utils"
)

var _ = Describe("GetSupplier", func() {
	var ctx context.Context

	BeforeEach(func() {
		test_utils.GetContext(&ctx)
	})

	Context("Get valid supplier", func() {
		It("Should Respond with success", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{
				SupplierType:             utils.Hlc,
				IsPhoneVerified:          true,
				SupplierCategoryMappings: []models.SupplierCategoryMapping{{CategoryID: 1}, {CategoryID: 2}},
				SupplierOpcMappings:      []models.SupplierOpcMapping{{ProcessingCenterID: 3}, {ProcessingCenterID: 4}},
			})
			supplierAddress := test_helper.CreateSupplierAddress(ctx, &models.SupplierAddress{SupplierID: supplier.ID})
			paymentDetails := test_helper.CreatePaymentAccountDetail(ctx, &models.PaymentAccountDetail{SupplierID: supplier.ID, IsDefault: true})
			resp, err := new(services.SupplierService).Get(ctx, &supplierpb.GetSupplierParam{Id: supplier.ID})

			Expect(err).To(BeNil())
			Expect(resp.Success).To(Equal(true))
			Expect(resp.Data.Email).To(Equal(supplier.Email))
			Expect(resp.Data.Name).To(Equal(supplier.Name))
			Expect(resp.Data.Phone).To(Equal(supplier.Phone))
			Expect(resp.Data.AlternatePhone).To(Equal(supplier.AlternatePhone))
			Expect(resp.Data.BusinessName).To(Equal(supplier.BusinessName))
			Expect(resp.Data.ShopImageUrl).To(Equal(supplier.ShopImageURL))
			Expect(resp.Data.Reason).To(Equal(supplier.Reason))
			Expect(resp.Data.IsPhoneVerified).To(Equal(true))
			Expect(resp.Data.CategoryIds).To(Equal([]uint64{1, 2}))
			Expect(resp.Data.OpcIds).To(Equal([]uint64{3, 4}))
			Expect(resp.Data.SupplierType).To(Equal(uint64(utils.Hlc)))
			Expect(resp.Data.Status).To(Equal(string(models.SupplierStatusPending)))

			Expect(len(resp.Data.SupplierAddresses)).To(Equal(1))
			Expect(resp.Data.SupplierAddresses[0].Firstname).To(Equal(supplierAddress.Firstname))
			Expect(resp.Data.SupplierAddresses[0].Address1).To(Equal(supplierAddress.Address1))
			Expect(resp.Data.SupplierAddresses[0].Phone).To(Equal(supplierAddress.Phone))

			Expect(len(resp.Data.PaymentAccountDetails)).To(Equal(1))
			Expect(resp.Data.PaymentAccountDetails[0].Id).To(Equal(paymentDetails.ID))
			Expect(resp.Data.PaymentAccountDetails[0].AccountName).To(Equal(paymentDetails.AccountName))
			Expect(resp.Data.PaymentAccountDetails[0].AccountNumber).To(Equal(paymentDetails.AccountNumber))
			Expect(resp.Data.PaymentAccountDetails[0].BankId).To(Equal(paymentDetails.BankID))
			Expect(resp.Data.PaymentAccountDetails[0].BranchName).To(Equal(paymentDetails.BranchName))
			Expect(resp.Data.PaymentAccountDetails[0].BankName).To(Equal(resp.Data.PaymentAccountDetails[0].BankName))
			Expect(resp.Data.PaymentAccountDetails[0].RoutingNumber).To(Equal(paymentDetails.RoutingNumber))
			Expect(resp.Data.PaymentAccountDetails[0].IsDefault).To(Equal(paymentDetails.IsDefault))
			Expect(resp.Data.PaymentAccountDetails[0].AccountType).To(Equal(uint64(paymentDetails.AccountType)))
			Expect(resp.Data.PaymentAccountDetails[0].AccountSubType).To(Equal(uint64(paymentDetails.AccountSubType)))
		})

		It("Should Respond with only non-deleted mapping", func() {
			deletedAt := time.Now()
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{
				SupplierCategoryMappings: []models.SupplierCategoryMapping{
					{CategoryID: 1, DeletedAt: &deletedAt},
					{CategoryID: 2},
				},
				SupplierType: utils.Hlc,
				SupplierOpcMappings: []models.SupplierOpcMapping{
					{ProcessingCenterID: 2},
					{ProcessingCenterID: 3, DeletedAt: &deletedAt},
					{ProcessingCenterID: 4},
				},
			})

			resp, err := new(services.SupplierService).Get(ctx, &supplierpb.GetSupplierParam{Id: supplier.ID})

			Expect(err).To(BeNil())
			Expect(resp.Success).To(Equal(true))
			Expect(resp.Data.Email).To(Equal(supplier.Email))
			Expect(resp.Data.Name).To(Equal(supplier.Name))
			Expect(resp.Data.CategoryIds).To(Equal([]uint64{2}))
			Expect(resp.Data.OpcIds).To(Equal([]uint64{2, 4}))
			Expect(resp.Data.SupplierType).To(Equal(uint64(utils.Hlc)))
		})
	})

	Context("Invalid supplier", func() {
		It("Should Respond with error", func() {
			res, err := new(services.SupplierService).Get(ctx, &supplierpb.GetSupplierParam{Id: 10})
			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
		})
	})
})
