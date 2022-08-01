package mocks

import (
	"context"

	mock "github.com/stretchr/testify/mock"
	supplierPb "github.com/voonik/goConnect/api/go/audit_log_service/supplier"

	"github.com/voonik/ss2/internal/app/helpers"
)

type AuditLogMock struct {
	mock.Mock
	Count map[string]int
}

func SetAuditLogMock() *AuditLogMock {
	mock := &AuditLogMock{Count: map[string]int{}}
	helpers.InjectMockAuditActionInstance(mock)

	return mock
}

func UnsetAuditLogMock() {
	helpers.InjectMockAuditActionInstance(nil)
}

func (_m *AuditLogMock) RecordAuditAction(ctx context.Context, auditRecord *supplierPb.AuditRecord) error {
	args := _m.Called(ctx, auditRecord)
	_m.Count["RecordAuditAction"] += 1
	return args.Error(0)
}
