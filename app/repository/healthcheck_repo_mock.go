// Code generated by mockery v2.42.0. DO NOT EDIT.

package repository

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockHealthcheckRepo is an autogenerated mock type for the HealthcheckRepo type
type MockHealthcheckRepo struct {
	mock.Mock
}

type MockHealthcheckRepo_Expecter struct {
	mock *mock.Mock
}

func (_m *MockHealthcheckRepo) EXPECT() *MockHealthcheckRepo_Expecter {
	return &MockHealthcheckRepo_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: ctx
func (_m *MockHealthcheckRepo) Execute(ctx context.Context) error {
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

// MockHealthcheckRepo_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type MockHealthcheckRepo_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockHealthcheckRepo_Expecter) Execute(ctx interface{}) *MockHealthcheckRepo_Execute_Call {
	return &MockHealthcheckRepo_Execute_Call{Call: _e.mock.On("Execute", ctx)}
}

func (_c *MockHealthcheckRepo_Execute_Call) Run(run func(ctx context.Context)) *MockHealthcheckRepo_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockHealthcheckRepo_Execute_Call) Return(_a0 error) *MockHealthcheckRepo_Execute_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockHealthcheckRepo_Execute_Call) RunAndReturn(run func(context.Context) error) *MockHealthcheckRepo_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockHealthcheckRepo creates a new instance of MockHealthcheckRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockHealthcheckRepo(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockHealthcheckRepo {
	mock := &MockHealthcheckRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
