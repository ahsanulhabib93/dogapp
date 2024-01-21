package attachment_service

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	attachmentpb "github.com/voonik/goConnect/api/go/ss2/attachment"
	aaaModels "github.com/voonik/goFramework/pkg/aaa/models"
	aaaMocks "github.com/voonik/goFramework/pkg/aaa/models/mocks"
	test_utils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/services"
)

var _ = Describe("Get attachment file types", func() {
	var ctx context.Context
	var appPreferenceMockInstance *aaaMocks.AppPreferenceInterface

	BeforeEach(func() {
		test_utils.GetContext(&ctx)
		appPreferenceMockInstance = new(aaaMocks.AppPreferenceInterface)
		aaaModels.InjectMockAppPreferenceServiceInstance(appPreferenceMockInstance)
	})

	Context("For supplier attachment type", func() {
		Context("When app preference is empty", func() {
			BeforeEach(func() {
				appPreferenceMockInstance.On("GetValue", ctx, "active_file_types", []string{}).Return([]string{}, nil) // Ensure mock returns an empty slice
			})
			It("should return empty response", func() {
				param := attachmentpb.GetAttachmentFileTypesParams{
					AttachableType: 1,
				}
				res, err := services.GetAttachmentServiceInstance().GetAttachmentFileTypes(ctx, &param)

				Expect(err).To(BeNil())
				Expect(res.FileTypes).To(Equal([]string{}))
			})
		})

		Context("When app preference has value", func() {
			BeforeEach(func() {
				appPreferenceMockInstance.On("GetValue", ctx, "active_file_types", []string{}).Return([]string{"TIN", "BIN"})
			})
			It("should return active file types", func() {
				param := attachmentpb.GetAttachmentFileTypesParams{
					AttachableType: 1,
				}
				res, err := services.GetAttachmentServiceInstance().GetAttachmentFileTypes(ctx, &param)

				Expect(err).To(BeNil())
				Expect(res.FileTypes).To(Equal([]string{"TIN", "BIN"}))
			})
		})

		Context("When app preference has wrong value", func() {
			BeforeEach(func() {
				appPreferenceMockInstance.On("GetValue", ctx, "active_file_types", []string{}).Return([]string{"TIN", "BIN", "AADHAR"})
			})
			It("should return active file types after ignoring undefined values", func() {
				param := attachmentpb.GetAttachmentFileTypesParams{
					AttachableType: 1,
				}
				res, err := services.GetAttachmentServiceInstance().GetAttachmentFileTypes(ctx, &param)

				Expect(err).To(BeNil())
				Expect(res.FileTypes).To(Equal([]string{"TIN", "BIN"}))
				Expect(res.FileTypes).NotTo(ContainElement("AADHAR"))
			})
		})
	})
})
