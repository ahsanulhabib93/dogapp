package seller_service_test

import (
	"context"
	"encoding/json"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	spb "github.com/voonik/goConnect/api/go/ss2/seller"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/goFramework/pkg/misc"
	testUtils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/services"
	"github.com/voonik/ss2/internal/app/test/test_helper"
	"github.com/voonik/ss2/internal/app/utils"
)

var _ = Describe("Create", func() {
	var ctx context.Context

	BeforeEach(func() {
		testUtils.GetContext(&ctx)
	})

	Context("Failure Cases", func() {
		It("Should return error without Seller Params", func() {
			expectedResponse := &spb.CreateResponse{Status: false, Message: "Error in seller creation: Missing All Seller Params"}

			params := spb.CreateParams{}
			res, err := new(services.SellerService).Create(ctx, &params)

			Expect(res.Status).To(Equal(false))
			Expect(res.Message).To(Equal(expectedResponse.Message))
			Expect(err).To(BeNil())
		})
		It("Should return error for missing seller params", func() {
			params := spb.CreateParams{Seller: &spb.SellerObject{}}
			res, err := new(services.SellerService).Create(ctx, &params)

			Expect(res.Status).To(Equal(false))
			Expect(res.Message).To(ContainSubstring("activation_state: non zero value required"))
			Expect(res.Message).To(ContainSubstring("brand_name: non zero value required"))
			Expect(res.Message).To(ContainSubstring("business_unit: non zero value required"))
			Expect(res.Message).To(ContainSubstring("color_code: non zero value required"))
			Expect(res.Message).To(ContainSubstring("hub: non zero value required"))
			Expect(res.Message).To(ContainSubstring("primary_email: non zero value required"))
			Expect(res.Message).To(ContainSubstring("user_id: non zero value required"))
			Expect(err).To(BeNil())
		})
		It("Should return error for invalid seller params", func() {
			expectedResponse := &spb.CreateResponse{Status: false, Message: "Error in seller creation: Invalid Seller Params: business_unit,color_code,activation_state"}

			params := spb.CreateParams{Seller: &spb.SellerObject{
				UserId:          101,
				PrimaryEmail:    "test@example.com",
				BusinessUnit:    649,
				BrandName:       "Test Brand",
				Hub:             "Test Hub",
				ColorCode:       "InvalidColour",
				ActivationState: 100,
			}}
			res, err := new(services.SellerService).Create(ctx, &params)
			Expect(res.Status).To(Equal(false))
			Expect(res.Message).To(Equal(expectedResponse.Message))
			Expect(err).To(BeNil())
		})
		It("Should return error for missing vendor address params", func() {
			params := spb.CreateParams{Seller: &spb.SellerObject{
				UserId:          101,
				PrimaryEmail:    "test@example.com",
				BusinessUnit:    uint64(utils.UNICORN),
				BrandName:       "Test Brand",
				Hub:             "Test Hub",
				ColorCode:       string(utils.Gold),
				ActivationState: uint64(utils.ACTIVATED),
				VendorAddresses: []*spb.VendorAddressObject{
					{
						Firstname: utils.EmptyString,
						Address1:  utils.EmptyString,
						Zipcode:   utils.EmptyString,
					},
				},
			}}

			res, err := new(services.SellerService).Create(ctx, &params)
			Expect(res.Status).To(Equal(false))
			Expect(res.Message).To(ContainSubstring("address1: non zero value required"))
			Expect(res.Message).To(ContainSubstring("firstname: non zero value required"))
			Expect(res.Message).To(ContainSubstring("zipcode: non zero value required"))
			Expect(res.Message).To(ContainSubstring("0:"))
			Expect(err).To(BeNil())
		})

		It("Should return error for non unique seller params", func() {
			seller1 := test_helper.CreateSeller(ctx, &models.Seller{UserID: 101, BrandName: "SomeBrand"})

			params := spb.CreateParams{Seller: &spb.SellerObject{
				UserId:           103,
				BrandName:        seller1.BrandName,
				PrimaryEmail:     seller1.PrimaryEmail,
				PrimaryPhone:     seller1.PrimaryPhone,
				ActivationState:  uint64(utils.ACTIVATED),
				Hub:              "Test Hub",
				DeliveryType:     utils.Ten,
				ProcessingType:   utils.Ten,
				BusinessUnit:     utils.Ten,
				FullfillmentType: utils.Ten,
				ColorCode:        string(utils.Gold),
				VendorAddresses: []*spb.VendorAddressObject{
					{
						Firstname: "John",
						Address1:  "123 Main St",
						Zipcode:   "12345",
						SellerId:  1,
					},
				},
			},
			}
			res, err := new(services.SellerService).Create(ctx, &params)
			Expect(res.Status).To(Equal(false))
			fmt.Println(res.Message)
			Expect(res.Message).To(Equal("Error in seller creation: Non Unique Seller Params: primary_email,primary_phone,brand_name"))
			Expect(err).To(BeNil())
		})
	})
	Context("Success Cases", func() {
		It("Should return success if seller is already registered", func() {
			seller := test_helper.CreateSeller(ctx, &models.Seller{UserID: 101})

			params := spb.CreateParams{Seller: &spb.SellerObject{
				UserId:          101,
				PrimaryEmail:    "test@example.com",
				BusinessUnit:    utils.Ten,
				BrandName:       "Test Brand",
				Hub:             "Test Hub",
				ColorCode:       string(utils.Gold),
				ActivationState: uint64(utils.ACTIVATED),
			}}
			res, err := new(services.SellerService).Create(ctx, &params)
			Expect(res.Status).To(Equal(true))
			Expect(res.Message).To(Equal("Seller already registered for UserID: 101"))
			Expect(res.UserId).To(Equal(seller.UserID))
			Expect(err).To(BeNil())
		})

		It("Should create seller", func() {
			ctx = misc.SetInContextThreadObject(ctx, &misc.ThreadObject{VaccountId: 1, PortalId: 1, UserData: &misc.UserData{
				UserId: 1,
				Name:   "Test User",
			}})

			params := spb.CreateParams{Seller: &spb.SellerObject{
				UserId:           101,
				BrandName:        "Test Brand",
				PrimaryEmail:     "test@example.com",
				PrimaryPhone:     "1234567890",
				ActivationState:  uint64(utils.ACTIVATED),
				Hub:              "Test Hub",
				DeliveryType:     utils.Ten,
				ProcessingType:   utils.Ten,
				BusinessUnit:     utils.Ten,
				FullfillmentType: utils.Ten,
				ColorCode:        string(utils.Gold),
				VendorAddresses: []*spb.VendorAddressObject{
					{
						Firstname: "John",
						Address1:  "123 Main St",
						Zipcode:   "12345",
						SellerId:  1,
					},
				},
			},
			}
			res, err := new(services.SellerService).Create(ctx, &params)

			seller := &models.Seller{UserID: params.Seller.UserId}
			database.DBAPM(ctx).Model(&models.Seller{}).Preload("SellerConfig").Preload("SellerPricingDetails").Preload("VendorAddresses").Find(seller)
			sellerConfig := &models.SellerConfig{SellerID: int(seller.ID)}
			database.DBAPM(ctx).Model(&models.SellerConfig{}).Find(sellerConfig)

			sellerPricing := &models.SellerPricingDetail{}
			database.DBAPM(ctx).Model(&models.SellerPricingDetail{}).Find(sellerPricing)

			Expect(res.Message).To(Equal("Seller registered successfully."))
			Expect(res.UserId).To(Equal(params.Seller.UserId))

			Expect(seller.UserID).To(Equal(params.Seller.UserId))
			Expect(seller.BrandName).To(Equal(params.Seller.BrandName))
			Expect(seller.CompanyName).To(Equal(params.Seller.BrandName))
			Expect(seller.PrimaryEmail).To(Equal(params.Seller.PrimaryEmail))
			Expect(seller.PrimaryPhone).To(Equal(params.Seller.PrimaryPhone))
			Expect(seller.ActivationState).To(Equal(utils.ActivationState(params.Seller.ActivationState)))
			Expect(seller.Slug).To(Equal(params.Seller.BrandName))
			Expect(seller.Hub).To(Equal(params.Seller.Hub))
			Expect(seller.DeliveryType).To(Equal(int(params.Seller.DeliveryType)))
			Expect(seller.ProcessingType).To(Equal(int(params.Seller.ProcessingType)))
			Expect(seller.BusinessUnit).To(Equal(utils.BusinessUnit(params.Seller.BusinessUnit)))
			Expect(seller.FullfillmentType).To(Equal(int(params.Seller.FullfillmentType)))
			Expect(seller.ColorCode).To(Equal(utils.ColorCode(params.Seller.ColorCode)))
			Expect(seller.IsDirect).To(BeTrue())
			Expect(seller.AggregatorID).To(Equal(int(params.Seller.UserId)))
			Expect(seller.AgentID).To(Equal(int(misc.ExtractThreadObject(ctx).GetUserData().GetUserId())))

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

			Expect(seller.VendorAddresses[0].Firstname).To(Equal(params.Seller.VendorAddresses[0].Firstname))
			Expect(seller.VendorAddresses[0].Address1).To(Equal(params.Seller.VendorAddresses[0].Address1))
			Expect(seller.VendorAddresses[0].Zipcode).To(Equal(params.Seller.VendorAddresses[0].Zipcode))
			Expect(seller.VendorAddresses[0].State).To(Equal(utils.DefaultState))
			Expect(seller.VendorAddresses[0].Country).To(Equal(utils.DefaultCountry))
			Expect(seller.VendorAddresses[0].AddressType).To(Equal(2))
			Expect(seller.VendorAddresses[0].SellerID).To(Equal(int(seller.ID)))
			Expect(seller.VendorAddresses[0].UUID).To(Not(BeNil()))
			Expect(seller.VendorAddresses[0].VerificationStatus).To(Equal(utils.Verified))

			Expect(err).To(BeNil())
		})
	})
})
