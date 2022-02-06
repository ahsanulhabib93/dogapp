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

var _ = Describe("AddSupplier", func() {
	var ctx context.Context

	BeforeEach(func() {
		test_utils.GetContext(&ctx)
	})

	Context("Adding new Supplier", func() {
		It("Should create supplier and return success response", func() {
			param := &supplierpb.SupplierParam{
				Name:         "Name",
				Email:        "Email",
				SupplierType: uint64(utils.Hlc),
				Firstname:    "Firstname",
				Lastname:     "Lastname",
				Address1:     "Address1",
				Address2:     "Address2",
				Landmark:     "Landmark",
				City:         "City",
				State:        "State",
				Country:      "Country",
				Zipcode:      "Zipcode",
				Phone:        "Phone",
				GstNumber:    "GstNumber",
				CategoryIds:  []uint64{1, 30},
			}
			res, err := new(services.SupplierService).Add(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Supplier Added Successfully"))

			supplier := &models.Supplier{}
			database.DBAPM(ctx).Model(&models.Supplier{}).Where("name = ?", param.Name).Preload("SupplierCategoryMappings").First(&supplier)
			Expect(supplier.Email).To(Equal(param.Email))
			Expect(supplier.SupplierType).To(Equal(utils.Hlc))
			Expect(len(supplier.SupplierCategoryMappings)).To(Equal(2))
			Expect(supplier.SupplierCategoryMappings[1].CategoryID).To(Equal(uint64(30)))
			Expect(supplier.Status).To(Equal(models.SupplierStatusPending))

			addresses := []*models.SupplierAddress{{}}
			database.DBAPM(ctx).Model(supplier).Association("SupplierAddresses").Find(&addresses)
			Expect(len(addresses)).To(Equal(1))
			address := addresses[0]

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
			Expect(address.IsDefault).To(Equal(true))
		})
	})

	Context("Adding Supplier without name", func() {
		It("Should return error response", func() {
			param := &supplierpb.SupplierParam{
				Email:        "Email",
				SupplierType: uint64(utils.Hlc),
				Address1:     "Address1",
				Zipcode:      "Zipcode",
			}
			res, err := new(services.SupplierService).Add(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Error while creating Supplier: Name can't be blank"))
		})
	})

	Context("Adding Supplier with existing name", func() {
		It("Should return error response", func() {
			supplier1 := test_helper.CreateSupplier(ctx, &models.Supplier{SupplierType: utils.Hlc})
			param := &supplierpb.SupplierParam{
				Name:         supplier1.Name,
				Email:        "Email",
				SupplierType: uint64(utils.Hlc),
				Address1:     "Address1",
				Zipcode:      "Zipcode",
			}
			res, err := new(services.SupplierService).Add(ctx, param)
			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Error while creating Supplier: Name should be unique"))
		})
	})

	Context("Adding Supplier without supplier type", func() {
		It("Should return error response", func() {
			param := &supplierpb.SupplierParam{
				Name:     "Name",
				Email:    "Email",
				Address1: "Address1",
				Zipcode:  "Zipcode",
			}
			res, err := new(services.SupplierService).Add(ctx, param)
			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Error while creating Supplier: supplier_type can't be blank"))
		})
	})

})
