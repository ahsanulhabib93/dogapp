package attachment_service

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/voonik/ss2/internal/app/test"
)

func TestAttachmentService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Attachment Suite")
}

var _ = AfterEach(func() {
	test.Cleaner.Clean("suppliers", "partner_service_mappings", "attachments")
})
