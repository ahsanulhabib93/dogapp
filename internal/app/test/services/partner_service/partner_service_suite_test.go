package partner_service_service_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/voonik/ss2/internal/app/test"
)

func TestPartnerService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PartnerService Suite")
}

var _ = AfterEach(func() {
	test.Cleaner.Clean("suppliers", "supplier_addresses", "supplier_category_mappings", "partner_service_mappings")
})
