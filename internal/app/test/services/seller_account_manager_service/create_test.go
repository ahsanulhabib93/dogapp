package seller_account_manager_service_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/voonik/goFramework/pkg/database"

	sampb "github.com/voonik/goConnect/api/go/ss2/seller_account_manager"
	testUtils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/services"
	"github.com/voonik/ss2/internal/app/test/test_helper"
	"github.com/voonik/ss2/internal/app/utils"
)

var _ = Describe("SellerAccountManager Create", func() {
	var ctx context.Context
	var seller *models.Seller

	BeforeEach(func() {
		testUtils.GetContext(&ctx)
		seller = test_helper.CreateSeller(ctx, &models.Seller{
			UserID:    123456,
			BrandName: "test_brand",
		})
	})

	Context("Failure Cases", func() {
		var createParams *sampb.AccountManagerObject
		BeforeEach(func() {
			createParams = &sampb.AccountManagerObject{
				Phone:    8801548654343,
				Role:     "KAM",
				SellerId: uint64(999888),
				Name:     "TEST_SAM",
			}
		})
		It("Should return error on invalid phone", func() {
			createParams.Phone = uint64(8801435)
			resp, err := new(services.SellerAccountManagerService).Create(ctx, createParams)
			Expect(err).To(BeNil())
			Expect(resp.Success).To(BeFalse())
			Expect(resp.Message).To(ContainSubstring("Invalid Seller:"))
		})
		It("Should return error on empty role", func() {
			createParams.Role = ""
			resp, err := new(services.SellerAccountManagerService).Create(ctx, createParams)
			Expect(err).To(BeNil())
			Expect(resp.Success).To(BeFalse())
			Expect(resp.Message).To(ContainSubstring("Invalid Seller:"))
		})
		It("Should return error on invalid seller", func() {
			createParams.SellerId = uint64(999888)
			resp, err := new(services.SellerAccountManagerService).Create(ctx, createParams)
			Expect(err).To(BeNil())
			Expect(resp.Success).To(BeFalse())
			Expect(resp.Message).To(ContainSubstring("Invalid Seller:"))
		})
	})

	Context("Succes Cases", func() {
		It("Should return success", func() {
			resp, err := new(services.SellerAccountManagerService).Create(ctx, &sampb.AccountManagerObject{
				Phone:    8801548654343,
				Role:     "KAM",
				SellerId: seller.ID,
				Name:     "TEST_SAM",
			})
			Expect(err).To(BeNil())
			Expect(resp.Success).To(BeTrue())
			Expect(resp.Message).To(Equal("SellerAccountManager added successfully"))
			Expect(resp.SellerUserId).To(Equal(seller.UserID))
			updatedSam := &models.SellerAccountManager{}
			database.DBAPM(ctx).Model(&models.SellerAccountManager{}).Where("seller_id = ? ", seller.ID).Find(updatedSam)
			Expect(updatedSam.Phone).To(Equal(int64(8801548654343)))
			Expect(updatedSam.Role).To(Equal("KAM"))
			Expect(updatedSam.Name).To(Equal("TEST_SAM"))
			Expect(updatedSam.Priority).To(Equal(utils.One))
		})
	})
})
