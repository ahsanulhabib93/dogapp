package key_account_manager_service_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/voonik/ss2/internal/app/test"
)

func TestKeyAccountManagerService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "KeyAccountManagerService Suite")
}

var _ = AfterEach(func() {
	test.Cleaner.Clean("suppliers", "key_account_managers")
})
