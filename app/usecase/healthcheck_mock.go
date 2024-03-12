// Code generated by mockery v2.42.0. DO NOT EDIT.

package usecase

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
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
func (_m *MockHealthcheck) Execute(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
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

func (_c *MockHealthcheck_Execute_Call) Return(_a0 error) *MockHealthcheck_Execute_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockHealthcheck_Execute_Call) RunAndReturn(run func(context.Context) error) *MockHealthcheck_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockHealthcheck creates a new instance of MockHealthcheck. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockHealthcheck(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockHealthcheck {
	mock := &MockHealthcheck{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
