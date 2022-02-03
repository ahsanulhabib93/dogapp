package supplier_service_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	supplierpb "github.com/voonik/goConnect/api/go/ss2/supplier"
	"github.com/voonik/goFramework/pkg/database"
	test_utils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/services"
	"github.com/voonik/ss2/internal/app/test/test_helper"
	"github.com/voonik/ss2/internal/app/utils"
)

var _ = Describe("EditSupplier", func() {
	var ctx context.Context

	BeforeEach(func() {
		test_utils.GetContext(&ctx)
	})

	Context("Editing existing Supplier", func() {
		It("Should update supplier and return success response", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			param := &supplierpb.SupplierObject{
				Id:           supplier.ID,
				Name:         "Name",
				Email:        "Email",
				SupplierType: uint64(utils.L1),
			}
			res, err := new(services.SupplierService).Edit(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Supplier Edited Successfully"))

			database.DBAPM(ctx).Model(&models.Supplier{}).First(&supplier, supplier.ID)
			Expect(supplier.Email).To(Equal(param.Email))
			Expect(supplier.Name).To(Equal(param.Name))
			Expect(supplier.SupplierType).To(Equal(utils.L1))
		})
	})

	Context("Editing only one field of existing Supplier", func() {
		It("Should update supplier name and return success response", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			param := &supplierpb.SupplierObject{
				Id:   supplier.ID,
				Name: "Name",
			}
			res, err := new(services.SupplierService).Edit(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Supplier Edited Successfully"))

			updatedSupplier := &models.Supplier{}
			database.DBAPM(ctx).Model(&models.Supplier{}).First(&updatedSupplier, supplier.ID)
			Expect(updatedSupplier.Email).To(Equal(supplier.Email))
			Expect(updatedSupplier.SupplierType).To(Equal(utils.Hlc))
			Expect(updatedSupplier.Name).To(Equal(param.Name))
			Expect(updatedSupplier.Status).To(Equal(models.SupplierStatusPending))
		})

		It("Should update supplier status and return success response", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{
				Name: "test-supplier",
			})
			param := &supplierpb.SupplierObject{
				Id:     supplier.ID,
				Status: models.SupplierStatusActive,
			}
			res, err := new(services.SupplierService).Edit(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Supplier Edited Successfully"))

			updatedSupplier := &models.Supplier{}
			database.DBAPM(ctx).Model(&models.Supplier{}).First(&updatedSupplier, supplier.ID)
			Expect(updatedSupplier.Email).To(Equal(supplier.Email))
			Expect(updatedSupplier.SupplierType).To(Equal(utils.Hlc))
			Expect(updatedSupplier.Name).To(Equal(supplier.Name))
			Expect(updatedSupplier.Status).To(Equal(models.SupplierStatusActive))
		})

		It("Should return error on invalid status", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			param := &supplierpb.SupplierObject{
				Id:     supplier.ID,
				Status: "no idea",
			}
			res, err := new(services.SupplierService).Edit(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Error while updating Supplier: Status should be Active/Pending"))
		})
	})

	Context("Editing invalid supplier", func() {
		It("Should return error response", func() {
			param := &supplierpb.SupplierObject{Id: 1000}
			res, err := new(services.SupplierService).Edit(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Supplier Not Found"))
		})
	})

	Context("Editing with other supplier name", func() {
		It("Should return error response", func() {
			supplier1 := test_helper.CreateSupplier(ctx, &models.Supplier{})
			supplier2 := test_helper.CreateSupplier(ctx, &models.Supplier{})
			param := &supplierpb.SupplierObject{
				Id:   supplier1.ID,
				Name: supplier2.Name,
			}
			res, err := new(services.SupplierService).Edit(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Error while updating Supplier: Name should be unique"))
		})
	})
})
