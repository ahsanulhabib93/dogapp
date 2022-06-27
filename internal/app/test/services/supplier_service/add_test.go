package supplier_service_test

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	opcPb "github.com/voonik/goConnect/api/go/oms/processing_center"
	supplierpb "github.com/voonik/goConnect/api/go/ss2/supplier"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/goFramework/pkg/misc"
	test_utils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/services"
	"github.com/voonik/ss2/internal/app/test/mocks"
	"github.com/voonik/ss2/internal/app/test/test_helper"
	"github.com/voonik/ss2/internal/app/utils"
)

var _ = Describe("AddSupplier", func() {
	var ctx context.Context
	var userId uint64 = uint64(101)

	BeforeEach(func() {
		test_utils.GetContext(&ctx)
		mocks.UnsetOpcMock()

		threadObject := &misc.ThreadObject{
			VaccountId:    1,
			PortalId:      1,
			CurrentActId:  1,
			XForwardedFor: "5079327",
			UserData: &misc.UserData{
				UserId: userId,
				Name:   "John",
				Email:  "john@gmail.com",
				Phone:  "8801855533367",
			},
		}
		ctx = misc.SetInContextThreadObject(ctx, threadObject)
	})

	Context("Adding new Supplier", func() {
		It("Should create supplier and return success response", func() {
			opcIds := []uint64{5000, 6000}
			mockOpc := mocks.SetOpcMock()
			mockOpc.On("GetProcessingCenterListWithOpcIds", ctx, opcIds).Return(&opcPb.ProcessingCenterListResponse{Data: []*opcPb.OpcDetail{
				{OpcId: 5000},
				{OpcId: 6000},
			}}, nil)

			param := &supplierpb.SupplierParam{
				Name:                     "Name",
				Email:                    "Email",
				SupplierType:             uint64(utils.Hlc),
				BusinessName:             "BusinessName",
				Phone:                    "8801234567890",
				AlternatePhone:           "8801234567891",
				ShopImageUrl:             "ss2/shop_images/test.png",
				Firstname:                "Firstname",
				Lastname:                 "Lastname",
				Address1:                 "Address1",
				Address2:                 "Address2",
				Landmark:                 "Landmark",
				City:                     "City",
				State:                    "State",
				Country:                  "Country",
				Zipcode:                  "Zipcode",
				GstNumber:                "GstNumber",
				NidNumber:                "123456789",
				TradeLicenseUrl:          "TradeLicenseUrl",
				NidFrontImageUrl:         "NidFrontImageUrl",
				NidBackImageUrl:          "NidBackImageUrl",
				ShopOwnerImageUrl:        "ShopOwnerImageUrl",
				GuarantorImageUrl:        "GuarantorImageUrl",
				GuarantorNidNumber:       "12345",
				GuarantorNidBackImageUrl: "GuarantorNidFrontImageUrl",
				ChequeImageUrl:           "ChequeImageUrl",
				CategoryIds:              []uint64{1, 30},
				OpcIds:                   opcIds,
			}
			res, err := new(services.SupplierService).Add(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Supplier Added Successfully"))

			supplier := &models.Supplier{}
			database.DBAPM(ctx).Model(&models.Supplier{}).Where("name = ?", param.Name).Preload("SupplierCategoryMappings").Preload("SupplierOpcMappings").First(&supplier)
			Expect(res.Id).To(Equal(supplier.ID))
			Expect(supplier.Email).To(Equal(param.Email))
			Expect(supplier.SupplierType).To(Equal(utils.Hlc))
			Expect(*supplier.UserID).To(Equal(userId))
			Expect(supplier.Status).To(Equal(models.SupplierStatusPending))
			Expect(supplier.BusinessName).To(Equal(param.BusinessName))
			Expect(supplier.Phone).To(Equal(param.Phone))
			Expect(supplier.AlternatePhone).To(Equal(param.AlternatePhone))
			Expect(supplier.ShopImageURL).To(Equal(param.ShopImageUrl))
			Expect(supplier.NidNumber).To(Equal(param.NidNumber))
			Expect(supplier.NidFrontImageUrl).To(Equal(param.NidFrontImageUrl))
			Expect(supplier.NidBackImageUrl).To(Equal(param.NidBackImageUrl))
			Expect(supplier.TradeLicenseUrl).To(Equal(param.TradeLicenseUrl))
			Expect(supplier.ShopOwnerImageUrl).To(Equal(param.ShopOwnerImageUrl))
			Expect(supplier.GuarantorImageUrl).To(Equal(param.GuarantorImageUrl))
			Expect(supplier.GuarantorNidNumber).To(Equal(param.GuarantorNidNumber))
			Expect(supplier.GuarantorNidBackImageUrl).To(Equal(param.GuarantorNidBackImageUrl))
			Expect(supplier.ChequeImageUrl).To(Equal(param.ChequeImageUrl))

			Expect(len(supplier.SupplierCategoryMappings)).To(Equal(2))
			Expect(supplier.SupplierCategoryMappings[1].CategoryID).To(Equal(uint64(30)))

			Expect(len(supplier.SupplierOpcMappings)).To(Equal(2))
			Expect(supplier.SupplierOpcMappings[1].ProcessingCenterID).To(Equal(uint64(6000)))

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

	Context("Should return error", func() {
		It("When NID number has invalid character", func() {
			param := &supplierpb.SupplierParam{
				Name:             "Name",
				Email:            "Email",
				SupplierType:     uint64(utils.Hlc),
				BusinessName:     "BusinessName",
				Phone:            "8801234567890",
				AlternatePhone:   "8801234567891",
				ShopImageUrl:     "ss2/shop_images/test.png",
				Firstname:        "Firstname",
				Lastname:         "Lastname",
				Address1:         "Address1",
				Address2:         "Address2",
				Landmark:         "Landmark",
				City:             "City",
				State:            "State",
				Country:          "Country",
				Zipcode:          "Zipcode",
				GstNumber:        "GstNumber",
				NidNumber:        "nid_number",
				TradeLicenseUrl:  "TradeLicenseUrl",
				NidFrontImageUrl: "NidFrontImageUrl",
				NidBackImageUrl:  "NidBackImageUrl",
				CategoryIds:      []uint64{1, 30},
			}
			res, err := new(services.SupplierService).Add(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Error while creating Supplier: NID number should only consist of digits"))
		})
	})

	Context("Adding new Supplier without address", func() {
		It("should create supplier and return success", func() {
			param := &supplierpb.SupplierParam{
				Name:         "Name",
				Email:        "Email",
				Phone:        "8801234567890",
				SupplierType: uint64(utils.Hlc),
			}
			res, err := new(services.SupplierService).Add(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Supplier Added Successfully"))

			supplier := &models.Supplier{}
			database.DBAPM(ctx).Model(&models.Supplier{}).Where("id = ?", res.Id).First(&supplier)
			Expect(supplier.Email).To(Equal(param.Email))
		})
	})

	Context("Adding Supplier without name", func() {
		It("Should return error response", func() {
			param := &supplierpb.SupplierParam{
				Email:        "Email",
				SupplierType: uint64(utils.Hlc),
				Phone:        "8801234567890",
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
		It("Should create supplier", func() {
			supplier1 := test_helper.CreateSupplier(ctx, &models.Supplier{SupplierType: utils.Hlc})
			param := &supplierpb.SupplierParam{
				Name:         supplier1.Name,
				Email:        "Email",
				Phone:        "8801234567890",
				SupplierType: uint64(utils.Hlc),
				Address1:     "Address1",
				Zipcode:      "Zipcode",
			}
			res, err := new(services.SupplierService).Add(ctx, param)
			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Supplier Added Successfully"))

			supplier := &models.Supplier{}
			database.DBAPM(ctx).Model(&models.Supplier{}).Where("id = ?", res.Id).First(&supplier)
			Expect(supplier.Name).To(Equal(supplier1.Name))
		})
	})

	Context("Adding Supplier without supplier type", func() {
		It("Should return error response", func() {
			param := &supplierpb.SupplierParam{
				Name:     "Name",
				Email:    "Email",
				Phone:    "8801234567890",
				Address1: "Address1",
				Zipcode:  "Zipcode",
			}
			res, err := new(services.SupplierService).Add(ctx, param)
			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Error while creating Supplier: supplier_type can't be blank"))
		})
	})

	Context("Adding Supplier with OPC Mapping", func() {
		It("Should return error response", func() {
			opcIds := []uint64{5000, 6000}
			mockOpc := mocks.SetOpcMock()
			mockOpc.On("GetProcessingCenterListWithOpcIds", ctx, opcIds).Return(&opcPb.ProcessingCenterListResponse{Data: []*opcPb.OpcDetail{
				{OpcId: 5000},
				{OpcId: 6000},
			}}, nil)

			param := &supplierpb.SupplierParam{
				Name:     "Name",
				Email:    "Email",
				Phone:    "8801234567890",
				Address1: "Address1",
				Zipcode:  "Zipcode",
				OpcIds:   opcIds,
			}
			res, err := new(services.SupplierService).Add(ctx, param)
			supplier := &models.Supplier{}
			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Error while creating Supplier: supplier_type can't be blank"))
			database.DBAPM(ctx).Model(&models.Supplier{}).Where("name = ?", param.Name).Preload("SupplierOpcMappings").First(&supplier)
			Expect(len(supplier.SupplierOpcMappings)).To(Equal(0))
		})

		It("Should return error response for invalid OPC ids", func() {
			opcIds := []uint64{5000, 6000}
			mockOpc := mocks.SetOpcMock()
			mockOpc.On("GetProcessingCenterListWithOpcIds", ctx, opcIds).Return(&opcPb.ProcessingCenterListResponse{Data: []*opcPb.OpcDetail{}}, nil)

			param := &supplierpb.SupplierParam{
				Name:     "Name",
				Email:    "Email",
				Address1: "Address1",
				Zipcode:  "Zipcode",
				OpcIds:   opcIds,
			}
			res, err := new(services.SupplierService).Add(ctx, param)
			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("invalid opc id #(5000)"))
		})
	})

	Context("Adding Supplier by SA user", func() {
		It("Should return with success response", func() {
			opcIds := []uint64{5000, 6000}
			mockOpc := mocks.SetOpcMock()
			mockOpc.On("GetProcessingCenterListWithOpcIds", ctx, opcIds).Return(&opcPb.ProcessingCenterListResponse{Data: []*opcPb.OpcDetail{
				{OpcId: 5000},
				{OpcId: 6000},
			}}, nil)

			mockOpc.On("GetProcessingCenterListWithUserId", ctx, userId).Return(&opcPb.ProcessingCenterListResponse{
				Data: []*opcPb.OpcDetail{
					{OpcId: 201},
					{OpcId: 202},
				},
			}, nil)

			param := &supplierpb.SupplierParam{
				Name:                 "Name",
				Phone:                "8801234567890",
				SupplierType:         uint64(utils.Hlc),
				OpcIds:               opcIds,
				CreateWithOpcMapping: true,
			}
			res, err := new(services.SupplierService).Add(ctx, param)
			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))

			var count int
			database.DBAPM(ctx).Model(&models.SupplierOpcMapping{}).Where("supplier_id = ?", res.Id).Count(&count)
			Expect(count).To(Equal(4))
		})

		It("Should return with success response on OMS remote call error", func() {
			opcIds := []uint64{5000, 6000}
			mockOpc := mocks.SetOpcMock()
			mockOpc.On("GetProcessingCenterListWithUserId", ctx, userId).Return(&opcPb.ProcessingCenterListResponse{}, errors.New("Failing here"))
			mockOpc.On("GetProcessingCenterListWithOpcIds", ctx, opcIds).Return(&opcPb.ProcessingCenterListResponse{Data: []*opcPb.OpcDetail{
				{OpcId: 5000},
				{OpcId: 6000},
			}}, nil)

			param := &supplierpb.SupplierParam{
				Name:                 "Name",
				SupplierType:         uint64(utils.Hlc),
				Phone:                "8801234567890",
				OpcIds:               opcIds,
				CreateWithOpcMapping: true,
			}

			res, err := new(services.SupplierService).Add(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))

			var count int
			database.DBAPM(ctx).Model(&models.SupplierOpcMapping{}).Where("supplier_id = ?", res.Id).Count(&count)
			Expect(count).To(Equal(2))
		})
	})

	Context("Adding Supplier with invalid phone number", func() {
		It("Should return error response", func() {
			param := &supplierpb.SupplierParam{
				Name:         "Name",
				Email:        "Email",
				Phone:        "1234567890",
				SupplierType: uint64(utils.Hlc),
			}
			res, err := new(services.SupplierService).Add(ctx, param)
			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Error while creating Supplier: Invalid Phone Number"))
		})
	})

	Context("Adding Supplier with empty phone number", func() {
		It("Should return error response", func() {
			param := &supplierpb.SupplierParam{
				Name:         "Name",
				Email:        "Email",
				SupplierType: uint64(utils.Hlc),
			}
			res, err := new(services.SupplierService).Add(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Error while creating Supplier: Phone Number can't be blank"))
		})
	})

	Context("Adding Supplier with duplicate phone number", func() {
		It("Should return error response", func() {
			test_helper.CreateSupplier(ctx, &models.Supplier{Phone: "8801234567890"})
			param := &supplierpb.SupplierParam{
				Name:         "Name",
				Email:        "Email",
				SupplierType: uint64(utils.Hlc),
				Phone:        "8801234567890",
			}
			res, err := new(services.SupplierService).Add(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Error while creating Supplier: Supplier Already Exists"))
		})
	})
})
