package supplier_service_test

import (
	"context"
	//"log"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	supplierpb "github.com/voonik/goConnect/api/go/ss2/supplier"
	test_utils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/services"
	"github.com/voonik/ss2/internal/app/test/test_helper"
	"github.com/voonik/ss2/internal/app/utils"
)

var _ = Describe("ListSupplier", func() {
	var ctx context.Context

	BeforeEach(func() {
		test_utils.GetContext(&ctx)
	})

	Context("Supplier List", func() {
		It("Should Respond with all the suppliers", func() {
			//log.Println("Supplier List")
			supplier1 := test_helper.CreateSupplier(ctx, &models.Supplier{
				SupplierCategoryMappings: []models.SupplierCategoryMapping{
					{CategoryID: 1},
					{CategoryID: 2},
				},
				SupplierType: utils.Hlc,
				SupplierSaMappings: []models.SupplierSaMapping{
					{SourcingAssociateId: 2},
					{SourcingAssociateId: 4},
				},
			})
			//log.Println("Supplier1 List", supplier1)

			supplier2 := test_helper.CreateSupplier(ctx, &models.Supplier{SupplierType: utils.L1})

			res, err := new(services.SupplierService).List(ctx, &supplierpb.ListParams{})
			Expect(err).To(BeNil())
			//log.Println("List Response", res.Data)
			Expect(len(res.Data)).To(Equal(2))
			supplierData1 := res.Data[0]
			Expect(supplierData1.Email).To(Equal(supplier1.Email))
			Expect(supplierData1.Name).To(Equal(supplier1.Name))
			Expect(supplierData1.CategoryIds).To(Equal([]uint64{1, 2}))
			Expect(supplierData1.SaIds).To(Equal([]uint64{2, 4}))
			Expect(supplierData1.SupplierType).To(Equal(uint64(utils.Hlc)))
			supplierData2 := res.Data[1]
			Expect(supplierData2.Email).To(Equal(supplier2.Email))
			Expect(supplierData2.Name).To(Equal(supplier2.Name))
			Expect(supplierData2.CategoryIds).To(Equal([]uint64{}))
			Expect(supplierData2.SaIds).To(Equal([]uint64{}))
			Expect(supplierData2.SupplierType).To(Equal(uint64(utils.L1)))
		})
	})
})
