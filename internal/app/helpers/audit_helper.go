package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	supplierPb "github.com/voonik/goConnect/api/go/audit_log_service/supplier"

	"github.com/voonik/goFramework/pkg/misc"
	"github.com/voonik/goFramework/pkg/pubsub/publisher"
	"github.com/voonik/goFramework/pkg/serviceapiconfig"
	"github.com/voonik/ss2/internal/app/utils"
)

type AuditActionType string

const (
	ActionUpdateSupplierStatus AuditActionType = "update_supplier_status"
	ActionUpdateSupplier       AuditActionType = "update_supplier"
	ActionCreateSupplier       AuditActionType = "create_supplier"
)

type AuditHelper struct{}
type AuditActionInterface interface {
	RecordAuditAction(ctx context.Context, auditRecord *supplierPb.AuditRecord) error
}

var auditAction AuditActionInterface

func InjectMockAuditActionInstance(mockObj AuditActionInterface) {
	auditAction = mockObj
}

func getAuditInstance() AuditActionInterface {
	if auditAction == nil {
		return new(AuditHelper)
	}
	return auditAction
}

func AuditAction(ctx context.Context, supplierId uint64, entity string, action AuditActionType, data interface{}) error {
	dump, err := json.Marshal(data)
	if err != nil {
		log.Println("AuditAction: Failed to create dump. Error: ", err.Error())
		return err
	}

	auditRecord := &supplierPb.AuditRecord{
		Source:     "ss2",
		Entity:     entity,
		ActionName: string(action),
		UserId:     *utils.GetCurrentUserID(ctx),
		SupplierId: supplierId,
		DataDump:   string(dump),
		VaccountId: uint64(utils.GetVaccount(ctx)),
	}

	if err := getAuditInstance().RecordAuditAction(ctx, auditRecord); err != nil {
		log.Println("AuditAction: Failed to publish audit log. Error: ", err.Error())
	}

	return nil
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
