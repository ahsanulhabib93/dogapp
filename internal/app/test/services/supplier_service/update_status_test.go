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
)

var _ = Describe("UpdateStatus", func() {
	var ctx context.Context

	BeforeEach(func() {
		test_utils.GetContext(&ctx)
	})

	Context("Update Supplier status", func() {
		It("Should be updated and return success response", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			param := &supplierpb.UpdateStatusParam{
				Id:     supplier.ID,
				Status: string(models.SupplierStatusFailed),
				Reason: "test reason",
			}
			res, err := new(services.SupplierService).UpdateStatus(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Supplier status updated successfully"))

			updatedSupplier := models.Supplier{}
			database.DBAPM(ctx).Model(&models.Supplier{}).First(&updatedSupplier, supplier.ID)
			Expect(updatedSupplier.Status).To(Equal(models.SupplierStatusFailed))
			Expect(updatedSupplier.Reason).To(Equal(param.Reason))
		})
	})

	Context("Updating invalid supplier", func() {
		It("Should return error response", func() {
			param := &supplierpb.UpdateStatusParam{
				Id:     1000,
				Status: string(models.SupplierStatusVerified),
			}
			res, err := new(services.SupplierService).UpdateStatus(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Supplier Not Found"))
		})
	})
})
