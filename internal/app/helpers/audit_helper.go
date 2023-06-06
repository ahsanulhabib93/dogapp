package helpers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/golang/protobuf/proto"
	supplierPb "github.com/voonik/goConnect/api/go/audit_log_service/supplier"
	"github.com/voonik/goFramework/pkg/misc"
	"github.com/voonik/goFramework/pkg/pubsub/publisher"
	"github.com/voonik/goFramework/pkg/serviceapiconfig"
	"github.com/voonik/ss2/internal/app/appPreference"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/utils"
)

type AuditHelper struct{}
type Auditor interface {
	RecordAuditAction(ctx context.Context, auditRecord proto.Message) error
}

func AuditAction(ctx context.Context, supplierId uint64, entity string, action models.AuditActionType, data interface{}, supplier models.Supplier) error {
	auditRecord, err := CreateAuditLog(ctx, supplierId, entity, action, data)
	if err != nil {
		return fmt.Errorf("[AuditAction] Failed to create audit log with error: %s", err.Error())
	}

	if err = getAuditInstance().RecordAuditAction(ctx, auditRecord); err != nil {
		return fmt.Errorf("[AuditAction] Failed to publish audit log with error: %s", err.Error())
	}

	if appPreference.ShouldSendSupplierLog(ctx) {
		if err = PublishSupplierLog(ctx, action, supplier, data); err != nil {
			return fmt.Errorf("[AuditAction] failed to publish supplier log with err: %s", err.Error())
		}
	}

	return nil
}

func CreateAuditLog(ctx context.Context, supplierId uint64, entity string, action models.AuditActionType, data interface{}) (*supplierPb.AuditRecord, error) {
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

func (a *AuditHelper) RecordAuditAction(ctx context.Context, auditRecord proto.Message) error {
	fmt.Print("auditRecord: ", auditRecord)

	transportconf := serviceapiconfig.NewClientOptions(
		serviceapiconfig.WithPubSubTopic(utils.SupplierAuditTopic),
		serviceapiconfig.WithPubSubUrl("/audit/log"),
		serviceapiconfig.WithPubSubKlass("AuditLogService::Supplier::AuditRecord"),
	)
	return publisher.ProduceMessage(ctx, auditRecord, &misc.PubSubMessage{}, utils.SupplierAuditTopic, "", "", transportconf)
}

func getAuditInstance() Auditor {
	if auditAction == nil {
		return new(AuditHelper)
	}
	return auditAction
}

var auditAction Auditor

func InjectMockAuditActionInstance(mockObj Auditor) {
	auditAction = mockObj
}
