package seller_bank_detail_service_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	sbdpb "github.com/voonik/goConnect/api/go/ss2/seller_bank_detail"
	aaaModels "github.com/voonik/goFramework/pkg/aaa/models"
	test_utils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/helpers"
	"github.com/voonik/ss2/internal/app/services"
	"github.com/voonik/ss2/internal/app/test/mocks"
)

var _ = Describe("Verify Bank Detail", func() {
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

	Context("Success Case", func() {
		It("Should return data", func() {

			param := sbdpb.VerifyBankDetailParams{}

			res, err := new(services.SellerBankDetailService).VerifyBankDetail(ctx, &param)
			Expect(res).To(BeNil())
			Expect(err).To(BeNil())
		})
	})
})
