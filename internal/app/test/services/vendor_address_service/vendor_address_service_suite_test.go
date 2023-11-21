package vendor_address_service_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/voonik/ss2/internal/app/test"
)

func TestVendorAddressService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "VendorAddressService Suite")
}

var _ = AfterEach(func() {
	test.Cleaner.Clean("sellers")
	test.Cleaner.Clean("vendor_addresses")
})
