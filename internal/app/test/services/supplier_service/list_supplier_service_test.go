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
)

var _ = Describe("ListSupplier", func() {
	var ctx context.Context

	BeforeEach(func() {
		test_utils.GetContext(&ctx)
	})

	Context("Supplier List", func() {
		It("It Should Respond with all the suppliers", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})

			res, err := new(services.SupplierService).List(ctx, &supplierpb.ListParams{})
			Expect(err).To(BeNil())
			Expect(len(res.Data)).To(Equal(1))

			data := res.Data[0]
			Expect(data.Email).To(Equal(supplier.Email))
			Expect(data.Name).To(Equal(supplier.Name))
			Expect(data.SupplierType).To(Equal(uint64(supplier.SupplierType)))
		})
	})
})
