package seller_service_test

import (
	"context"
	"encoding/json"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	spb "github.com/voonik/goConnect/api/go/ss2/seller"
	"github.com/voonik/goConnect/api/go/vigeon/notify"
	"github.com/voonik/goFramework/pkg/database"
	test_utils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/services"
	"github.com/voonik/ss2/internal/app/test/mocks"
	"github.com/voonik/ss2/internal/app/test/test_helper"
	"github.com/voonik/ss2/internal/app/utils"
)

var _ = Describe("Create", func() {
	var ctx context.Context
	var mockEmail *mocks.VigeonAPIHelperInterface

	BeforeEach(func() {
		test_utils.GetContext(&ctx)
		mockEmail = mocks.SetVigeonAPIHelperMock()
	})

	AfterEach(func() {
		mocks.UnsetVigeonHelperMock()
	})

	Context("Failure Cases", func() {
		It("Should return error without Seller Params", func() {
			err := fmt.Errorf("Missing All Seller Params")
			expectedResponse := &spb.CreateResponse{Status: false, Message: "Error in seller creation: Missing All Seller Params"}
			mockEmail.On("SendEmailAPI", ctx, notify.EmailParam{
				ToEmail:   "Mokam<noreply@shopf.co>",
				FromEmail: "smk@shopf.co",
				Subject:   "New Seller Registration Failed",
				Content:   fmt.Sprintf("Seller Registration failed because of the error <br><br> %s  <br><br> with the response <br><br> %s", err.Error(), expectedResponse),
			}).Return(&notify.EmailResp{}, nil)

			params := spb.CreateParams{}
			res, err := new(services.SellerService).Create(ctx, &params)

			Expect(res.Status).To(Equal(false))
			Expect(res.Message).To(Equal(expectedResponse.Message))
			Expect(err).To(BeNil())

			Expect(mockEmail.Count["SendEmailAPI"]).To(Equal(1))
		})
		It("Should return error for missing params", func() {
			err := fmt.Errorf("Missing Seller Params: user_id,primary_email,business_unit,brand_name,hub,color_code,activation_state")
			expectedResponse := &spb.CreateResponse{Status: false, Message: "Error in seller creation: Missing Seller Params: user_id,primary_email,business_unit,brand_name,hub,color_code,activation_state"}
			mockEmail.On("SendEmailAPI", ctx, notify.EmailParam{
				ToEmail:   "Mokam<noreply@shopf.co>",
				FromEmail: "smk@shopf.co",
				Subject:   "New Seller Registration Failed",
				Content:   fmt.Sprintf("Seller Registration failed because of the error <br><br> %s  <br><br> with the response <br><br> %s", err.Error(), expectedResponse),
			}).Return(&notify.EmailResp{}, nil)

			params := spb.CreateParams{Seller: &spb.SellerObject{}}
			res, err := new(services.SellerService).Create(ctx, &params)

			Expect(res.Status).To(Equal(false))
			Expect(res.Message).To(Equal(expectedResponse.Message))
			Expect(err).To(BeNil())

			Expect(mockEmail.Count["SendEmailAPI"]).To(Equal(1))
		})
		It("Should return error for invalid params", func() {
			err := fmt.Errorf("Invalid Seller Params: business_unit,color_code,activation_state")
			expectedResponse := &spb.CreateResponse{Status: false, Message: "Error in seller creation: Invalid Seller Params: business_unit,color_code,activation_state"}
			mockEmail.On("SendEmailAPI", ctx, notify.EmailParam{
				ToEmail:   "Mokam<noreply@shopf.co>",
				FromEmail: "smk@shopf.co",
				Subject:   "New Seller Registration Failed",
				Content:   fmt.Sprintf("Seller Registration failed because of the error <br><br> %s  <br><br> with the response <br><br> %s", err.Error(), expectedResponse),
			}).Return(&notify.EmailResp{}, nil)
			params := spb.CreateParams{Seller: &spb.SellerObject{
				UserId:          101,
				PrimaryEmail:    "test@example.com",
				BusinessUnit:    100,
				BrandName:       "Test Brand",
				Hub:             "Test Hub",
				ColorCode:       "InvalidColour",
				ActivationState: 100,
			}}
			res, err := new(services.SellerService).Create(ctx, &params)
			Expect(res.Status).To(Equal(false))
			Expect(res.Message).To(Equal(expectedResponse.Message))
			Expect(err).To(BeNil())

			Expect(mockEmail.Count["SendEmailAPI"]).To(Equal(1))
		})
	})
	Context("Success Cases", func() {
		It("Should return success if seller is already registered", func() {
			seller := test_helper.CreateSeller(ctx, &models.Seller{UserID: 101})
			mockEmail.On("SendEmailAPI", ctx, notify.EmailParam{
				ToEmail:   "Mokam<noreply@shopf.co>",
				FromEmail: "smk@shopf.co",
				Subject:   "New Seller Registered successfully",
				Content:   fmt.Sprintf("Seller Registered (<b>email:</b> %s, <b>agent_email:</b>  ) with the response <br><br> status:true message:\"Seller already registered.\" user_id:101", seller.PrimaryEmail),
			}).Return(&notify.EmailResp{}, nil)
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
			Expect(res.Message).To(Equal("Seller already registered."))
			Expect(res.UserId).To(Equal(seller.UserID))
			Expect(err).To(BeNil())

			Expect(mockEmail.Count["SendEmailAPI"]).To(Equal(1))
		})
		It("Should create seller", func() {
			mockEmail.On("SendEmailAPI", ctx, notify.EmailParam{
				ToEmail:   "Mokam<noreply@shopf.co>",
				FromEmail: "smk@shopf.co",
				Subject:   "New Seller Registered successfully",
				Content:   "Seller Registered (<b>email:</b> test@example.com, <b>agent_email:</b> someEmail@email.com ) with the response <br><br> status:true message:\"Seller registered successfully.\" user_id:101",
			}).Return(&notify.EmailResp{}, nil)
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
				VendorAddresses: []*spb.VendorAddressObject{
					{
						Firstname:            "John",
						Lastname:             "Doe",
						Address1:             "123 Main St",
						Address2:             "Apt 101",
						City:                 "Anytown",
						Zipcode:              "12345",
						AlternativePhone:     "555-1234",
						Company:              "Acme Inc.",
						State:                "AnyState",
						Country:              "AnyCountry",
						AddressType:          1,
						SellerId:             1,
						DefaultAddress:       true,
						AddressProofFileName: "proof.jpg",
						VerificationStatus:   "Verified",
						ExtraData:            "{\"key\":\"value\"}",
						Uuid:                 "123e4567-e89b-12d3-a456-426614174000",
					},
				},
			},
				AgentId:    7,
				AgentEmail: "someEmail@email.com",
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

			Expect(seller.VendorAddresses[0].Firstname).To(Equal(params.Seller.VendorAddresses[0].Firstname))
			Expect(seller.VendorAddresses[0].Lastname).To(Equal(params.Seller.VendorAddresses[0].Lastname))
			Expect(seller.VendorAddresses[0].Address1).To(Equal(params.Seller.VendorAddresses[0].Address1))
			Expect(seller.VendorAddresses[0].Address2).To(Equal(params.Seller.VendorAddresses[0].Address2))
			Expect(seller.VendorAddresses[0].City).To(Equal(params.Seller.VendorAddresses[0].City))
			Expect(seller.VendorAddresses[0].Zipcode).To(Equal(params.Seller.VendorAddresses[0].Zipcode))
			Expect(seller.VendorAddresses[0].AlternativePhone).To(Equal(params.Seller.VendorAddresses[0].AlternativePhone))
			Expect(seller.VendorAddresses[0].Company).To(Equal(params.Seller.VendorAddresses[0].Company))
			Expect(seller.VendorAddresses[0].State).To(Equal(params.Seller.VendorAddresses[0].State))
			Expect(seller.VendorAddresses[0].Country).To(Equal(params.Seller.VendorAddresses[0].Country))
			Expect(seller.VendorAddresses[0].AddressType).To(Equal(int(params.Seller.VendorAddresses[0].AddressType)))
			Expect(seller.VendorAddresses[0].DefaultAddress).To(Equal(params.Seller.VendorAddresses[0].DefaultAddress))
			Expect(seller.VendorAddresses[0].AddressProofFileName).To(Equal(params.Seller.VendorAddresses[0].AddressProofFileName))
			Expect(seller.VendorAddresses[0].VerificationStatus).To(Equal(utils.VerificationStatus(params.Seller.VendorAddresses[0].VerificationStatus)))
			Expect(seller.VendorAddresses[0].ExtraData).To(Equal(params.Seller.VendorAddresses[0].ExtraData))
			Expect(seller.VendorAddresses[0].UUID).To(Equal(params.Seller.VendorAddresses[0].Uuid))

			Expect(err).To(BeNil())

			Expect(mockEmail.Count["SendEmailAPI"]).To(Equal(1))
		})
	})
})
