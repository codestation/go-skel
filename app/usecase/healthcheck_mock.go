// Code generated by mockery v2.23.1. DO NOT EDIT.

package usecase

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	model "megpoid.dev/go/go-skel/app/model"
)

// MockHealthcheck is an autogenerated mock type for the Healthcheck type
type MockHealthcheck struct {
	mock.Mock
}

type MockHealthcheck_Expecter struct {
	mock *mock.Mock
}

func (_m *MockHealthcheck) EXPECT() *MockHealthcheck_Expecter {
	return &MockHealthcheck_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: ctx
func (_m *MockHealthcheck) Execute(ctx context.Context) *model.HealthcheckResult {
	ret := _m.Called(ctx)

	var r0 *model.HealthcheckResult
	if rf, ok := ret.Get(0).(func(context.Context) *model.HealthcheckResult); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.HealthcheckResult)
		}
	}

	return r0
}

// MockHealthcheck_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type MockHealthcheck_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockHealthcheck_Expecter) Execute(ctx interface{}) *MockHealthcheck_Execute_Call {
	return &MockHealthcheck_Execute_Call{Call: _e.mock.On("Execute", ctx)}
}

func (_c *MockHealthcheck_Execute_Call) Run(run func(ctx context.Context)) *MockHealthcheck_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockHealthcheck_Execute_Call) Return(_a0 *model.HealthcheckResult) *MockHealthcheck_Execute_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockHealthcheck_Execute_Call) RunAndReturn(run func(context.Context) *model.HealthcheckResult) *MockHealthcheck_Execute_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockHealthcheck interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockHealthcheck creates a new instance of MockHealthcheck. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockHealthcheck(t mockConstructorTestingTNewMockHealthcheck) *MockHealthcheck {
	mock := &MockHealthcheck{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
