package seller_pricing_service_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/voonik/ss2/internal/app/test"
)

func TestSellerService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SellerPricingDetailService Suite")
}

var _ = AfterEach(func() {
	test.Cleaner.Clean("sellers")
	test.Cleaner.Clean("seller_pricing_details")
})
