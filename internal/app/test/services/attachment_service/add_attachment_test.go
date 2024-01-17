package attachment_service

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/test/test_helper"
	"github.com/voonik/ss2/internal/app/utils"

	attachmentpb "github.com/voonik/goConnect/api/go/ss2/attachment"
	test_utils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/services"
)

var _ = Describe("Add attachment", func() {
	var ctx context.Context

	BeforeEach(func() {
		test_utils.GetContext(&ctx)
	})

	Context("Failure case", func() {
		It("Invalid attachable type", func() {
			param := attachmentpb.AddAttachmentParams{
				AttachableId:    0,
				FileType:        "TIN",
				FileUrl:         "",
				ReferenceNumber: "0",
				AttachableType:  3,
			}
			res, err := services.GetAttachmentServiceInstance().AddAttachment(ctx, &param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Invalid attachable type"))
		})

		It("Invalid file type", func() {
			param := attachmentpb.AddAttachmentParams{
				AttachableId:    0,
				FileType:        "INVALID",
				FileUrl:         "",
				ReferenceNumber: "0",
				AttachableType:  1,
			}
			res, err := services.GetAttachmentServiceInstance().AddAttachment(ctx, &param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Invalid file type"))
		})

		It("Incompatible attachable type and file type", func() {
			param := attachmentpb.AddAttachmentParams{
				AttachableId:    0,
				FileType:        "TIN",
				FileUrl:         "",
				ReferenceNumber: "0",
				AttachableType:  2,
			}
			res, err := services.GetAttachmentServiceInstance().AddAttachment(ctx, &param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Incompatible attachable type and file type"))
		})

		It("Attachable not found", func() {
			param := attachmentpb.AddAttachmentParams{
				AttachableId:    100,
				FileType:        "TIN",
				FileUrl:         "",
				ReferenceNumber: "0",
				AttachableType:  1,
			}
			res, err := services.GetAttachmentServiceInstance().AddAttachment(ctx, &param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Attachable not found"))
		})
	})

	Context("Success case", func() {
		It("should create TIN attachment for supplier", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			param := attachmentpb.AddAttachmentParams{
				AttachableId:    supplier.ID,
				FileType:        "TIN",
				FileUrl:         "high_security/google",
				ReferenceNumber: "1234xyz",
				AttachableType:  1,
			}
			res, err := services.GetAttachmentServiceInstance().AddAttachment(ctx, &param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Attachment added successfully"))
			attachment := &models.Attachment{}
			database.DBAPM(ctx).Model(&models.Attachment{}).Find(&attachment, "attachable_id = ? and reference_number = ?", supplier.ID, "1234xyz")
			Expect(attachment.AttachableType).To(Equal(utils.AttachableTypeSupplier))
			Expect(attachment.AttachableID).To(Equal(supplier.ID))
			Expect(attachment.FileURL).To(Equal("high_security/google"))
			Expect(attachment.ReferenceNumber).To(Equal("1234xyz"))
			Expect(attachment.FileType).To(Equal(utils.TIN))
		})
		It("should create TradeLicense attachment for partner_service_mapping", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			psMapping := test_helper.CreatePartnerServiceMapping(ctx, &models.PartnerServiceMapping{ServiceType: utils.Transporter, ServiceLevel: utils.Captive, SupplierId: supplier.ID})
			param := attachmentpb.AddAttachmentParams{
				AttachableId:    psMapping.ID,
				FileType:        "TradeLicense",
				FileUrl:         "high_security/google",
				ReferenceNumber: "tl_123",
				AttachableType:  2,
			}
			res, err := services.GetAttachmentServiceInstance().AddAttachment(ctx, &param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Attachment added successfully"))
			attachment := &models.Attachment{}
			database.DBAPM(ctx).Model(&models.Attachment{}).Find(&attachment, "attachable_id = ? and reference_number = ?", psMapping.ID, "tl_123")
			Expect(attachment.AttachableType).To(Equal(utils.AttachableTypePartnerServiceMapping))
			Expect(attachment.AttachableID).To(Equal(psMapping.ID))
			Expect(attachment.FileURL).To(Equal("high_security/google"))
			Expect(attachment.ReferenceNumber).To(Equal("tl_123"))
			Expect(attachment.FileType).To(Equal(utils.TradeLicense))
		})
	})
})
