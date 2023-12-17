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
			var Sam1, Sam2, Sam3 *models.SellerAccountManager
			BeforeEach(func() {
				Sam1 = test_helper.CreateSellerAccountManager(ctx, 1, "SAM 1", 98765, "example@example.com", 1, "sourcing-associate")
				Sam2 = test_helper.CreateSellerAccountManager(ctx, 1, "SAM 2", 98766, "example2@example.com", 1, "non-sourcing-associate")
				Sam3 = test_helper.CreateSellerAccountManager(ctx, 1, "SAM 3", 98766, "example2@example.com", 2, "non-sourcing-associate")
			})
			It("Should return data ordered by role and priority with success message for proper params", func() {

				resp, err := new(services.SellerAccountManagerService).List(ctx, &sampb.ListParams{SellerId: 1})

				Expect(err).To(BeNil())
				Expect(resp.Status).To(Equal("success"))

				Expect(resp.AccountManager[2].Id).To(Equal(Sam1.ID))
				Expect(resp.AccountManager[2].Email).To(Equal(Sam1.Email))
				Expect(resp.AccountManager[2].Phone).To(Equal(uint64(Sam1.Phone)))
				Expect(resp.AccountManager[2].Name).To(Equal(Sam1.Name))
				Expect(resp.AccountManager[2].Priority).To(Equal(uint64(Sam1.Priority)))
				Expect(resp.AccountManager[2].Role).To(Equal(Sam1.Role))

				Expect(resp.AccountManager[1].Id).To(Equal(Sam3.ID))
				Expect(resp.AccountManager[1].Email).To(Equal(Sam3.Email))
				Expect(resp.AccountManager[1].Phone).To(Equal(uint64(Sam3.Phone)))
				Expect(resp.AccountManager[1].Name).To(Equal(Sam3.Name))
				Expect(resp.AccountManager[1].Priority).To(Equal(uint64(Sam3.Priority)))
				Expect(resp.AccountManager[1].Role).To(Equal(Sam3.Role))

				Expect(resp.AccountManager[0].Id).To(Equal(Sam2.ID))
				Expect(resp.AccountManager[0].Email).To(Equal(Sam2.Email))
				Expect(resp.AccountManager[0].Phone).To(Equal(uint64(Sam2.Phone)))
				Expect(resp.AccountManager[0].Name).To(Equal(Sam2.Name))
				Expect(resp.AccountManager[0].Priority).To(Equal(uint64(Sam2.Priority)))
				Expect(resp.AccountManager[0].Role).To(Equal(Sam2.Role))
			})
		})
	})
})
