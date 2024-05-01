package seller_service_test

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	spb "github.com/voonik/goConnect/api/go/ss2/seller"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/goFramework/pkg/misc"
	testUtils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/services"
)

var _ = Describe("Index", func() {
	var ctx context.Context
	var seller *models.Seller
	var seller2 *models.Seller
	BeforeEach(func() {
		testUtils.GetContext(&ctx)
		seller = &models.Seller{
			UserID:           101,
			BusinessUnit:     1,
			BrandName:        "testBrand",
			FullfillmentType: 2,
		}
		seller2 = &models.Seller{
			UserID:           102,
			BusinessUnit:     1,
			BrandName:        "testBrand",
			FullfillmentType: 1,
		}
		database.DBAPM(ctx).Model(&models.Seller{}).Create(seller)
		database.DBAPM(ctx).Model(&models.Seller{}).Create(seller2)
	})
	Context("Success cases", func() {
		Context("When UserID is passed", func() {
			It("Should filter sellers with userID even if current user is present", func() {
				ctx = misc.SetInContextThreadObject(ctx, &misc.ThreadObject{VaccountId: 1, PortalId: 1, UserData: &misc.UserData{
					UserId: 1,
					Name:   "Test User",
					Phone:  "8801485743298",
				}})
				resp, err := new(services.SellerService).Index(ctx, &spb.GetSellerParams{
					UserId: []uint64{seller2.UserID},
				})
				Expect(err).To(BeNil())
				Expect(resp.Status).To(Equal("success"))
				Expect(resp.Seller).To(HaveLen(1))
				sellerData := resp.Seller[0]
				Expect(sellerData.Id).To(Equal(seller2.ID))
				Expect(sellerData.UserId).To(Equal(seller2.UserID))
				Expect(sellerData.BusinessUnit).To(Equal(uint64(seller2.BusinessUnit)))
			})
		})
		Context("When Current User is not present and userID is not passed", func() {
			It("Should filter seller with fulfilment type as 2", func() {
				resp, err := new(services.SellerService).Index(ctx, &spb.GetSellerParams{
					BusinessUnits: []uint64{uint64(seller.BusinessUnit)},
					BrandName:     seller.BrandName,
				})
				Expect(err).To(BeNil())
				Expect(resp.Status).To(Equal("success"))
				Expect(resp.Seller).To(HaveLen(1))
				sellerData := resp.Seller[0]
				Expect(sellerData.Id).To(Equal(seller.ID))
				Expect(sellerData.BusinessUnit).To(Equal(uint64(seller.BusinessUnit)))
			})
		})
		Context("When Current User is present and userID is not passed", func() {
			var sam *models.SellerAccountManager
			BeforeEach(func() {
				sam = &models.SellerAccountManager{
					Phone:    8801485743298,
					Role:     "sourcing_associate",
					SellerID: seller2.ID,
					Name:     "testSam",
				}
				database.DBAPM(ctx).Model(&models.SellerAccountManager{}).Create(sam)
			})
			It("Should return Sellers Mapped to Current user i.e (SAM)", func() {
				ctx = misc.SetInContextThreadObject(ctx, &misc.ThreadObject{VaccountId: 1, PortalId: 1, UserData: &misc.UserData{
					UserId: 1,
					Name:   "Test User",
					Phone:  "8801485743298",
				}})
				resp, err := new(services.SellerService).Index(ctx, &spb.GetSellerParams{})
				Expect(err).To(BeNil())
				Expect(resp.Status).To(Equal("success"))
				Expect(resp.Seller).To(HaveLen(1))
				sellerData := resp.Seller[0]
				Expect(sellerData.Id).To(Equal(seller2.ID))
				Expect(sellerData.BusinessUnit).To(Equal(uint64(seller2.BusinessUnit)))
			})
		})
	})
	Context("Failure cases", func() {
		Context("when SAM select query fails", func() {
			It("Should return error", func() {
				ctx = misc.SetInContextThreadObject(ctx, &misc.ThreadObject{VaccountId: 1, PortalId: 1, UserData: &misc.UserData{
					UserId: 1,
					Name:   "Test User",
					Phone:  "8801485743298",
				}})
				database.DBAPM(ctx).Model(&models.SellerAccountManager{}).DropColumn("role")
				resp, err := new(services.SellerService).Index(ctx, &spb.GetSellerParams{})
				fmt.Println(resp)
				database.DBAPM(ctx).AutoMigrate(&models.SellerAccountManager{})
				Expect(err).To(BeNil())
				Expect(resp.Status).To(Equal("failure"))
				Expect(resp.Message).To(ContainSubstring("Unknown column 'seller_account_managers.role' in 'where clause'"))
			})
		})
		Context("when Seller select query fails", func() {
			It("Should return error", func() {
				database.DBAPM(ctx).Model(&models.Seller{}).DropColumn("vaccount_id")
				resp, err := new(services.SellerService).Index(ctx, &spb.GetSellerParams{})
				fmt.Println(resp)
				database.DBAPM(ctx).AutoMigrate(&models.Seller{})
				Expect(err).To(BeNil())
				Expect(resp.Status).To(Equal("failure"))
				Expect(resp.Message).To(ContainSubstring("Unknown column 'sellers.vaccount_id' in 'where clause'"))
			})
		})
	})
})
