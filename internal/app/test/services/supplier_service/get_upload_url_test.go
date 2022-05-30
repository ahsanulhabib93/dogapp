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

var _ = Describe("GetUploadUrl", func() {
	var ctx context.Context

	BeforeEach(func() {
		test_utils.GetContext(&ctx)
	})

	Context("Get Upload url", func() {
		BeforeEach(func() {
			cloudStorageInterface := new(mocks.CloudStorageInterface)
			cloudstorage.InjectGcsMockInstance(cloudStorageInterface)
			cloudStorageInterface.On("GetUploadURL", ctx, utils.GetBucketName(ctx), mock.AnythingOfType("string"), mock.AnythingOfType("time.Time")).Return("https://test/ss2/image.xyz", nil)
		})

		It("Should return path and file url", func() {
			param := &supplierpb.GetUploadUrlParam{UploadType: "SupplierShopImage"}
			res, err := new(services.SupplierService).GetUploadURL(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Fetched upload url successfully"))
			Expect(res.Path).To(HavePrefix("ss2/shop_images/shop_images-"))
			Expect(res.Path).To(HaveSuffix(".jpg"))
			Expect(res.Url).To(Equal("https://test/ss2/image.xyz"))
		})

		It("Should return path and file url for NID front image", func() {
			param := &supplierpb.GetUploadUrlParam{UploadType: "SupplierNIDFrontImage"}
			res, err := new(services.SupplierService).GetUploadURL(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Fetched upload url successfully"))
			Expect(res.Path).To(HavePrefix("ss2/nid_front_images/nid_front_images-"))
			Expect(res.Url).To(Equal("https://test/ss2/image.xyz"))
		})

		It("Should return path and file url for Agreement PDF", func() {
			param := &supplierpb.GetUploadUrlParam{UploadType: "SupplierAgreement"}
			res, err := new(services.SupplierService).GetUploadURL(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Fetched upload url successfully"))
			Expect(res.Path).To(HavePrefix("ss2/agreements/agreements-"))
			Expect(res.Path).To(HaveSuffix(".pdf"))
			Expect(res.Url).To(Equal("https://test/ss2/image.xyz"))
		})

		It("Should return path and file url for Guarantor image", func() {
			param := &supplierpb.GetUploadUrlParam{UploadType: "GuarantorImage"}
			res, err := new(services.SupplierService).GetUploadURL(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Fetched upload url successfully"))
			Expect(res.Path).To(HavePrefix("ss2/guarantor_images/guarantor_images-"))
			Expect(res.Path).To(HaveSuffix(".jpg"))
			Expect(res.Url).To(Equal("https://test/ss2/image.xyz"))
		})
	})

	Context("For invalid upload type", func() {

		It("Should return error", func() {
			param := &supplierpb.GetUploadUrlParam{UploadType: "test"}
			res, err := new(services.SupplierService).GetUploadURL(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Invalid File Type"))
		})
	})

})
