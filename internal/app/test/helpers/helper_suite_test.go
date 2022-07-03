package helper_tests

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/voonik/ss2/internal/app/test"
)

func TestHelper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Helper Suite")
}

var _ = AfterEach(func() {
	test.Cleaner.Clean("suppliers")
})
