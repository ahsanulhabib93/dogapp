package seller_service_test

import (
	"context"
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	userMappingPb "github.com/voonik/goConnect/api/go/oms/user_mapping"
	spb "github.com/voonik/goConnect/api/go/ss2/seller"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/goFramework/pkg/misc"
	testUtils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/helpers"
	"github.com/voonik/ss2/internal/app/helpers/mocks"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/services"
)

var _ = Describe("SellerByBrandName", func() {
	var ctx context.Context
	var seller *models.Seller
	var seller2 *models.Seller
	var omsApiMock *mocks.OmsApiHelperInterface
	var omsResponse *userMappingPb.UserMappingResponse
	BeforeEach(func() {
		testUtils.GetContext(&ctx)
		omsResponse = &userMappingPb.UserMappingResponse{
			Data: []*userMappingPb.UserMappingData{
				{
					UserId:        999,
					BusinessUnits: []uint64{1},
					OpcIds:        []uint64{1},
					ZoneIds:       []uint64{1},
					UserData: &userMappingPb.UserData{
						Roles: []uint64{1},
					},
				},
			},
		}
		ctx = misc.SetInContextThreadObject(ctx, &misc.ThreadObject{VaccountId: 1, PortalId: 1, UserData: &misc.UserData{
			UserId: 999,
			Name:   "Test User",
			Phone:  "8801485743298",
		}})
		seller = &models.Seller{
			UserID:           101,
			BusinessUnit:     1,
			BrandName:        "testBrand",
			FullfillmentType: 2,
		}
		seller2 = &models.Seller{
			UserID:           102,
			BusinessUnit:     2,
			BrandName:        "testBrand",
			FullfillmentType: 1,
		}
		database.DBAPM(ctx).Model(&models.Seller{}).Create(seller)
		database.DBAPM(ctx).Model(&models.Seller{}).Create(seller2)
		omsApiMock = new(mocks.OmsApiHelperInterface)
		helpers.InjectMockOmsAPIHelperInstance(omsApiMock)
	})
	AfterEach(func() {
		helpers.InjectMockOmsAPIHelperInstance(nil)
	})
	Context("Success case", func() {
		BeforeEach(func() {
			omsApiMock.On("FetchUserMappingData", ctx, []uint64{999}).Return(omsResponse, nil)
		})
		It("Should filter seller with bu based on OMS userMapping Data", func() {
			resp, err := new(services.SellerService).SellerByBrandName(ctx, &spb.GetSellerParams{
				BrandName: seller.BrandName,
			})
			Expect(err).To(BeNil())
			Expect(resp.Status).To(Equal("success"))
			Expect(resp.Seller).To(HaveLen(1))
			sellerData := resp.Seller[0]
			Expect(sellerData.Id).To(Equal(seller.ID))
			Expect(sellerData.BusinessUnit).To(Equal(uint64(seller.BusinessUnit)))
		})
	})
	Context("Failure case", func() {
		Context("When Oms Api Fails", func() {
			BeforeEach(func() {
				omsApiMock.On("FetchUserMappingData", ctx, []uint64{999}).Return(nil, errors.New("oms api failed"))
			})
			It("Should return error", func() {
				resp, err := new(services.SellerService).SellerByBrandName(ctx, &spb.GetSellerParams{
					BrandName: seller.BrandName,
				})
				Expect(err).To(BeNil())
				Expect(resp.Status).To(Equal("failure"))
				Expect(resp.Message).To(Equal("oms api failed"))
				Expect(resp.Seller).To(HaveLen(0))
			})
		})
		Context("When DB select query fails", func() {
			BeforeEach(func() {
				omsApiMock.On("FetchUserMappingData", ctx, []uint64{999}).Return(omsResponse, nil)
			})
			It("Should return error", func() {
				database.DBAPM(ctx).Model(&models.Seller{}).DropColumn("brand_name")
				resp, err := new(services.SellerService).SellerByBrandName(ctx, &spb.GetSellerParams{
					BrandName: seller.BrandName,
				})
				database.DBAPM(ctx).AutoMigrate(&models.Seller{})
				fmt.Println("BRANDDDDD")
				fmt.Println(resp)
				Expect(err).To(BeNil())
				Expect(resp.Status).To(Equal("failure"))
				Expect(resp.Message).To(ContainSubstring("Unknown column 'brand_name' in 'where clause'"))
				Expect(resp.Seller).To(HaveLen(0))
			})
		})
	})
})
