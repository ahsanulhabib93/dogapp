package seller_service_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	spb "github.com/voonik/goConnect/api/go/ss2/seller"
	aaaModels "github.com/voonik/goFramework/pkg/aaa/models"
	test_utils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/helpers"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/services"
	"github.com/voonik/ss2/internal/app/test/mocks"
	"github.com/voonik/ss2/internal/app/test/test_helper"
)

var _ = Describe("Seller Phone Relation", func() {
	var ctx context.Context
	var mockAudit *mocks.AuditLogMock

	BeforeEach(func() {
		test_utils.GetContext(&ctx)
		mocks.UnsetOpcMock()

		mockAudit = mocks.SetAuditLogMock()
		mockAudit.On("RecordAuditAction", ctx, mock.Anything).Return(nil)
	})

	AfterEach(func() {
		mocks.UnsetAuditLogMock()
		helpers.InjectMockAPIHelperInstance(nil)
		helpers.InjectMockIdentityUserApiHelperInstance(nil)
		aaaModels.InjectMockAppPreferenceServiceInstance(nil)
	})

	Context("When primary phone is not passed as param", func() {
		It("Should return error", func() {
			param := spb.SellerPhoneRelationParams{}
			res, err := new(services.SellerService).SellerPhoneRelation(ctx, &param)
			Expect(len(res.Seller)).To(Equal(0))
			Expect(res.Status).To(Equal("failure"))
			Expect(res.Message).To(Equal("no valid param"))
			Expect(err).To(BeNil())
		})
	})

	Context("When invalid primary phone is passed as param", func() {
		It("Should return error", func() {
			param := spb.SellerPhoneRelationParams{
				Phone: []string{"123"},
			}
			res, err := new(services.SellerService).SellerPhoneRelation(ctx, &param)
			Expect(len(res.Seller)).To(Equal(0))
			Expect(res.Status).To(Equal("failure"))
			Expect(res.Message).To(Equal("seller not found"))
			Expect(err).To(BeNil())
		})
	})

	Context("When valid primary phone is passed as param", func() {
		It("Should return seller data", func() {
			seller1 := test_helper.CreateSeller(ctx, &models.Seller{})
			param := spb.SellerPhoneRelationParams{
				Phone: []string{seller1.PrimaryPhone},
			}
			res, err := new(services.SellerService).SellerPhoneRelation(ctx, &param)
			Expect(len(res.Seller)).To(Equal(1))
			Expect(res.Seller[0].UserId).To(Equal(seller1.UserID))
			Expect(res.Seller[0].PrimaryPhone).To(Equal(seller1.PrimaryPhone))
			Expect(res.Status).To(Equal("success"))
			Expect(res.Message).To(Equal("fetched seller details successfully"))
			Expect(err).To(BeNil())
		})
	})

	Context("When multiple primary phone numbers are passed as param", func() {
		It("Should return seller data", func() {
			seller1 := test_helper.CreateSeller(ctx, &models.Seller{})
			seller2 := test_helper.CreateSeller(ctx, &models.Seller{})
			param := spb.SellerPhoneRelationParams{
				Phone: []string{seller1.PrimaryPhone, seller2.PrimaryPhone},
			}
			res, err := new(services.SellerService).SellerPhoneRelation(ctx, &param)
			Expect(len(res.Seller)).To(Equal(2))
			Expect(res.Seller[0].UserId).To(Equal(seller1.UserID))
			Expect(res.Seller[0].PrimaryPhone).To(Equal(seller1.PrimaryPhone))
			Expect(res.Seller[1].UserId).To(Equal(seller2.UserID))
			Expect(res.Seller[1].PrimaryPhone).To(Equal(seller2.PrimaryPhone))
			Expect(res.Status).To(Equal("success"))
			Expect(res.Message).To(Equal("fetched seller details successfully"))
			Expect(err).To(BeNil())
		})

		It("Should return seller data for valid phone numbers", func() {
			seller1 := test_helper.CreateSeller(ctx, &models.Seller{})
			param := spb.SellerPhoneRelationParams{
				Phone: []string{seller1.PrimaryPhone},
			}

			res, err := new(services.SellerService).SellerPhoneRelation(ctx, &param)
			Expect(len(res.Seller)).To(Equal(1))
			Expect(res.Seller[0].UserId).To(Equal(seller1.UserID))
			Expect(res.Seller[0].PrimaryPhone).To(Equal(seller1.PrimaryPhone))
			Expect(res.Status).To(Equal("success"))
			Expect(res.Message).To(Equal("fetched seller details successfully"))
			Expect(err).To(BeNil())
		})
	})
})
