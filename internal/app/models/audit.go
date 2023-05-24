package models

type AuditActionType string

const (
	ActionUpdateSupplierStatus      AuditActionType = "update_supplier_status"
	ActionUpdateSupplier            AuditActionType = "update_supplier"
	ActionCreateSupplier            AuditActionType = "create_supplier"
	ActionVerifySupplierPhoneNumber AuditActionType = "verify_supplier_phone_number"
	ActionRemoveSupplierDocuments   AuditActionType = "remove_supplier_document"
)
