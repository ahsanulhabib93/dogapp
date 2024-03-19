// Code generated by mockery v2.33.3. DO NOT EDIT.

package mocks

import (
	context "context"

	attendance "github.com/voonik/goConnect/api/go/sr_service/attendance"

	mock "github.com/stretchr/testify/mock"

	oms_seller "github.com/voonik/goConnect/api/go/oms/oms_seller"

	otp "github.com/voonik/goConnect/api/go/vigeon2/otp"

	payment_gateway "github.com/voonik/goConnect/api/go/paywell_token/payment_gateway"

	product "github.com/voonik/goConnect/api/go/cmt/product"

	users_detail "github.com/voonik/goConnect/api/go/cre_admin/users_detail"
)

// APIHelperInterface is an autogenerated mock type for the APIHelperInterface type
type APIHelperInterface struct {
	mock.Mock
}

// CmtApproveItems provides a mock function with given fields: _a0, _a1
func (_m *APIHelperInterface) CmtApproveItems(_a0 context.Context, _a1 *product.ApproveItemParams) *product.ItemCountResponse {
	ret := _m.Called(_a0, _a1)

	var r0 *product.ItemCountResponse
	if rf, ok := ret.Get(0).(func(context.Context, *product.ApproveItemParams) *product.ItemCountResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*product.ItemCountResponse)
		}
	}

	return r0
}

// CreateOmsSeller provides a mock function with given fields: ctx, param
func (_m *APIHelperInterface) CreateOmsSeller(ctx context.Context, param *oms_seller.SellerParams) *oms_seller.SellerResponse {
	ret := _m.Called(ctx, param)

	var r0 *oms_seller.SellerResponse
	if rf, ok := ret.Get(0).(func(context.Context, *oms_seller.SellerParams) *oms_seller.SellerResponse); ok {
		r0 = rf(ctx, param)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*oms_seller.SellerResponse)
		}
	}

	return r0
}

// CreatePaywellCard provides a mock function with given fields: ctx, params
func (_m *APIHelperInterface) CreatePaywellCard(ctx context.Context, params *payment_gateway.CreateCardRequest) *payment_gateway.CreateCardResponse {
	ret := _m.Called(ctx, params)

	var r0 *payment_gateway.CreateCardResponse
	if rf, ok := ret.Get(0).(func(context.Context, *payment_gateway.CreateCardRequest) *payment_gateway.CreateCardResponse); ok {
		r0 = rf(ctx, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*payment_gateway.CreateCardResponse)
		}
	}

	return r0
}

// FindTalentXUserByPhone provides a mock function with given fields: _a0, _a1
func (_m *APIHelperInterface) FindTalentXUserByPhone(_a0 context.Context, _a1 string) []*attendance.EmployeeRecord {
	ret := _m.Called(_a0, _a1)

	var r0 []*attendance.EmployeeRecord
	if rf, ok := ret.Get(0).(func(context.Context, string) []*attendance.EmployeeRecord); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*attendance.EmployeeRecord)
		}
	}

	return r0
}

// FindUserByPhone provides a mock function with given fields: _a0, _a1
func (_m *APIHelperInterface) FindUserByPhone(_a0 context.Context, _a1 string) *users_detail.UserInfo {
	ret := _m.Called(_a0, _a1)

	var r0 *users_detail.UserInfo
	if rf, ok := ret.Get(0).(func(context.Context, string) *users_detail.UserInfo); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*users_detail.UserInfo)
		}
	}

	return r0
}

// SendOtpAPI provides a mock function with given fields: _a0, _a1
func (_m *APIHelperInterface) SendOtpAPI(_a0 context.Context, _a1 otp.OtpParam) *otp.OtpResponse {
	ret := _m.Called(_a0, _a1)

	var r0 *otp.OtpResponse
	if rf, ok := ret.Get(0).(func(context.Context, otp.OtpParam) *otp.OtpResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*otp.OtpResponse)
		}
	}

	return r0
}

// VerifyOtpAPI provides a mock function with given fields: _a0, _a1
func (_m *APIHelperInterface) VerifyOtpAPI(_a0 context.Context, _a1 otp.VerifyOtpParam) *otp.OtpResponse {
	ret := _m.Called(_a0, _a1)

	var r0 *otp.OtpResponse
	if rf, ok := ret.Get(0).(func(context.Context, otp.VerifyOtpParam) *otp.OtpResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*otp.OtpResponse)
		}
	}

	return r0
}

// NewAPIHelperInterface creates a new instance of APIHelperInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAPIHelperInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *APIHelperInterface {
	mock := &APIHelperInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
