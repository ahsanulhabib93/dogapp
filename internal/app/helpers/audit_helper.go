package helpers

import (
	"context"
	"encoding/json"
	"fmt"

	supplierPb "github.com/voonik/goConnect/api/go/audit_log_service/supplier"

	"github.com/voonik/goFramework/pkg/misc"
	"github.com/voonik/goFramework/pkg/pubsub/publisher"
	"github.com/voonik/goFramework/pkg/serviceapiconfig"
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
	RecordAuditAction(ctx context.Context, auditRecord *supplierPb.AuditRecord) error
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

func (a *AuditHelper) RecordAuditAction(ctx context.Context, auditRecord *supplierPb.AuditRecord) error {
	fmt.Print("auditRecord: ", auditRecord)

	transportconf := serviceapiconfig.NewClientOptions(
		serviceapiconfig.WithPubSubTopic(utils.SupplierAuditTopic),
		serviceapiconfig.WithPubSubUrl("/audit/log"),
		serviceapiconfig.WithPubSubKlass("AuditLogService::Supplier::AuditRecord"),
	)
	return publisher.ProduceMessage(ctx, auditRecord, &misc.PubSubMessage{}, utils.SupplierAuditTopic, "", "", transportconf)
}
