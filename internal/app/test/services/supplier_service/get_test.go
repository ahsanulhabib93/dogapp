package supplier_service_test

import (
	"context"

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
				SupplierCategoryMappings: []models.SupplierCategoryMapping{
					{CategoryID: 1},
					{CategoryID: 2},
				},
				SupplierType: utils.Hlc,
				SupplierSaMappings: []models.SupplierSaMapping{
					{SourcingAssociateId: 3},
					{SourcingAssociateId: 4},
				},
			})
			supplierAddress := test_helper.CreateSupplierAddress(ctx, &models.SupplierAddress{SupplierID: supplier.ID})
			paymentDetails := test_helper.CreatePaymentAccountDetail(ctx, &models.PaymentAccountDetail{SupplierID: supplier.ID, IsDefault: true})

			resp, err := new(services.SupplierService).Get(ctx, &supplierpb.GetSupplierParam{Id: supplier.ID})

			Expect(err).To(BeNil())
			Expect(resp.Success).To(Equal(true))
			Expect(resp.Data.Email).To(Equal(supplier.Email))
			Expect(resp.Data.Name).To(Equal(supplier.Name))
			Expect(resp.Data.CategoryIds).To(Equal([]uint64{1, 2}))
			Expect(resp.Data.SaIds).To(Equal([]uint64{3, 4}))
			Expect(len(resp.Data.SupplierAddresses)).To(Equal(1))
			Expect(len(resp.Data.PaymentAccountDetails)).To(Equal(1))
			Expect(resp.Data.SupplierAddresses[0].Id).To(Equal(supplierAddress.ID))
			Expect(resp.Data.PaymentAccountDetails[0].Id).To(Equal(paymentDetails.ID))
			Expect(resp.Data.SupplierType).To(Equal(uint64(utils.Hlc)))
			Expect(resp.Data.Status).To(Equal(models.SupplierStatusPending))
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
