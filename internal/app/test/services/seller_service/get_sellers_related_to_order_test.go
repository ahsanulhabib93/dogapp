package seller_service_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	spb "github.com/voonik/goConnect/api/go/ss2/seller"
	aaaModels "github.com/voonik/goFramework/pkg/aaa/models"
	test_utils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/helpers"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/services"
	"github.com/voonik/ss2/internal/app/test/mocks"
	"github.com/voonik/ss2/internal/app/test/test_helper"
)

var _ = Describe("Get Sellers Related To Order", func() {
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

	Context("When user id is not passed as param", func() {
		It("Should return error", func() {

			param := spb.GetSellersRelatedToOrderParams{}
			res, err := new(services.SellerService).GetSellersRelatedToOrder(ctx, &param)
			Expect(len(res.Seller)).To(Equal(0))
			Expect(res.Status).To(Equal("failure"))
			Expect(res.Message).To(Equal("no valid param"))
			Expect(err).To(BeNil())
		})
	})

	Context("When invalid user id is passed as param", func() {
		It("Should return error", func() {

			param := spb.GetSellersRelatedToOrderParams{
				SellerIds: []uint64{101},
			}

			res, err := new(services.SellerService).GetSellersRelatedToOrder(ctx, &param)
			Expect(len(res.Seller)).To(Equal(0))
			Expect(res.Status).To(Equal("success"))
			Expect(res.Message).To(Equal("seller not found"))
			Expect(err).To(BeNil())
		})
	})

	Context("When user id is passed as param", func() {
		It("Should return seller data", func() {
			seller1 := test_helper.CreateSeller(ctx, &models.Seller{
				UserID:    uint64(101),
				BrandName: "test_brand",
			})
			param := spb.GetSellersRelatedToOrderParams{
				SellerIds: []uint64{101},
			}

			res, err := new(services.SellerService).GetSellersRelatedToOrder(ctx, &param)
			Expect(len(res.Seller)).To(Equal(1))
			Expect(res.Seller[0].UserId).To(Equal(seller1.UserID))
			Expect(res.Seller[0].BrandName).To(Equal(seller1.BrandName))
			Expect(res.Status).To(Equal("success"))
			Expect(res.Message).To(Equal("fetched seller details successfully"))
			Expect(err).To(BeNil())
		})
	})

	Context("When multiple user ids are passed as param", func() {
		It("Should return seller data", func() {
			seller1 := test_helper.CreateSeller(ctx, &models.Seller{
				UserID:    uint64(100),
				BrandName: "test_brand1",
			})
			seller2 := test_helper.CreateSeller(ctx, &models.Seller{
				UserID:    uint64(102),
				BrandName: "test_brand2",
			})
			param := spb.GetSellersRelatedToOrderParams{
				SellerIds: []uint64{100, 102},
			}

			res, err := new(services.SellerService).GetSellersRelatedToOrder(ctx, &param)
			Expect(len(res.Seller)).To(Equal(2))
			Expect(res.Seller[0].UserId).To(Equal(seller1.UserID))
			Expect(res.Seller[0].BrandName).To(Equal(seller1.BrandName))
			Expect(res.Seller[1].UserId).To(Equal(seller2.UserID))
			Expect(res.Seller[1].BrandName).To(Equal(seller2.BrandName))
			Expect(res.Status).To(Equal("success"))
			Expect(res.Message).To(Equal("fetched seller details successfully"))
			Expect(err).To(BeNil())
		})
	})

	Context("When multiple user ids are passed as param", func() {
		It("Should return seller data for valid user ids", func() {
			seller1 := test_helper.CreateSeller(ctx, &models.Seller{
				UserID:    uint64(100),
				BrandName: "test_brand1",
			})
			param := spb.GetSellersRelatedToOrderParams{
				SellerIds: []uint64{100, 103},
			}

			res, err := new(services.SellerService).GetSellersRelatedToOrder(ctx, &param)
			Expect(len(res.Seller)).To(Equal(1))
			Expect(res.Seller[0].UserId).To(Equal(seller1.UserID))
			Expect(res.Seller[0].BrandName).To(Equal(seller1.BrandName))
			Expect(res.Seller[0].ReturnExchangePolicy.Return.GetReturnDaysStartsFrom()).To(Equal("delivery"))
			Expect(res.Seller[0].SellerConfig.MaxQuantity).To(Equal(uint64(1000)))
			Expect(res.Status).To(Equal("success"))
			Expect(res.Message).To(Equal("fetched seller details successfully"))
			Expect(err).To(BeNil())
		})
	})

})
