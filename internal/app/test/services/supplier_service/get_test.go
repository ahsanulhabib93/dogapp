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

			resp, err := new(services.SupplierService).Get(ctx, &supplierpb.GetSupplierParam{Id: supplier.ID})
			Expect(err).To(BeNil())
			Expect(resp.Email).To(Equal(supplier.Email))
			Expect(resp.Name).To(Equal(supplier.Name))
			Expect(resp.CategoryIds).To(Equal([]uint64{1, 2}))
			Expect(resp.SaIds).To(Equal([]uint64{3, 4}))
			Expect(resp.SupplierType).To(Equal(uint64(utils.Hlc)))
			Expect(resp.Status).To(Equal(models.SupplierStatusPending))
		})
	})

	Context("Invalid supplier", func() {
		It("Should Respond with error", func() {
			_, err := new(services.SupplierService).Get(ctx, &supplierpb.GetSupplierParam{Id: 10})
			Expect(err).To(Equal(services.NotFound))
		})
	})
})
