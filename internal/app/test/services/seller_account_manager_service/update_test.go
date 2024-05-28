package seller_account_manager_service_test

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/voonik/goFramework/pkg/database"

	sampb "github.com/voonik/goConnect/api/go/ss2/seller_account_manager"
	testUtils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/services"
	"github.com/voonik/ss2/internal/app/test/test_helper"
)

var _ = Describe("SellerAccountManager Update", func() {
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
		It("Should return if id is empty in params", func() {
			resp, err := new(services.SellerAccountManagerService).Update(ctx, &sampb.AccountManagerObject{})
			Expect(err).To(BeNil())
			Expect(resp.Success).To(BeFalse())
			Expect(resp.Message).To(Equal("id cannot be empty"))
		})
		It("Should return if id is invalid", func() {
			resp, err := new(services.SellerAccountManagerService).Update(ctx, &sampb.AccountManagerObject{Id: uint64(9999)})
			Expect(err).To(BeNil())
			Expect(resp.Success).To(BeFalse())
			fmt.Println(resp.Message)
			Expect(resp.Message).To(Equal("record not found"))
		})
	})

	Context("Succes Cases", func() {
		It("Should return success on update", func() {
			sam := test_helper.CreateSellerAccountManager(ctx, seller.ID, "SAM 1", 8801548654342, "example@example.com", 1, "sourcing-associate")
			resp, err := new(services.SellerAccountManagerService).Update(ctx, &sampb.AccountManagerObject{
				Id:    sam.ID,
				Phone: 8801548654343,
				Email: "example@samm.com",
			})
			Expect(err).To(BeNil())
			Expect(resp.Success).To(BeTrue())
			Expect(resp.SellerUserId).To(Equal(seller.UserID))
			Expect(resp.Message).To(Equal("update successfull"))
			updatedSam := &models.SellerAccountManager{}
			database.DBAPM(ctx).Model(&models.SellerAccountManager{}).Where("id = ? ", sam.ID).Find(updatedSam)
			Expect(updatedSam.Phone).To(Equal(int64(8801548654343)))
			Expect(updatedSam.Email).To(Equal("example@samm.com"))
			Expect(updatedSam.Role).To(Equal("sourcing-associate"))
			Expect(updatedSam.Name).To(Equal("SAM 1"))
		})
	})
})
