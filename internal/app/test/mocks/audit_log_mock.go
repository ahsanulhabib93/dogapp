package mocks

import (
	"context"

	"github.com/golang/protobuf/proto"
	mock "github.com/stretchr/testify/mock"
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

func (_m *AuditLogMock) RecordAuditAction(ctx context.Context, auditRecord proto.Message) error {
	args := _m.Called(ctx, auditRecord)
	_m.Count["RecordAuditAction"] += 1
	return args.Error(0)
}
