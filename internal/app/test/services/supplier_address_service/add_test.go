package supplier_address_service_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	addresspb "github.com/voonik/goConnect/api/go/ss2/supplier_address"
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
			param := &addresspb.SupplierAddressParam{
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
				Phone:      "01123456789",
				GstNumber:  "GstNumber",
				IsDefault:  false,
			}
			res, err := new(services.SupplierAddressService).Add(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Supplier Address Added Successfully"))

			addresses := []*models.SupplierAddress{{}}
			database.DBAPM(ctx).Model(supplier).Association("SupplierAddresses").Find(&addresses)
			Expect(len(addresses)).To(Equal(2))

			address1 := addresses[0]
			Expect(address1.IsDefault).To(Equal(true))

			address2 := addresses[1]
			Expect(address2.Firstname).To(Equal(param.Firstname))
			Expect(address2.Lastname).To(Equal(param.Lastname))
			Expect(address2.Address1).To(Equal(param.Address1))
			Expect(address2.Address2).To(Equal(param.Address2))
			Expect(address2.Landmark).To(Equal(param.Landmark))
			Expect(address2.City).To(Equal(param.City))
			Expect(address2.State).To(Equal(param.State))
			Expect(address2.Country).To(Equal(param.Country))
			Expect(address2.Zipcode).To(Equal(param.Zipcode))
			Expect(address2.Phone).To(Equal(param.Phone))
			Expect(address2.GstNumber).To(Equal(param.GstNumber))
			Expect(address2.IsDefault).To(Equal(false))
		})
	})

	Context("While adding address for existing Supplier without previous address", func() {
		It("Should create address and return success response", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			param := &addresspb.SupplierAddressParam{
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
				Phone:      "01123456789",
				GstNumber:  "GstNumber",
			}
			res, err := new(services.SupplierAddressService).Add(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Supplier Address Added Successfully"))

			addresses := []*models.SupplierAddress{{}}
			database.DBAPM(ctx).Model(supplier).Association("SupplierAddresses").Find(&addresses)
			Expect(len(addresses)).To(Equal(1))

			address1 := addresses[0]
			Expect(address1.Firstname).To(Equal(param.Firstname))
			Expect(address1.Lastname).To(Equal(param.Lastname))
			Expect(address1.Address1).To(Equal(param.Address1))
			Expect(address1.Address2).To(Equal(param.Address2))
			Expect(address1.Landmark).To(Equal(param.Landmark))
			Expect(address1.City).To(Equal(param.City))
			Expect(address1.State).To(Equal(param.State))
			Expect(address1.Country).To(Equal(param.Country))
			Expect(address1.Zipcode).To(Equal(param.Zipcode))
			Expect(address1.Phone).To(Equal(param.Phone))
			Expect(address1.GstNumber).To(Equal(param.GstNumber))
			Expect(address1.IsDefault).To(Equal(true))
		})
	})

	Context("While adding default address for existing Supplier", func() {
		It("Should create address and return success response", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			test_helper.CreateSupplierAddress(ctx, &models.SupplierAddress{SupplierID: supplier.ID, IsDefault: true})
			test_helper.CreateSupplierAddress(ctx, &models.SupplierAddress{SupplierID: supplier.ID, IsDefault: false})
			param := &addresspb.SupplierAddressParam{
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
				Phone:      "01123456789",
				GstNumber:  "GstNumber",
				IsDefault:  true,
			}
			res, err := new(services.SupplierAddressService).Add(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Supplier Address Added Successfully"))

			addresses := []*models.SupplierAddress{{}}
			database.DBAPM(ctx).Model(supplier).Association("SupplierAddresses").Find(&addresses)
			Expect(len(addresses)).To(Equal(3))

			address1 := addresses[0]
			Expect(address1.IsDefault).To(Equal(false))
			address2 := addresses[1]
			Expect(address2.IsDefault).To(Equal(false))

			address3 := addresses[2]
			Expect(address3.Firstname).To(Equal(param.Firstname))
			Expect(address3.Lastname).To(Equal(param.Lastname))
			Expect(address3.Address1).To(Equal(param.Address1))
			Expect(address3.Address2).To(Equal(param.Address2))
			Expect(address3.Landmark).To(Equal(param.Landmark))
			Expect(address3.City).To(Equal(param.City))
			Expect(address3.State).To(Equal(param.State))
			Expect(address3.Country).To(Equal(param.Country))
			Expect(address3.Zipcode).To(Equal(param.Zipcode))
			Expect(address3.Phone).To(Equal(param.Phone))
			Expect(address3.GstNumber).To(Equal(param.GstNumber))
			Expect(address3.IsDefault).To(Equal(true))
		})
	})

	Context("While adding address for invalid Supplier ID", func() {
		It("Should return error response", func() {
			param := &addresspb.SupplierAddressParam{
				SupplierId: 1000,
				Firstname:  "Firstname",
				Lastname:   "Lastname",
				Phone:      "01123456789",
			}
			res, err := new(services.SupplierAddressService).Add(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Supplier Not Found"))
		})
	})

	Context("While adding address without zipcode", func() {
		It("Should return success response", func() {
			supplier := test_helper.CreateSupplierWithAddress(ctx, &models.Supplier{})
			param := &addresspb.SupplierAddressParam{
				SupplierId: supplier.ID,
				Firstname:  "Firstname",
				Lastname:   "Lastname",
				Address1:   "Address1",
				Phone:      "01123456789",
			}
			res, err := new(services.SupplierAddressService).Add(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Supplier Address Added Successfully"))
		})
	})

	Context("While adding address without address1", func() {
		It("Should return error response", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			param := &addresspb.SupplierAddressParam{
				SupplierId: supplier.ID,
				Firstname:  "Firstname",
				Lastname:   "Lastname",
				Zipcode:    "Zipcode",
				Phone:      "01123456789",
			}
			res, err := new(services.SupplierAddressService).Add(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Error while creating Supplier Address: Address1 can't be blank"))
		})
	})

	Context("While adding address with invalid phone number", func() {
		It("Should return error response", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			param := &addresspb.SupplierAddressParam{
				SupplierId: supplier.ID,
				Firstname:  "Firstname",
				Lastname:   "Lastname",
				Zipcode:    "Zipcode",
				Address1:   "Address1",
				Phone:      "123456789",
			}

			res, err := new(services.SupplierAddressService).Add(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Error while creating Supplier Address: Invalid Phone Number"))
		})
	})

	Context("While adding address with empty phone number", func() {
		It("Should return error response", func() {
			supplier := test_helper.CreateSupplierWithAddress(ctx, &models.Supplier{})
			param := &addresspb.SupplierAddressParam{
				SupplierId: supplier.ID,
				Firstname:  "Firstname",
				Lastname:   "Lastname",
				Zipcode:    "Zipcode",
				Address1:   "Address1",
			}
			res, err := new(services.SupplierAddressService).Add(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Error while creating Supplier Address: Phone Number can't be blank"))
		})
	})
})
