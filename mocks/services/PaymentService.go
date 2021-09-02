// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	service "go.simple/structure/services"
)

// PaymentService is an autogenerated mock type for the PaymentService type
type PaymentService struct {
	mock.Mock
}

// CreateTransaction provides a mock function with given fields: req
func (_m *PaymentService) CreateTransaction(req service.CreateTransactionRequest) (*service.UserTransactionResponse, error) {
	ret := _m.Called(req)

	var r0 *service.UserTransactionResponse
	if rf, ok := ret.Get(0).(func(service.CreateTransactionRequest) *service.UserTransactionResponse); ok {
		r0 = rf(req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*service.UserTransactionResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(service.CreateTransactionRequest) error); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteUserTransaction provides a mock function with given fields: id, userID
func (_m *PaymentService) DeleteUserTransaction(id int64, userID int64) error {
	ret := _m.Called(id, userID)

	var r0 error
	if rf, ok := ret.Get(0).(func(int64, int64) error); ok {
		r0 = rf(id, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetUserTransactions provides a mock function with given fields: userID, accountID
func (_m *PaymentService) GetUserTransactions(userID int64, accountID int64) ([]service.UserTransactionResponse, error) {
	ret := _m.Called(userID, accountID)

	var r0 []service.UserTransactionResponse
	if rf, ok := ret.Get(0).(func(int64, int64) []service.UserTransactionResponse); ok {
		r0 = rf(userID, accountID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]service.UserTransactionResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64, int64) error); ok {
		r1 = rf(userID, accountID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateUserTransaction provides a mock function with given fields: req
func (_m *PaymentService) UpdateUserTransaction(req service.UpdateTransactionRequest) (*service.UserTransactionResponse, error) {
	ret := _m.Called(req)

	var r0 *service.UserTransactionResponse
	if rf, ok := ret.Get(0).(func(service.UpdateTransactionRequest) *service.UserTransactionResponse); ok {
		r0 = rf(req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*service.UserTransactionResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(service.UpdateTransactionRequest) error); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
