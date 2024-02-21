package seller_service_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	spb "github.com/voonik/goConnect/api/go/ss2/seller"
	test_utils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/services"
	"github.com/voonik/ss2/internal/app/test/test_helper"
)

var _ = Describe("Update", func() {
	var ctx context.Context

	BeforeEach(func() {
		test_utils.GetContext(&ctx)
	})

	Context("Failure Cases", func() {
		It("Should return error", func() {
			param := spb.CreateParams{}
			res, err := new(services.SellerService).Create(ctx, &param)
			Expect(res.Status).To(Equal(false))
			Expect(res.Message).To(Equal("Failed to register the seller. Please try again."))
			Expect(err).To(BeNil())
		})
	})
	Context("Success Cases", func() {
		It("Should return success if seller is already registered", func() {
			seller := test_helper.CreateSeller(ctx, &models.Seller{UserID: 101})
			param := spb.CreateParams{Seller: &spb.SellerObject{UserId: 101}}
			res, err := new(services.SellerService).Create(ctx, &param)
			Expect(res.Status).To(Equal(true))
			Expect(res.Message).To(Equal("Seller already registered."))
			Expect(res.UserId).To(Equal(seller.UserID))
			Expect(err).To(BeNil())
		})
	})
})
