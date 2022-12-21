package supplier_address_service_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	addresspb "github.com/voonik/goConnect/api/go/ss2/supplier_address"
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

	Context("Supplier Address List When Filtered with Supplier id", func() {
		It("Should Respond with corresponding addresses", func() {
			test_helper.CreateSupplierWithAddress(ctx, &models.Supplier{})
			supplier2 := test_helper.CreateSupplier(ctx, &models.Supplier{SupplierType: utils.L1})
			supplierAddress1 := test_helper.CreateSupplierAddress(ctx, &models.SupplierAddress{SupplierID: supplier2.ID})
			supplierAddress2 := test_helper.CreateSupplierAddress(ctx, &models.SupplierAddress{SupplierID: supplier2.ID})

			res, err := new(services.SupplierAddressService).List(ctx, &addresspb.ListSupplierAddressParams{SupplierId: supplier2.ID})
			Expect(err).To(BeNil())
			Expect(len(res.Data)).To(Equal(2))

			addressData1 := res.Data[0]
			Expect(addressData1.Firstname).To(Equal(supplierAddress1.Firstname))
			Expect(addressData1.Lastname).To(Equal(supplierAddress1.Lastname))
			Expect(addressData1.Address1).To(Equal(supplierAddress1.Address1))
			Expect(addressData1.Address2).To(Equal(supplierAddress1.Address2))
			Expect(addressData1.Landmark).To(Equal(supplierAddress1.Landmark))
			Expect(addressData1.City).To(Equal(supplierAddress1.City))
			Expect(addressData1.State).To(Equal(supplierAddress1.State))
			Expect(addressData1.Country).To(Equal(supplierAddress1.Country))
			Expect(addressData1.Zipcode).To(Equal(supplierAddress1.Zipcode))
			Expect(addressData1.Phone).To(Equal(supplierAddress1.Phone))
			Expect(addressData1.GstNumber).To(Equal(supplierAddress1.GstNumber))
			Expect(addressData1.IsDefault).To(Equal(false))

			addressData2 := res.Data[1]
			Expect(addressData2.Firstname).To(Equal(supplierAddress2.Firstname))
			Expect(addressData2.Lastname).To(Equal(supplierAddress2.Lastname))
			Expect(addressData2.Address1).To(Equal(supplierAddress2.Address1))
			Expect(addressData2.Address2).To(Equal(supplierAddress2.Address2))
			Expect(addressData2.Landmark).To(Equal(supplierAddress2.Landmark))
			Expect(addressData2.City).To(Equal(supplierAddress2.City))
			Expect(addressData2.State).To(Equal(supplierAddress2.State))
			Expect(addressData2.Country).To(Equal(supplierAddress2.Country))
			Expect(addressData2.Zipcode).To(Equal(supplierAddress2.Zipcode))
			Expect(addressData2.Phone).To(Equal(supplierAddress2.Phone))
			Expect(addressData2.GstNumber).To(Equal(supplierAddress2.GstNumber))
			Expect(addressData2.IsDefault).To(Equal(false))
		})
	})

	Context("Supplier Address List When Filtered with  id", func() {
		It("Should Respond with corresponding addresses", func() {
			test_helper.CreateSupplierWithAddress(ctx, &models.Supplier{})
			supplier2 := test_helper.CreateSupplier(ctx, &models.Supplier{SupplierType: utils.L1})
			supplierAddress1 := test_helper.CreateSupplierAddress(ctx, &models.SupplierAddress{SupplierID: supplier2.ID})
			supplierAddress2 := test_helper.CreateSupplierAddress(ctx, &models.SupplierAddress{SupplierID: supplier2.ID})

			res, err := new(services.SupplierAddressService).List(ctx, &addresspb.ListSupplierAddressParams{Id: supplierAddress1.ID})
			Expect(err).To(BeNil())
			Expect(len(res.Data)).To(Equal(1))

			addressData1 := res.Data[0]
			Expect(addressData1.Firstname).To(Equal(supplierAddress1.Firstname))
			Expect(addressData1.Lastname).To(Equal(supplierAddress1.Lastname))
			Expect(addressData1.Address1).To(Equal(supplierAddress1.Address1))
			Expect(addressData1.Address2).To(Equal(supplierAddress1.Address2))
			Expect(addressData1.Landmark).To(Equal(supplierAddress1.Landmark))
			Expect(addressData1.City).To(Equal(supplierAddress1.City))
			Expect(addressData1.State).To(Equal(supplierAddress1.State))
			Expect(addressData1.Country).To(Equal(supplierAddress1.Country))
			Expect(addressData1.Zipcode).To(Equal(supplierAddress1.Zipcode))
			Expect(addressData1.Phone).To(Equal(supplierAddress1.Phone))
			Expect(addressData1.GstNumber).To(Equal(supplierAddress1.GstNumber))
			Expect(addressData1.IsDefault).To(Equal(false))
		})
	})
})
