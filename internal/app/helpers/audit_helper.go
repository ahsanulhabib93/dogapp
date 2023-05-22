package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/golang/protobuf/proto"
	"github.com/shopuptech/event-bus-logs-go/core"
	"github.com/shopuptech/event-bus-logs-go/models/supplier"
	"github.com/shopuptech/event-bus-logs-go/ss2"
	supplierPb "github.com/voonik/goConnect/api/go/audit_log_service/supplier"
	"github.com/voonik/goFramework/pkg/misc"
	"github.com/voonik/goFramework/pkg/pubsub/publisher"
	"github.com/voonik/goFramework/pkg/serviceapiconfig"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/utils"
)

type AuditActionType string

const (
	ActionUpdateSupplierStatus      AuditActionType = "update_supplier_status"
	ActionUpdateSupplier            AuditActionType = "update_supplier"
	ActionCreateSupplier            AuditActionType = "create_supplier"
	ActionVerifySupplierPhoneNumber AuditActionType = "verify_supplier_phone_number"
	ActionRemoveSupplierDocuments   AuditActionType = "remove_supplier_document"
)

type AuditHelper struct{}
type Auditor interface {
	RecordAuditAction(ctx context.Context, auditRecord proto.Message) error
}

var auditAction Auditor

func InjectMockAuditActionInstance(mockObj Auditor) {
	auditAction = mockObj
}

func getAuditInstance() Auditor {
	if auditAction == nil {
		return new(AuditHelper)
	}
	return auditAction
}

func AuditAction(ctx context.Context, supplierId uint64, entity string, action AuditActionType, data interface{}) error {
	auditRecord, err := CreateAuditLog(ctx, supplierId, entity, action, data)
	if err != nil {
		return fmt.Errorf("[AuditAction] Failed to create audit log with error: %s", err.Error())
	}

	if err = getAuditInstance().RecordAuditAction(ctx, auditRecord); err != nil {
		return fmt.Errorf("[AuditAction] Failed to publish audit log with error: %s", err.Error())
	}

	return nil
}

func CreateAuditLog(ctx context.Context, supplierId uint64, entity string, action AuditActionType, data interface{}) (*supplierPb.AuditRecord, error) {
	dump, err := json.Marshal(data)

	if err != nil {
		return nil, fmt.Errorf("failed to create dump with error: %s", err.Error())
	}

	var userId uint64
	if v := utils.GetCurrentUserID(ctx); v != nil {
		userId = *v
	}

	auditRecord := &supplierPb.AuditRecord{
		Source:     "ss2",
		Entity:     entity,
		ActionName: string(action),
		UserId:     userId,
		SupplierId: supplierId,
		DataDump:   string(dump),
		VaccountId: uint64(utils.GetVaccount(ctx)),
	}
	return auditRecord, nil
}

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
		//CategoryIds:               supplier.SupplierCategoryMappings,
		//OpcIds:                    supplier.SupplierOpcMappings,
		//CreatedAt:                 supplier.CreatedAt,
		//UpdatedAt:                 supplier.UpdatedAt,
	}

	return key, value
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
			//Warehouses:     p.PaymentAccountDetailWarehouseMappings,
			//CreatedAt:      p.CreatedAt,
			//UpdatedAt:      p.UpdatedAt,
		}

		details = append(details, detail)
	}

	return details
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

func (a *AuditHelper) RecordAuditAction(ctx context.Context, auditRecord proto.Message) error {
	fmt.Print("auditRecord: ", auditRecord)

	transportconf := serviceapiconfig.NewClientOptions(
		serviceapiconfig.WithPubSubTopic(utils.SupplierAuditTopic),
		serviceapiconfig.WithPubSubUrl("/audit/log"),
		serviceapiconfig.WithPubSubKlass("AuditLogService::Supplier::AuditRecord"),
	)
	return publisher.ProduceMessage(ctx, auditRecord, &misc.PubSubMessage{}, utils.SupplierAuditTopic, "", "", transportconf)
}
