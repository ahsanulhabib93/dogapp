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

var _ = Describe("AddSupplierAddress", func() {
	var ctx context.Context

	BeforeEach(func() {
		test_utils.GetContext(&ctx)
	})

	Context("While adding address for existing Supplier", func() {
		It("Should create address and return success response", func() {
			supplier := test_helper.CreateSupplierWithAddress(ctx, &models.Supplier{})
			param := &supplierpb.SupplierAddressParam{
				SupplierId: supplier.ID,
				Firstname:  "Firstname",
				Lastname:   "Lastname",
				Address1:   "Address1",
				Address2:   "Address2",
				Landmark:   "Landmark",
				City:       "City",
				State:      "State",
				Country:    "Country",
				Zipcode:    "Zipcode",
				Phone:      "Phone",
				GstNumber:  "GstNumber",
				IsDefault:  false,
			}
			res, err := new(services.SupplierService).AddSupplierAddress(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("SupplierAddress Added Successfully"))

			addresses := []*models.SupplierAddress{{}}
			database.DBAPM(ctx).Model(supplier).Association("SupplierAddresses").Find(&addresses)
			Expect(len(addresses)).To(Equal(2))
			address := addresses[1]

			Expect(address.Firstname).To(Equal(param.Firstname))
			Expect(address.Lastname).To(Equal(param.Lastname))
			Expect(address.Address1).To(Equal(param.Address1))
			Expect(address.Address2).To(Equal(param.Address2))
			Expect(address.Landmark).To(Equal(param.Landmark))
			Expect(address.City).To(Equal(param.City))
			Expect(address.State).To(Equal(param.State))
			Expect(address.Country).To(Equal(param.Country))
			Expect(address.Zipcode).To(Equal(param.Zipcode))
			Expect(address.Phone).To(Equal(param.Phone))
			Expect(address.GstNumber).To(Equal(param.GstNumber))
			Expect(address.IsDefault).To(Equal(false))
		})
	})

	Context("While adding address for invalid Supplier ID", func() {
		It("Should create address and return success response", func() {
			param := &supplierpb.SupplierAddressParam{
				SupplierId: 1000,
				Firstname:  "Firstname",
				Lastname:   "Lastname",
			}
			res, err := new(services.SupplierService).AddSupplierAddress(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Supplier Not Found"))
		})
	})
})
