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

var _ = Describe("Validate Field", func() {
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

	Context("When no param is passed", func() {
		It("Should return false as status", func() {
			param := spb.ValidateFieldParams{}
			res, err := new(services.SellerService).ValidateField(ctx, &param)
			Expect(res.Status).To(BeFalse())
			Expect(err).To(BeNil())
		})
	})

	Context("When existing data is passed", func() {
		It("Should return false as status", func() {
			test_helper.CreateSeller(ctx, &models.Seller{
				UserID:    uint64(102),
				BrandName: "test_brand",
			})
			param := spb.ValidateFieldParams{
				Data: map[string]string{"user_id": "102"},
			}
			res, err := new(services.SellerService).ValidateField(ctx, &param)
			Expect(res.Status).To(BeFalse())
			Expect(err).To(BeNil())
		})
	})

	Context("When non existing data is passed", func() {
		It("Should return true as status", func() {
			param := spb.ValidateFieldParams{
				Data: map[string]string{"user_id": "102"},
			}
			res, err := new(services.SellerService).ValidateField(ctx, &param)
			Expect(res.Status).To(BeTrue())
			Expect(err).To(BeNil())
		})
	})

	Context("When non existing field is passed", func() {
		It("Should return true as status", func() {
			param := spb.ValidateFieldParams{
				Data: map[string]string{"new_field": "102"},
			}
			res, err := new(services.SellerService).ValidateField(ctx, &param)
			Expect(res.Status).To(BeTrue())
			Expect(err).To(BeNil())
		})
	})
})
