package seller_service_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	spb "github.com/voonik/goConnect/api/go/ss2/seller"
	test_utils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/services"
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
})
