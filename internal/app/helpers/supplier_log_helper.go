package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/shopuptech/event-bus-logs-go/core"
	"github.com/shopuptech/event-bus-logs-go/models/supplier"
	"github.com/shopuptech/event-bus-logs-go/ss2"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/publisher"
	"github.com/voonik/ss2/internal/app/utils"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const supplierLogTopic = "supplier-service-log"

func PublishSupplierLog(ctx context.Context, action models.AuditActionType, supplier models.Supplier, data interface{}) error {
	key, value := supplierLog(ctx, supplier, supplierMetadata(ctx, action, data))

	_, err := publisher.Publish(ctx, supplierLogTopic, key, value)
	if err != nil {
		return fmt.Errorf("failed to publish supplier log with err: %s", err.Error())
	}

	return nil
}

func supplierMetadata(ctx context.Context, action models.AuditActionType, data interface{}) map[string]string {
	m := make(map[string]string)

	var userId uint64
	if v := utils.GetCurrentUserID(ctx); v != nil {
		userId = *v
	}

	m["source"] = "ss2"
	m["user_id"] = strconv.FormatUint(userId, 10)
	m["action_name"] = string(action)

	var dataMap map[string]string
	d, _ := json.Marshal(data)
	json.Unmarshal(d, &dataMap)

	for k, v := range dataMap {
		m[k] = v
	}

	return m
}

func supplierLog(ctx context.Context, supplier models.Supplier, metadata map[string]string) (*ss2.SupplierLogKey, *ss2.SupplierLogValue) {
	logContext := &core.Context{
		VaccountId:   int32(utils.GetVaccount(ctx)),
		PortalId:     int32(utils.GetPortalId(ctx)),
		CurrentActId: int32(utils.GetCurrentActId(ctx)),
		XRequestId:   utils.GetXRequestId(ctx),
	}

	event := &core.Event{
		Id:       "",
		Metadata: metadata,
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
		CreatedAt:                 timestamppb.New(supplier.CreatedAt),
		UpdatedAt:                 timestamppb.New(supplier.UpdatedAt),
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
			CreatedAt:      timestamppb.New(p.CreatedAt),
			UpdatedAt:      timestamppb.New(p.UpdatedAt),
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
			CreatedAt: timestamppb.New(s.CreatedAt),
			UpdatedAt: timestamppb.New(s.UpdatedAt),
		}

		addresses = append(addresses, address)
	}
	return addresses
}
