package helpers

import (
	"context"
	"strconv"

	"github.com/shopuptech/event-bus-logs-go/core"
	"github.com/shopuptech/event-bus-logs-go/models/supplier"
	"github.com/shopuptech/event-bus-logs-go/ss2"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/utils"
)

func CreateSupplierLog(ctx context.Context, supplier models.Supplier, metadata map[string]string) (*ss2.SupplierLogKey, *ss2.SupplierLogValue) {
	logContext := &core.Context{
		VaccountId:   int32(utils.GetVaccount(ctx)),
		PortalId:     int32(utils.GetPortalId(ctx)),
		CurrentActId: int32(utils.GetCurrentActId(ctx)),
		XRequestId:   utils.GetXRequestId(ctx),
	}

	//var userId uint64
	//if v := utils.GetCurrentUserID(ctx); v != nil {
	//	userId = *v
	//}

	//metadata := make(map[string]string)
	//metadata["actionName"] = string(action)
	//metadata["userId"] = strconv.FormatUint(userId, 64)
	//metadata["source"] = "ss2"

	event := &core.Event{
		Id:             "",
		Ordering:       0,
		Timestamp:      nil,
		SequenceNumber: 0,
		ReferenceId:    "",
		Metadata:       metadata,
	}

	key := &ss2.SupplierLogKey{
		Event:      event,
		Context:    logContext,
		SupplierId: supplier.ID,
	}

	value := &ss2.SupplierLogValue{
		Id:                        supplier.ID,
		Name:                      supplier.Name,
		Status:                    string(supplier.Status),
		Reason:                    supplier.Reason,
		Email:                     supplier.Email,
		Phone:                     supplier.Phone,
		AlternatePhone:            supplier.AlternatePhone,
		BusinessName:              supplier.BusinessName,
		IsPhoneVerified:           strconv.FormatBool(*supplier.IsPhoneVerified),
		ShopImageUrl:              supplier.ShopImageURL,
		NidNumber:                 supplier.NidNumber,
		NidFrontImageUrl:          supplier.NidFrontImageUrl,
		NidBackImageUrl:           supplier.NidBackImageUrl,
		TradeLicenseUrl:           supplier.TradeLicenseUrl,
		AgreementUrl:              supplier.AgreementUrl,
		ShopOwnerImageUrl:         supplier.ShopOwnerImageUrl,
		GuarantorImageUrl:         supplier.GuarantorImageUrl,
		GuarantorNidNumber:        supplier.GuarantorNidNumber,
		GuarantorNidFrontImageUrl: supplier.GuarantorNidFrontImageUrl,
		GuarantorNidBackImageUrl:  supplier.GuarantorNidBackImageUrl,
		ChequeImageUrl:            supplier.ChequeImageUrl,
		SupplierType:              ss2.SupplierType(supplier.SupplierType),
		SupplierAddresses:         getSupplierAddresses(supplier.SupplierAddresses),
		PaymentAccountDetails:     getPaymentAccountDetails(supplier.PaymentAccountDetails),
		CategoryIds:               getCategoryIds(supplier.SupplierCategoryMappings),
		OpcIds:                    getOpcIds(supplier.SupplierOpcMappings),
		//CreatedAt:                 supplier.CreatedAt,
		//UpdatedAt:                 supplier.UpdatedAt,
	}

	return key, value
}

func getOpcIds(opcMapping []models.SupplierOpcMapping) []uint64 {
	var opcIds []uint64

	for _, o := range opcMapping {
		opcIds = append(opcIds, o.ProcessingCenterID)
	}

	return opcIds
}

func getCategoryIds(categoryMapping []models.SupplierCategoryMapping) []uint64 {
	var categoryIds []uint64

	for _, c := range categoryMapping {
		categoryIds = append(categoryIds, c.CategoryID)
	}

	return categoryIds
}

func getPaymentAccountDetails(paymentAccountDetails []models.PaymentAccountDetail) []*supplier.PaymentAccountDetail {
	var details []*supplier.PaymentAccountDetail

	for _, p := range paymentAccountDetails {
		detail := &supplier.PaymentAccountDetail{
			Id:             p.ID,
			AccountType:    supplier.AccountType(p.AccountType),
			AccountSubType: supplier.AccountSubType(p.AccountSubType),
			AccountName:    p.AccountName,
			AccountNumber:  p.AccountNumber,
			BankId:         p.BankID,
			BranchName:     p.BranchName,
			RoutingNumber:  p.RoutingNumber,
			IsDefault:      p.IsDefault,
			Warehouses:     getWarehouses(p.PaymentAccountDetailWarehouseMappings),
			//CreatedAt:      p.CreatedAt,
			//UpdatedAt:      p.UpdatedAt,
		}

		details = append(details, detail)
	}

	return details
}

func getWarehouses(paymentAccountDetailWarehouseMappings []*models.PaymentAccountDetailWarehouseMapping) []uint64 {
	var warehouses []uint64

	for _, w := range paymentAccountDetailWarehouseMappings {
		warehouses = append(warehouses, w.WarehouseID)
	}

	return warehouses
}

func getSupplierAddresses(supplierAddresses []models.SupplierAddress) []*supplier.SupplierAddress {
	var addresses []*supplier.SupplierAddress

	for _, s := range supplierAddresses {
		address := &supplier.SupplierAddress{
			Id:        s.ID,
			Firstname: s.Firstname,
			Lastname:  s.Lastname,
			Address1:  s.Address1,
			Address2:  s.Address2,
			Landmark:  s.Landmark,
			City:      s.City,
			State:     s.State,
			Country:   s.Country,
			Zipcode:   s.Zipcode,
			Phone:     s.Phone,
			GstNumber: s.GstNumber,
			IsDefault: strconv.FormatBool(s.IsDefault),
			//CreatedAt: s.CreatedAt,
			//UpdatedAt: s.UpdatedAt,
		}

		addresses = append(addresses, address)
	}
	return addresses
}
