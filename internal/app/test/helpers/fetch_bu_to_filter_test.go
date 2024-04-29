package helper_tests

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	userMappingPb "github.com/voonik/goConnect/api/go/oms/user_mapping"
	"github.com/voonik/goFramework/pkg/misc"
	testUtils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/helpers"
	"github.com/voonik/ss2/internal/app/helpers/mocks"
)

var _ = Describe("FetchBuToFilter", func() {
	var ctx context.Context
	var omsApiMock *mocks.OmsApiHelperInterface
	var omsResponse *userMappingPb.UserMappingResponse
	BeforeEach(func() {
		testUtils.GetContext(&ctx)
		omsApiMock = new(mocks.OmsApiHelperInterface)
		helpers.InjectMockOmsAPIHelperInstance(omsApiMock)
	})
	AfterEach(func() {
		helpers.InjectMockOmsAPIHelperInstance(nil)
	})
	Context("When current user is not present", func() {
		It("Should return input BUs", func() {
			bu, err := helpers.FetchBuToFilter(ctx, []uint64{1})
			Expect(err).To(BeNil())
			Expect(bu).To(Equal([]uint64{1}))
		})
	})
	Context("When current user is present", func() {
		BeforeEach(func() {
			ctx = misc.SetInContextThreadObject(ctx, &misc.ThreadObject{VaccountId: 1, PortalId: 1, UserData: &misc.UserData{
				UserId: 999,
				Name:   "Test User",
				Phone:  "8801485743298",
			}})
		})
		Context("When inputBUs is empty", func() {
			BeforeEach(func() {
				omsResponse = &userMappingPb.UserMappingResponse{
					Data: []*userMappingPb.UserMappingData{
						{
							UserId:        999,
							BusinessUnits: []uint64{1},
							OpcIds:        []uint64{1},
							ZoneIds:       []uint64{1},
						},
					},
				}
			})
			It("Should return BUs from oms response", func() {
				omsApiMock.On("FetchUserMappingData", ctx, []uint64{999}).Return(omsResponse, nil)
				bu, err := helpers.FetchBuToFilter(ctx, []uint64{})
				Expect(err).To(BeNil())
				Expect(bu).To(Equal(omsResponse.Data[0].BusinessUnits))
			})
		})
		Context("When inputBUs is present", func() {
			Context("When OMS BU is empty", func() {
				It("Should return BUs from passed as arguments", func() {
					omsResponse.Data[0].BusinessUnits = []uint64{}
					omsApiMock.On("FetchUserMappingData", ctx, []uint64{999}).Return(omsResponse, nil)
					bu, err := helpers.FetchBuToFilter(ctx, []uint64{1, 2, 3})
					omsResponse.Data[0].BusinessUnits = []uint64{1}
					Expect(err).To(BeNil())
					Expect(bu).To(Equal([]uint64{1, 2, 3}))
				})
			})
			Context("When OMS response is not present for the user", func() {
				It("Should return BUs from passed as arguments", func() {
					omsResponse.Data[0].UserId = 1000
					omsApiMock.On("FetchUserMappingData", ctx, []uint64{999}).Return(omsResponse, nil)
					bu, err := helpers.FetchBuToFilter(ctx, []uint64{1, 2, 3})
					omsResponse.Data[0].UserId = 999
					omsResponse.Data[0].BusinessUnits = []uint64{1}
					Expect(err).To(BeNil())
					Expect(bu).To(Equal([]uint64{1, 2, 3}))
				})
			})
			Context("When OMS response has user data", func() {
				It("Should return intersection of input and oms response", func() {
					omsApiMock.On("FetchUserMappingData", ctx, []uint64{999}).Return(omsResponse, nil)
					bu, err := helpers.FetchBuToFilter(ctx, []uint64{1, 2, 3})
					Expect(err).To(BeNil())
					Expect(bu).To(Equal([]uint64{1}))
				})
			})
		})
	})
})
