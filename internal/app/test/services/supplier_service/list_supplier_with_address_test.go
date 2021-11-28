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

var _ = Describe("ListSupplierWithAddress", func() {
	var ctx context.Context

	BeforeEach(func() {
		test_utils.GetContext(&ctx)
	})

	Context("Without any filters", func() {
		It("Should return all the suppliers with addresses", func() {
			supplier1 := test_helper.CreateSupplier(ctx, &models.Supplier{SupplierType: utils.Hlc})
			supplierAddress1 := test_helper.CreateSupplierAddress(ctx, &models.SupplierAddress{SupplierID: supplier1.ID})
			supplier2 := test_helper.CreateSupplier(ctx, &models.Supplier{SupplierType: utils.L1})
			test_helper.CreateSupplierAddress(ctx, &models.SupplierAddress{SupplierID: supplier2.ID})
			test_helper.CreateSupplierAddress(ctx, &models.SupplierAddress{SupplierID: supplier2.ID})

			res, err := new(services.SupplierService).ListWithSupplierAddresses(ctx, &supplierpb.ListParams{})
			Expect(err).To(BeNil())
			Expect(len(res.Data)).To(Equal(2))

			supplierData1 := res.Data[0]
			Expect(supplierData1.Email).To(Equal(supplier1.Email))
			Expect(supplierData1.Name).To(Equal(supplier1.Name))
			Expect(supplierData1.SupplierType).To(Equal(uint64(utils.Hlc)))

			Expect(len(supplierData1.SupplierAddresses)).To(Equal(1))
			addressData := supplierData1.SupplierAddresses[0]
			Expect(addressData.Firstname).To(Equal(supplierAddress1.Firstname))
			Expect(addressData.Lastname).To(Equal(supplierAddress1.Lastname))
			Expect(addressData.Address1).To(Equal(supplierAddress1.Address1))
			Expect(addressData.Address2).To(Equal(supplierAddress1.Address2))
			Expect(addressData.Landmark).To(Equal(supplierAddress1.Landmark))
			Expect(addressData.City).To(Equal(supplierAddress1.City))
			Expect(addressData.State).To(Equal(supplierAddress1.State))
			Expect(addressData.Country).To(Equal(supplierAddress1.Country))
			Expect(addressData.Zipcode).To(Equal(supplierAddress1.Zipcode))
			Expect(addressData.Phone).To(Equal(supplierAddress1.Phone))
			Expect(addressData.GstNumber).To(Equal(supplierAddress1.GstNumber))
			Expect(addressData.IsDefault).To(Equal(false))

			supplierData2 := res.Data[1]
			Expect(supplierData2.Email).To(Equal(supplier2.Email))
			Expect(supplierData2.Name).To(Equal(supplier2.Name))
			Expect(supplierData2.SupplierType).To(Equal(uint64(utils.L1)))

			Expect(len(supplierData2.SupplierAddresses)).To(Equal(2))
			Expect(supplierData2.SupplierAddresses[0].IsDefault).To(Equal(false))
			Expect(supplierData2.SupplierAddresses[1].IsDefault).To(Equal(false))
		})
	})

	Context("With Supplier Id filter", func() {
		It("Should return corresponding supplier addresses", func() {
			supplier1 := test_helper.CreateSupplier(ctx, &models.Supplier{SupplierType: utils.Hlc})
			test_helper.CreateSupplierAddress(ctx, &models.SupplierAddress{SupplierID: supplier1.ID})
			supplier2 := test_helper.CreateSupplier(ctx, &models.Supplier{SupplierType: utils.L1})
			test_helper.CreateSupplierAddress(ctx, &models.SupplierAddress{SupplierID: supplier2.ID})
			test_helper.CreateSupplierAddress(ctx, &models.SupplierAddress{SupplierID: supplier2.ID})

			res, err := new(services.SupplierService).ListWithSupplierAddresses(ctx, &supplierpb.ListParams{Id: supplier2.ID})
			Expect(err).To(BeNil())
			Expect(len(res.Data)).To(Equal(1))

			supplierData1 := res.Data[0]
			Expect(supplierData1.Email).To(Equal(supplier2.Email))
			Expect(supplierData1.Name).To(Equal(supplier2.Name))
			Expect(supplierData1.SupplierType).To(Equal(uint64(utils.L1)))
			Expect(len(supplierData1.SupplierAddresses)).To(Equal(2))
		})
	})

	Context("With Supplier name filter", func() {
		It("Should return corresponding supplier addresses", func() {
			supplier1 := test_helper.CreateSupplier(ctx, &models.Supplier{SupplierType: utils.Hlc})
			test_helper.CreateSupplierAddress(ctx, &models.SupplierAddress{SupplierID: supplier1.ID})
			supplier2 := test_helper.CreateSupplier(ctx, &models.Supplier{Name: "string 123", SupplierType: utils.L1})
			test_helper.CreateSupplierAddress(ctx, &models.SupplierAddress{SupplierID: supplier2.ID})
			test_helper.CreateSupplierAddress(ctx, &models.SupplierAddress{SupplierID: supplier2.ID})

			res, err := new(services.SupplierService).ListWithSupplierAddresses(ctx, &supplierpb.ListParams{Name: "str"})
			Expect(err).To(BeNil())
			Expect(len(res.Data)).To(Equal(1))

			supplierData1 := res.Data[0]
			Expect(supplierData1.Email).To(Equal(supplier2.Email))
			Expect(supplierData1.Name).To(Equal(supplier2.Name))
			Expect(supplierData1.SupplierType).To(Equal(uint64(utils.L1)))
			Expect(len(supplierData1.SupplierAddresses)).To(Equal(2))
		})
	})

	Context("With Supplier email filter", func() {
		It("Should return corresponding supplier addresses", func() {
			supplier1 := test_helper.CreateSupplier(ctx, &models.Supplier{SupplierType: utils.Hlc})
			test_helper.CreateSupplierAddress(ctx, &models.SupplierAddress{SupplierID: supplier1.ID})
			supplier2 := test_helper.CreateSupplier(ctx, &models.Supplier{SupplierType: utils.L1})
			test_helper.CreateSupplierAddress(ctx, &models.SupplierAddress{SupplierID: supplier2.ID})
			test_helper.CreateSupplierAddress(ctx, &models.SupplierAddress{SupplierID: supplier2.ID})

			res, err := new(services.SupplierService).ListWithSupplierAddresses(ctx, &supplierpb.ListParams{Email: supplier2.Email})
			Expect(err).To(BeNil())
			Expect(len(res.Data)).To(Equal(1))

			supplierData1 := res.Data[0]
			Expect(supplierData1.Email).To(Equal(supplier2.Email))
			Expect(supplierData1.Name).To(Equal(supplier2.Name))
			Expect(supplierData1.SupplierType).To(Equal(uint64(utils.L1)))
			Expect(len(supplierData1.SupplierAddresses)).To(Equal(2))
		})
	})

	Context("With Phone filter", func() {
		It("Should return corresponding supplier addresses", func() {
			supplier1 := test_helper.CreateSupplier(ctx, &models.Supplier{SupplierType: utils.Hlc})
			test_helper.CreateSupplierAddress(ctx, &models.SupplierAddress{SupplierID: supplier1.ID})
			supplier2 := test_helper.CreateSupplier(ctx, &models.Supplier{SupplierType: utils.L1})
			address1 := test_helper.CreateSupplierAddress(ctx, &models.SupplierAddress{SupplierID: supplier2.ID})
			test_helper.CreateSupplierAddress(ctx, &models.SupplierAddress{SupplierID: supplier2.ID})

			res, err := new(services.SupplierService).ListWithSupplierAddresses(ctx, &supplierpb.ListParams{Phone: address1.Phone})
			Expect(err).To(BeNil())
			Expect(len(res.Data)).To(Equal(1))

			supplierData1 := res.Data[0]
			Expect(supplierData1.Email).To(Equal(supplier2.Email))
			Expect(supplierData1.Name).To(Equal(supplier2.Name))
			Expect(supplierData1.SupplierType).To(Equal(uint64(utils.L1)))
			Expect(len(supplierData1.SupplierAddresses)).To(Equal(2))
		})
	})

	Context("With City filter", func() {
		It("Should return corresponding supplier addresses", func() {
			supplier1 := test_helper.CreateSupplier(ctx, &models.Supplier{SupplierType: utils.Hlc})
			test_helper.CreateSupplierAddress(ctx, &models.SupplierAddress{SupplierID: supplier1.ID})
			supplier2 := test_helper.CreateSupplier(ctx, &models.Supplier{SupplierType: utils.L1})
			address1 := test_helper.CreateSupplierAddress(ctx, &models.SupplierAddress{SupplierID: supplier2.ID})
			test_helper.CreateSupplierAddress(ctx, &models.SupplierAddress{SupplierID: supplier2.ID})

			res, err := new(services.SupplierService).ListWithSupplierAddresses(ctx, &supplierpb.ListParams{City: address1.City})
			Expect(err).To(BeNil())
			Expect(len(res.Data)).To(Equal(1))

			supplierData1 := res.Data[0]
			Expect(supplierData1.Email).To(Equal(supplier2.Email))
			Expect(supplierData1.Name).To(Equal(supplier2.Name))
			Expect(supplierData1.SupplierType).To(Equal(uint64(utils.L1)))
			Expect(len(supplierData1.SupplierAddresses)).To(Equal(2))
		})
	})
})
