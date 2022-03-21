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

var _ = Describe("EditSupplierAddress", func() {
	var ctx context.Context

	BeforeEach(func() {
		test_utils.GetContext(&ctx)
	})

	Context("Editing all attributes of existing Supplier address", func() {
		It("Should update address and return success response", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			defaultAddress := test_helper.CreateSupplierAddress(ctx, &models.SupplierAddress{SupplierID: supplier.ID, IsDefault: true})
			address := test_helper.CreateSupplierAddress(ctx, &models.SupplierAddress{SupplierID: supplier.ID, IsDefault: false})
			param := &addresspb.SupplierAddressObject{
				Id:        address.ID,
				Firstname: "Firstname",
				Lastname:  "Lastname",
				Address1:  "Address1",
				Address2:  "Address2",
				Landmark:  "Landmark",
				City:      "City",
				State:     "State",
				Country:   "Country",
				Zipcode:   "Zipcode",
				Phone:     "01123456789",
				GstNumber: "GstNumber",
				IsDefault: false,
			}
			res, err := new(services.SupplierAddressService).Edit(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Supplier Address Edited Successfully"))

			database.DBAPM(ctx).Model(&models.SupplierAddress{}).First(&address, address.ID)
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

			database.DBAPM(ctx).Model(&models.SupplierAddress{}).First(&defaultAddress, defaultAddress.ID)
			Expect(defaultAddress.IsDefault).To(Equal(true))
		})
	})

	Context("Updating default address as non-default", func() {
		It("Should return error response", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			address := test_helper.CreateSupplierAddress(ctx, &models.SupplierAddress{SupplierID: supplier.ID, IsDefault: true})
			param := &addresspb.SupplierAddressObject{
				Id:        address.ID,
				IsDefault: false,
			}
			res, err := new(services.SupplierAddressService).Edit(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Default address is required"))
		})
	})

	Context("Updating non-default address as default", func() {
		It("Should update address and return success response", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			address1 := test_helper.CreateSupplierAddress(ctx, &models.SupplierAddress{SupplierID: supplier.ID, IsDefault: true})
			address2 := test_helper.CreateSupplierAddress(ctx, &models.SupplierAddress{SupplierID: supplier.ID, IsDefault: false})
			param := &addresspb.SupplierAddressObject{
				Id:        address2.ID,
				IsDefault: true,
			}
			res, err := new(services.SupplierAddressService).Edit(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Supplier Address Edited Successfully"))

			database.DBAPM(ctx).Model(&models.SupplierAddress{}).First(&address2, address2.ID)
			Expect(address2.IsDefault).To(Equal(true))

			database.DBAPM(ctx).Model(&models.SupplierAddress{}).First(&address1, address1.ID)
			Expect(address1.IsDefault).To(Equal(false))
		})
	})

	Context("Editing only name of existing record", func() {
		It("Should update address and return success response", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			address := test_helper.CreateSupplierAddress(ctx, &models.SupplierAddress{SupplierID: supplier.ID, IsDefault: true})
			param := &addresspb.SupplierAddressObject{
				Id:        address.ID,
				Firstname: "Firstname",
				Lastname:  "Lastname",
				IsDefault: true,
			}
			res, err := new(services.SupplierAddressService).Edit(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Supplier Address Edited Successfully"))

			updatedAddress := &models.SupplierAddress{}
			database.DBAPM(ctx).Model(&models.SupplierAddress{}).First(&updatedAddress, address.ID)
			Expect(updatedAddress.Firstname).To(Equal(param.Firstname))
			Expect(updatedAddress.Lastname).To(Equal(param.Lastname))
			Expect(updatedAddress.Address1).To(Equal(address.Address1))
			Expect(updatedAddress.Address2).To(Equal(address.Address2))
			Expect(updatedAddress.Landmark).To(Equal(address.Landmark))
			Expect(updatedAddress.City).To(Equal(address.City))
			Expect(updatedAddress.State).To(Equal(address.State))
			Expect(updatedAddress.Country).To(Equal(address.Country))
			Expect(updatedAddress.Zipcode).To(Equal(address.Zipcode))
			Expect(updatedAddress.Phone).To(Equal(address.Phone))
			Expect(updatedAddress.GstNumber).To(Equal(address.GstNumber))
			Expect(updatedAddress.IsDefault).To(Equal(address.IsDefault))
		})
	})

	Context("Editing invalid supplier address", func() {
		It("Should return error response", func() {
			param := &addresspb.SupplierAddressObject{Id: 1000}
			res, err := new(services.SupplierAddressService).Edit(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Supplier Address Not Found"))
		})
	})

	Context("Editing supplier address with invalid number", func() {
		It("Should return error response", func() {
			suplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			supplierAddress := test_helper.CreateSupplierAddress(ctx, &models.SupplierAddress{SupplierID: suplier.ID})
			param := &addresspb.SupplierAddressObject{
				Id:    supplierAddress.ID,
				Phone: "123789",
			}

			res, err := new(services.SupplierAddressService).Edit(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Error while updating Supplier Address: Invalid Phone Number"))
		})
	})
})
