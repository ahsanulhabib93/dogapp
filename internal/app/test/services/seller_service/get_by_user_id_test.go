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
	"github.com/voonik/ss2/internal/app/services"
	"github.com/voonik/ss2/internal/app/test/mocks"
)

var _ = Describe("Get seller by user ID", func() {
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

	Context("When params are given", func() {
		It("Should return seller details", func() {
			param := spb.GetByUserIDParams{}
			res, err := new(services.SellerService).GetByUserID(ctx, &param)
			Expect(res).To(BeNil())
			Expect(err).To(BeNil())
		})
	})
})
