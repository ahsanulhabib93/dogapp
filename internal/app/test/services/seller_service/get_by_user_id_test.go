package seller_service_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	spb "github.com/voonik/goConnect/api/go/ss2/seller"
	aaaModels "github.com/voonik/goFramework/pkg/aaa/models"
	"github.com/voonik/goFramework/pkg/database"
	test_utils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/helpers"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/services"
	"github.com/voonik/ss2/internal/app/test/mocks"
	"github.com/voonik/ss2/internal/app/test/test_helper"
	"github.com/voonik/ss2/internal/app/utils"
)

var _ = Describe("Get seller by user ID", func() {
	var ctx context.Context
	var mockAudit *mocks.AuditLogMock
	userId := uint64(101)

	BeforeEach(func() {
		test_utils.GetContext(&ctx)
		mocks.UnsetOpcMock()

		test_helper.SetContextUser(&ctx, userId, []string{})
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
			param := spb.GetByUserIDParams{}
			res, err := new(services.SellerService).GetByUserID(ctx, &param)
			Expect(res).To(Equal(&spb.GetByUserIDResponse{}))
			Expect(err).To(BeNil())
		})
	})

	Context("When params are given", func() {
		It("Should return seller details", func() {
			seller := test_helper.CreateSeller(ctx, &models.Seller{
				UserID:    userId,
				BrandName: "test_brand",
			})

			seller.VendorAddresses = []*models.VendorAddress{{SellerID: int(seller.ID), Zipcode: "456789", Address1: "Addr-456789", VerificationStatus: "VERIFIED"}}
			seller.BusinessType = utils.Manufacturer
			seller.ColorCode = utils.Black
			database.DBAPM(ctx).Save(&seller)

			param := spb.GetByUserIDParams{
				UserId: userId,
			}
			res, err := new(services.SellerService).GetByUserID(ctx, &param)
			Expect(res.Seller.BrandName).To(Equal(seller.BrandName))
			Expect(res.Seller.UserId).To(Equal(seller.UserID))
			Expect(res.Seller.VendorAddresses[0].Zipcode).To(Equal("456789"))
			Expect(res.Seller.VendorAddresses[0].Address1).To(Equal("Addr-456789"))
			Expect(res.Seller.ReturnExchangePolicy.Return.GetReturnDaysStartsFrom()).To(Equal("delivery"))
			Expect(res.Seller.SellerConfig.MaxQuantity).To(Equal(uint64(1000)))
			Expect(err).To(BeNil())
		})
	})

	Context("When params are given", func() {
		It("Should return seller details", func() {
			param := spb.GetByUserIDParams{
				UserId: 123,
			}
			res, err := new(services.SellerService).GetByUserID(ctx, &param)
			Expect(res).To(Equal(&spb.GetByUserIDResponse{}))
			Expect(err).To(BeNil())
		})
	})
})
