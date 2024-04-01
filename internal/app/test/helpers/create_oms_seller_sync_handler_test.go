package helper_tests

import (
	"context"
	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/shopuptech/go-jobs/v3/jobs"
	"github.com/stretchr/testify/mock"
	omsPb "github.com/voonik/goConnect/api/go/oms/seller"
	testUtils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/helpers"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/test/mocks"
)

var _ = Describe("CreateOMSSellerSyncHandler", func() {
	var ctx context.Context
	var params *models.Seller
	var apiHelperInstance *mocks.APIHelperInterface

	Context("When valid seller data is passed", func() {

		BeforeEach(func() {
			testUtils.GetContext(&ctx)
			apiHelperInstance = new(mocks.APIHelperInterface)
			helpers.InjectMockAPIHelperInstance(apiHelperInstance)
			vendorAddress := []*models.VendorAddress{{
				Firstname: "John",
				Lastname:  "Doe",
				Address1:  "Dhaka",
				City:      "Dhaka",
			},
				{
					Firstname: "Vendor",
					Lastname:  "2",
					Address1:  "Manyata",
					City:      "Bangalore",
				}}
			sellerConfig := &models.SellerConfig{
				CODConfirmationNeeded: true,
				PenaltyPolicy:         "Penalty",
				ItemsPerPackage:       1,
				PickupType:            1,
				QCFrequency:           1,
			}
			params = &models.Seller{
				UserID:              123,
				FullfillmentType:    1,
				BrandName:           "Sample Brand",
				CompanyName:         "Sample Company",
				Slug:                "sample-slug",
				AggregatorID:        567,
				PrimaryEmail:        "sample@example.com",
				PrimaryPhone:        "1234567890",
				TinNumber:           "ABCDE12345F",
				PanNumber:           "ABCDE1234F",
				SellerInvoiceNumber: 1001,
				SellerType:          2,
				SellerRating:        4.5,
				DeliveryType:        3,
				ProcessingType:      4,
				BusinessUnit:        5,
				VendorAddresses:     vendorAddress,
				SellerConfig:        sellerConfig,
			}
			apiResp := omsPb.SellerResponse{
				Success: true,
				Message: "OMS seller created successfully",
			}
			mockMatcher := mock.MatchedBy(func(omsSellerParams *omsPb.SellerParams) bool {
				return omsSellerParams.Seller.SellerId == params.ID &&
					omsSellerParams.Seller.UserId == params.UserID &&
					omsSellerParams.Seller.BrandName == params.BrandName &&
					omsSellerParams.VendorAddresses[0].Firstname == vendorAddress[0].Firstname &&
					omsSellerParams.VendorAddresses[0].Lastname == vendorAddress[0].Lastname &&
					omsSellerParams.Seller.CodConfirmationNeeded == sellerConfig.CODConfirmationNeeded &&
					omsSellerParams.Seller.PickupType == uint64(sellerConfig.PickupType)
			})
			apiHelperInstance.On("CreateOmsSeller", ctx, mockMatcher).Return(&apiResp)
		})

		It("Should create oms_seller", func() {
			payload, _ := json.Marshal(params)
			err := helpers.CreateOMSSellerSyncHandler(ctx, jobs.Job{Body: payload})
			Expect(err).To(BeNil())
		})
	})
})
