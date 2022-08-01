// Code generated by mockery v2.12.2. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	"github.com/voonik/goFramework/pkg/rest"
	"github.com/voonik/ss2/internal/app/helpers"
)

type ApiCallHelperInterface struct {
	mock.Mock
	Count map[string]int
}

func SetApiCallerMock() *ApiCallHelperInterface {
	mock := &ApiCallHelperInterface{Count: map[string]int{}}
	helpers.InjectMockApiCallHelperInstance(mock)

	return mock
}

func UnsetApiCallerMock() {
	helpers.InjectMockApiCallHelperInstance(nil)
}

func (_m *ApiCallHelperInterface) Get(ctx context.Context, url string, headers map[string]string) (*rest.Response, error) {
	args := _m.Called(ctx, url, headers)
	_m.Count["Get"] += 1
	return args.Get(0).(*rest.Response), args.Error(1)
}
