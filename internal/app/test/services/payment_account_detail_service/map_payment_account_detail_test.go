package payment_account_detail_service_test

import (
	"context"
	"sort"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	paymentpb "github.com/voonik/goConnect/api/go/ss2/payment_account_detail"
	"github.com/voonik/goFramework/pkg/database"
	test_utils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/services"
	"github.com/voonik/ss2/internal/app/test/test_helper"
	"github.com/voonik/ss2/internal/app/utils"
)

var _ = Describe("MapPaymentAccountDetail", func() {
	var ctx context.Context
	var supplier1 *models.Supplier
	var accountDetail1 *models.PaymentAccountDetail
	var accountDetail2 *models.PaymentAccountDetail
	BeforeEach(func() {
		test_utils.GetContext(&ctx)
		supplier1 = test_helper.CreateSupplier(ctx, &models.Supplier{})
		accountDetail1 = test_helper.CreatePaymentAccountDetail(ctx, &models.PaymentAccountDetail{SupplierID: supplier1.ID, AccountType: utils.Mfs, IsDefault: true})
		bank := test_helper.CreateBank(ctx, &models.Bank{})
		accountDetail2 = test_helper.CreatePaymentAccountDetail(ctx, &models.PaymentAccountDetail{SupplierID: supplier1.ID, AccountType: utils.Bank, BankID: bank.ID})
	})

	Context("PaymentAccountDetailWarehouseMappings", func() {
		BeforeEach(func() {
			test_helper.CreatePaymentAccountDetailWarehouseMappings(ctx, &models.PaymentAccountDetailWarehouseMapping{WarehouseID: 10, PaymentAccountDetailID: accountDetail1.ID})
			test_helper.CreatePaymentAccountDetailWarehouseMappings(ctx, &models.PaymentAccountDetailWarehouseMapping{WarehouseID: 11, PaymentAccountDetailID: accountDetail1.ID})
			test_helper.CreatePaymentAccountDetailWarehouseMappings(ctx, &models.PaymentAccountDetailWarehouseMapping{WarehouseID: 12, PaymentAccountDetailID: accountDetail1.ID})
			test_helper.CreatePaymentAccountDetailWarehouseMappings(ctx, &models.PaymentAccountDetailWarehouseMapping{WarehouseID: 10, PaymentAccountDetailID: accountDetail2.ID})
		})
		Context("With Proper params", func() {
			It("Should Add & Delete Mappiggs according to given warehouse_ids", func() {
				res, err := new(services.PaymentAccountDetailService).MapPaymentAccountDetail(ctx, &paymentpb.MappingParam{Id: accountDetail1.ID, MappableType: "warehouses", MappableIds: []uint64{10, 11, 50}})
				Expect(err).To(BeNil())
				Expect(res.Success).To(BeTrue())
				Expect(res.Message).To(Equal("Mapping Updated Successfully"))

				paymentAccountDetailWarehouseMappings := []*models.PaymentAccountDetailWarehouseMapping{}
				database.DBAPM(ctx).Model(&accountDetail1).Association("PaymentAccountDetailWarehouseMappings").Find(&paymentAccountDetailWarehouseMappings)

				sort.Slice(paymentAccountDetailWarehouseMappings, func(i, j int) bool {
					return paymentAccountDetailWarehouseMappings[i].WarehouseID < paymentAccountDetailWarehouseMappings[j].WarehouseID
				})

				Expect(len(paymentAccountDetailWarehouseMappings)).To(Equal(3))
				Expect(paymentAccountDetailWarehouseMappings[0].WarehouseID).To(Equal(uint64(10)))
				Expect(paymentAccountDetailWarehouseMappings[1].WarehouseID).To(Equal(uint64(11)))
				Expect(paymentAccountDetailWarehouseMappings[2].WarehouseID).To(Equal(uint64(50)))
			})
		})

		Context("With Invalid Params", func() {
			It("Invalid Mapping Type - Should return false response", func() {
				res, err := new(services.PaymentAccountDetailService).MapPaymentAccountDetail(ctx, &paymentpb.MappingParam{Id: accountDetail1.ID, MappableType: "howdy", MappableIds: []uint64{10, 11, 50}})
				Expect(err).To(BeNil())
				Expect(res.Success).To(BeFalse())
				Expect(res.Message).To(Equal("Invalid mapping_type"))
			})

			It("Invalid AccountDetail Id - Should return false response", func() {
				res, err := new(services.PaymentAccountDetailService).MapPaymentAccountDetail(ctx, &paymentpb.MappingParam{Id: accountDetail2.ID + 1, MappableType: "warehouses", MappableIds: []uint64{10, 11, 50}})
				Expect(err).To(BeNil())
				Expect(res.Success).To(BeFalse())
				Expect(res.Message).To(Equal("PaymentAccountDetail Not Found"))
			})
		})

	})
})
