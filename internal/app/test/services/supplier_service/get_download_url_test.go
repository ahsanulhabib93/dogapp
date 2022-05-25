package supplier_service_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	supplierpb "github.com/voonik/goConnect/api/go/ss2/supplier"
	"github.com/voonik/goFramework/pkg/cloudstorage"
	"github.com/voonik/goFramework/pkg/cloudstorage/mocks"
	test_utils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/services"
	"github.com/voonik/ss2/internal/app/utils"
)

var _ = Describe("GetDownloadUrl", func() {
	var ctx context.Context

	BeforeEach(func() {
		test_utils.GetContext(&ctx)
	})

	Context("Get Download url", func() {
		var filePath string

		BeforeEach(func() {
			filePath = "ss2/shop_images/shop_images-1651488390922/433/image.png"
			cloudStorageInterface := new(mocks.CloudStorageInterface)
			cloudstorage.InjectGcsMockInstance(cloudStorageInterface)
			cloudStorageInterface.On("GetObjectSignedURLForDownload", ctx, utils.GetBucketName(ctx), filePath, mock.AnythingOfType("time.Time")).Return("https://test/ss2/image.png", nil)
		})

		It("Should return path and file url", func() {
			param := &supplierpb.GetDownloadUrlParam{Path: "ss2/shop_images/shop_images-1651488390922/433/image.png"}
			res, err := new(services.SupplierService).GetDownloadURL(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Fetched url successfully"))
			Expect(res.Path).To(HavePrefix(filePath))
			Expect(res.Url).To(Equal("https://test/ss2/image.png"))
		})
	})

})
