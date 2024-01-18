package attachment_service

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/test/test_helper"
	"github.com/voonik/ss2/internal/app/utils"
	"time"

	attachmentpb "github.com/voonik/goConnect/api/go/ss2/attachment"
	test_utils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/services"
)

var _ = Describe("Remove attachment", func() {
	var ctx context.Context

	BeforeEach(func() {
		test_utils.GetContext(&ctx)
	})

	Context("When attachment not found", func() {
		It("should raise failed response if wrong attachable id", func() {
			supplier1 := test_helper.CreateSupplier(ctx, &models.Supplier{})
			attachment1 := test_helper.CreateAttachment(ctx, &models.Attachment{
				AttachableType: utils.AttachableTypeSupplier,
				AttachableID:   supplier1.ID,
				FileType:       utils.TIN,
			})
			param := attachmentpb.RemoveAttachmentParams{
				AttachableId:   supplier1.ID + 1,
				AttachableType: 1,
				AttachmentId:   attachment1.ID,
			}
			res, err := services.GetAttachmentServiceInstance().RemoveAttachment(ctx, &param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Attachment not found"))
		})
		It("should raise failed response if wrong attachment id", func() {
			supplier1 := test_helper.CreateSupplier(ctx, &models.Supplier{})
			attachment1 := test_helper.CreateAttachment(ctx, &models.Attachment{
				AttachableType: utils.AttachableTypeSupplier,
				AttachableID:   supplier1.ID,
				FileType:       utils.TIN,
			})
			param := attachmentpb.RemoveAttachmentParams{
				AttachableId:   supplier1.ID,
				AttachableType: 1,
				AttachmentId:   attachment1.ID + 1,
			}
			res, err := services.GetAttachmentServiceInstance().RemoveAttachment(ctx, &param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Attachment not found"))
		})
		It("should raise failed response if wrong attachable type", func() {
			supplier1 := test_helper.CreateSupplier(ctx, &models.Supplier{})
			attachment1 := test_helper.CreateAttachment(ctx, &models.Attachment{
				AttachableType: utils.AttachableTypeSupplier,
				AttachableID:   supplier1.ID,
				FileType:       utils.TIN,
			})
			param := attachmentpb.RemoveAttachmentParams{
				AttachableId:   supplier1.ID,
				AttachableType: 2,
				AttachmentId:   attachment1.ID,
			}
			res, err := services.GetAttachmentServiceInstance().RemoveAttachment(ctx, &param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Attachment not found"))
		})
	})

	Context("When attachment is found", func() {
		It("should delete the attachment", func() {
			supplier1 := test_helper.CreateSupplier(ctx, &models.Supplier{})
			attachment1 := test_helper.CreateAttachment(ctx, &models.Attachment{
				AttachableType: utils.AttachableTypeSupplier,
				AttachableID:   supplier1.ID,
				FileType:       utils.TIN,
			})
			param := attachmentpb.RemoveAttachmentParams{
				AttachableId:   supplier1.ID,
				AttachableType: 1,
				AttachmentId:   attachment1.ID,
			}
			res, err := services.GetAttachmentServiceInstance().RemoveAttachment(ctx, &param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Attachment removed successfully"))
			attachment := &models.Attachment{}
			database.DBAPM(ctx).Model(&models.Attachment{}).Unscoped().First(attachment, attachment1.ID)
			Expect(attachment.ID).To(Equal(attachment1.ID))
			Expect(attachment.DeletedAt.Day()).To(Equal(time.Now().Day()))
		})
	})
})
