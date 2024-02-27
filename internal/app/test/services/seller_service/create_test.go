package seller_service_test

import (
	"context"
	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	spb "github.com/voonik/goConnect/api/go/ss2/seller"
	"github.com/voonik/goFramework/pkg/database"
	test_utils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/services"
	"github.com/voonik/ss2/internal/app/test/test_helper"
	"github.com/voonik/ss2/internal/app/utils"
)

var _ = Describe("Create", func() {
	var ctx context.Context

	BeforeEach(func() {
		test_utils.GetContext(&ctx)
	})

	Context("Failure Cases", func() {
		It("Should return error without Seller Params", func() {
			params := spb.CreateParams{}
			res, err := new(services.SellerService).Create(ctx, &params)

			Expect(res.Status).To(Equal(false))
			Expect(res.Message).To(Equal("Missing Seller Params"))
			Expect(err).To(BeNil())
		})
	})
	Context("Success Cases", func() {
		It("Should return success if seller is already registered", func() {
			seller := test_helper.CreateSeller(ctx, &models.Seller{UserID: 101})
			params := spb.CreateParams{Seller: &spb.SellerObject{UserId: 101}}
			res, err := new(services.SellerService).Create(ctx, &params)
			Expect(res.Status).To(Equal(true))
			Expect(res.Message).To(Equal("Seller already registered."))
			Expect(res.UserId).To(Equal(seller.UserID))
			Expect(err).To(BeNil())
		})
		It("Should create seller", func() {
			params := spb.CreateParams{Seller: &spb.SellerObject{
				UserId:                 101,
				BrandName:              "Test Brand",
				CompanyName:            "Test Company",
				PrimaryEmail:           "test@example.com",
				PrimaryPhone:           "1234567890",
				ActivationState:        uint64(utils.ACTIVATED),
				Slug:                   "test-brand",
				Hub:                    "Test Hub",
				Slot:                   "Test Slot",
				DeliveryType:           utils.Ten,
				ProcessingType:         utils.Ten,
				BusinessUnit:           utils.Ten,
				FullfillmentType:       utils.Ten,
				ColorCode:              string(utils.Gold),
				TinNumber:              "123456789",
				SellerCloseDay:         "Friday",
				AcceptedPaymentMethods: "Cash",
			}, AgentId: 7}
			res, err := new(services.SellerService).Create(ctx, &params)

			seller := &models.Seller{UserID: params.Seller.UserId}
			database.DBAPM(ctx).Model(&models.Seller{}).Preload("SellerConfig").Preload("SellerPricingDetails").Find(seller)
			sellerConfig := &models.SellerConfig{SellerID: int(seller.ID)}
			database.DBAPM(ctx).Model(&models.SellerConfig{}).Find(sellerConfig)

			sellerPricing := &models.SellerPricingDetail{}
			database.DBAPM(ctx).Model(&models.SellerPricingDetail{}).Find(sellerPricing)

			Expect(res.Message).To(Equal("Seller registered successfully."))
			Expect(res.UserId).To(Equal(params.Seller.UserId))

			Expect(seller.UserID).To(Equal(params.Seller.UserId))
			Expect(seller.BrandName).To(Equal(params.Seller.BrandName))
			Expect(seller.CompanyName).To(Equal(params.Seller.CompanyName))
			Expect(seller.PrimaryEmail).To(Equal(params.Seller.PrimaryEmail))
			Expect(seller.PrimaryPhone).To(Equal(params.Seller.PrimaryPhone))
			Expect(seller.ActivationState).To(Equal(utils.ActivationState(params.Seller.ActivationState)))
			Expect(seller.Slug).To(Equal(params.Seller.Slug))
			Expect(seller.Hub).To(Equal(params.Seller.Hub))
			Expect(seller.Slot).To(Equal(params.Seller.Slot))
			Expect(seller.DeliveryType).To(Equal(int(params.Seller.DeliveryType)))
			Expect(seller.ProcessingType).To(Equal(int(params.Seller.ProcessingType)))
			Expect(seller.BusinessUnit).To(Equal(int(params.Seller.BusinessUnit)))
			Expect(seller.FullfillmentType).To(Equal(int(params.Seller.FullfillmentType)))
			Expect(seller.ColorCode).To(Equal(utils.ColorCode(params.Seller.ColorCode)))
			Expect(seller.TinNumber).To(Equal(params.Seller.TinNumber))
			Expect(seller.SellerCloseDay).To(Equal(params.Seller.SellerCloseDay))
			Expect(seller.AcceptedPaymentMethods).To(Equal(params.Seller.AcceptedPaymentMethods))
			Expect(seller.AffiliateURL).To(Equal(utils.DefaultAffiliateURL))
			Expect(seller.IsDirect).To(BeTrue())
			Expect(seller.AggregatorID).To(Equal(int(params.Seller.UserId)))
			Expect(seller.AgentID).To(Equal(int(params.AgentId)))

			var returnExchangePolicy map[string]interface{}
			json.Unmarshal(seller.ReturnExchangePolicy, &returnExchangePolicy)

			var policyDetails = map[string]interface{}{
				"return": map[string]interface{}{
					"default_duration":        float64(15),
					"return_days_starts_from": "delivery",
				},
				"exchange": map[string]interface{}{
					"default_duration":        float64(15),
					"return_days_starts_from": "delivery",
				},
			}

			Expect(returnExchangePolicy).To(Equal(policyDetails))

			var dataMappingJson map[string]interface{}
			json.Unmarshal(seller.DataMapping, &dataMappingJson)
			Expect(dataMappingJson).To(Equal(utils.SellerDataMapping))

			Expect(seller.SellerConfig.SellerID).To(Equal(int(seller.ID)))
			Expect(seller.SellerConfig.CODConfirmationNeeded).To(Equal(true))
			Expect(seller.SellerConfig.ItemsPerPackage).To(Equal(int(utils.DefaultSellerItemsPerPackage)))
			Expect(seller.SellerConfig.MaxQuantity).To(Equal(int(utils.DefaultSellerMaxQuantity)))
			Expect(seller.SellerConfig.SellerStockEnabled).To(Equal(true))
			Expect(seller.SellerConfig.AllowPriceUpdate).To(Equal(true))
			Expect(seller.SellerConfig.PickupType).To(Equal(int(utils.DefaultSellerPickupType)))
			Expect(seller.SellerConfig.AllowVendorCoupons).To(Equal(true))

			Expect(seller.SellerPricingDetails[0].SellerID).To(Equal(int(seller.ID)))

			Expect(err).To(BeNil())
		})
	})
})
