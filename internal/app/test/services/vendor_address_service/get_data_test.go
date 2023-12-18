package vendor_address_service_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	vapb "github.com/voonik/goConnect/api/go/ss2/vendor_address"
	aaaModels "github.com/voonik/goFramework/pkg/aaa/models"
	test_utils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/helpers"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/services"
	"github.com/voonik/ss2/internal/app/test/mocks"
	"github.com/voonik/ss2/internal/app/test/test_helper"
)

var _ = Describe("Get Data", func() {
	var ctx context.Context
	var mockAudit *mocks.AuditLogMock

	BeforeEach(func() {
		test_utils.GetContext(&ctx)
		mocks.UnsetOpcMock()

		mockAudit = mocks.SetAuditLogMock()
		mockAudit.On("RecordAuditAction", ctx, mock.Anything).Return(nil)
	})

	AfterEach(func() {
		mocks.UnsetAuditLogMock()
		helpers.InjectMockAPIHelperInstance(nil)
		helpers.InjectMockIdentityUserApiHelperInstance(nil)
		aaaModels.InjectMockAppPreferenceServiceInstance(nil)
	})

	Context("When no params are given", func() {
		It("Should return error", func() {
			param := vapb.GetDataParams{}
			res, err := new(services.VendorAddressService).GetData(ctx, &param)
			Expect(res.Status).To(Equal("failure"))
			Expect(res.Message).To(Equal("param not specified"))
			Expect(err).To(BeNil())
		})
	})

	Context("When seller id is passed in the param", func() {
		It("Should return vendor address details", func() {
			seller1 := test_helper.CreateSeller(ctx, &models.Seller{})
			seller2 := test_helper.CreateSeller(ctx, &models.Seller{})
			vendorAddress1 := test_helper.CreateVendorAddress(ctx, &models.VendorAddress{SellerID: int(seller1.ID)})
			vendorAddress2 := test_helper.CreateVendorAddress(ctx, &models.VendorAddress{SellerID: int(seller2.ID)})
			param := vapb.GetDataParams{
				SellerIds: []uint64{seller1.ID, seller2.ID},
			}
			res, err := new(services.VendorAddressService).GetData(ctx, &param)
			Expect(len(res.VendorAddress)).To(Equal(2))
			Expect(res.VendorAddress[0].Firstname).To(Equal(vendorAddress1.Firstname))
			Expect(res.VendorAddress[0].Lastname).To(Equal(vendorAddress1.Lastname))
			Expect(res.VendorAddress[1].Firstname).To(Equal(vendorAddress2.Firstname))
			Expect(res.VendorAddress[1].Lastname).To(Equal(vendorAddress2.Lastname))
			Expect(res.Status).To(Equal("success"))
			Expect(res.Message).To(Equal("fetched vendor address successfully"))
			Expect(err).To(BeNil())
		})
	})

	Context("When vendor address id is passed in the param", func() {
		It("Should return vendor address details", func() {
			vendorAddress1 := test_helper.CreateVendorAddress(ctx, &models.VendorAddress{})
			vendorAddress2 := test_helper.CreateVendorAddress(ctx, &models.VendorAddress{})
			param := vapb.GetDataParams{
				Ids: []uint64{vendorAddress1.ID, vendorAddress2.ID},
			}

			res, err := new(services.VendorAddressService).GetData(ctx, &param)
			Expect(len(res.VendorAddress)).To(Equal(2))
			Expect(res.VendorAddress[0].Firstname).To(Equal(vendorAddress1.Firstname))
			Expect(res.VendorAddress[0].Lastname).To(Equal(vendorAddress1.Lastname))
			Expect(res.VendorAddress[1].Firstname).To(Equal(vendorAddress2.Firstname))
			Expect(res.VendorAddress[1].Lastname).To(Equal(vendorAddress2.Lastname))
			Expect(res.Status).To(Equal("success"))
			Expect(res.Message).To(Equal("fetched vendor address successfully"))
			Expect(err).To(BeNil())
		})
	})

	Context("When user id is passed in the param", func() {
		It("Should return vendor address details", func() {
			seller1 := test_helper.CreateSeller(ctx, &models.Seller{})
			seller2 := test_helper.CreateSeller(ctx, &models.Seller{})
			vendorAddress1 := test_helper.CreateVendorAddress(ctx, &models.VendorAddress{SellerID: int(seller1.ID)})
			vendorAddress2 := test_helper.CreateVendorAddress(ctx, &models.VendorAddress{SellerID: int(seller2.ID)})
			param := vapb.GetDataParams{
				UserIds: []uint64{seller1.UserID, seller2.UserID},
			}

			res, err := new(services.VendorAddressService).GetData(ctx, &param)
			Expect(len(res.VendorAddress)).To(Equal(2))
			Expect(res.VendorAddress[0].Firstname).To(Equal(vendorAddress1.Firstname))
			Expect(res.VendorAddress[0].Lastname).To(Equal(vendorAddress1.Lastname))
			Expect(res.VendorAddress[1].Firstname).To(Equal(vendorAddress2.Firstname))
			Expect(res.VendorAddress[1].Lastname).To(Equal(vendorAddress2.Lastname))
			Expect(res.Status).To(Equal("success"))
			Expect(res.Message).To(Equal("fetched vendor address successfully"))
			Expect(err).To(BeNil())
		})
	})

	Context("When multiple params are passed in the param", func() {
		It("Should return error", func() {
			seller1 := test_helper.CreateSeller(ctx, &models.Seller{})
			seller2 := test_helper.CreateSeller(ctx, &models.Seller{UserID: 101})
			test_helper.CreateVendorAddress(ctx, &models.VendorAddress{SellerID: int(seller1.ID)})
			test_helper.CreateVendorAddress(ctx, &models.VendorAddress{SellerID: int(seller2.ID)})
			vendorAddress3 := test_helper.CreateVendorAddress(ctx, &models.VendorAddress{})
			param := vapb.GetDataParams{
				Ids:       []uint64{vendorAddress3.ID},
				UserIds:   []uint64{seller1.UserID},
				SellerIds: []uint64{seller2.ID},
			}

			res, err := new(services.VendorAddressService).GetData(ctx, &param)
			Expect(res.Status).To(Equal("failure"))
			Expect(res.Message).To(Equal("specify any one param"))
			Expect(err).To(BeNil())
		})

		It("Should return error", func() {
			seller1 := test_helper.CreateSeller(ctx, &models.Seller{})
			seller2 := test_helper.CreateSeller(ctx, &models.Seller{UserID: 101})
			test_helper.CreateVendorAddress(ctx, &models.VendorAddress{SellerID: int(seller1.ID)})
			test_helper.CreateVendorAddress(ctx, &models.VendorAddress{SellerID: int(seller2.ID)})
			param := vapb.GetDataParams{
				UserIds:   []uint64{seller1.UserID},
				SellerIds: []uint64{seller2.ID},
			}

			res, err := new(services.VendorAddressService).GetData(ctx, &param)
			Expect(res.Status).To(Equal("failure"))
			Expect(res.Message).To(Equal("specify any one param"))
			Expect(err).To(BeNil())
		})
	})

	Context("When invalid user_id is passed in the param", func() {
		It("Should return error", func() {
			param := vapb.GetDataParams{
				UserIds: []uint64{101},
			}
			res, err := new(services.VendorAddressService).GetData(ctx, &param)
			Expect(res.Status).To(Equal("failure"))
			Expect(res.Message).To(Equal("seller not found with the user id"))
			Expect(err).To(BeNil())
		})
	})

	Context("When invalid user_id is passed in the param", func() {
		It("Should return error", func() {
			seller1 := test_helper.CreateSeller(ctx, &models.Seller{})
			param := vapb.GetDataParams{
				UserIds: []uint64{seller1.UserID},
			}
			res, err := new(services.VendorAddressService).GetData(ctx, &param)
			Expect(res.Status).To(Equal("success"))
			Expect(res.Message).To(Equal("vendor address not found"))
			Expect(err).To(BeNil())
		})
	})
})
