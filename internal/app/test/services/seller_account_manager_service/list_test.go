package seller_account_manager_service_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	sampb "github.com/voonik/goConnect/api/go/ss2/seller_account_manager"
	test_utils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/services"
	"github.com/voonik/ss2/internal/app/test/test_helper"
	"github.com/voonik/ss2/internal/app/utils"
)

var _ = Describe("SellerAccountManager", func() {
	var ctx context.Context

	BeforeEach(func() {
		test_utils.GetContext(&ctx)
	})

	Context("List", func() {
		It("Should return empty response for no params", func() {
			resp, err := new(services.SellerAccountManagerService).List(ctx, &sampb.ListParams{})

			Expect(err).To(BeNil())
			Expect(resp.Status).To(Equal(utils.EmptyString))
			Expect(resp.AccountManager).To(BeNil())
		})

		Context("With Data", func() {
			var Sam1 *models.SellerAccountManager
			BeforeEach(func() {
				Sam1 = test_helper.CreateSellerAccountManager(ctx, 1, "SAM 1", 98765, "example@example.com", 1, "sourcing-associate")
			})
			It("Should return data with success message for proper params", func() {

				resp, err := new(services.SellerAccountManagerService).List(ctx, &sampb.ListParams{SellerId: 1})

				Expect(err).To(BeNil())
				Expect(resp.Status).To(Equal("Success"))
				Expect(resp.AccountManager.Id).To(Equal(Sam1.ID))
				Expect(resp.AccountManager.Email).To(Equal(Sam1.Email))
				Expect(resp.AccountManager.Phone).To(Equal(uint64(Sam1.Phone)))
				Expect(resp.AccountManager.Name).To(Equal(Sam1.Name))
				Expect(resp.AccountManager.Priority).To(Equal(uint64(Sam1.Priority)))
				Expect(resp.AccountManager.Role).To(Equal(Sam1.Role))
			})
		})
	})
})
