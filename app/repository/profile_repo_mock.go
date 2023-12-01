// Code generated by mockery v2.23.1. DO NOT EDIT.

package repository

import (
	context "context"

	clause "go.megpoid.dev/go-skel/pkg/clause"

	exp "github.com/doug-martin/goqu/v9/exp"

	mock "github.com/stretchr/testify/mock"

	model "go.megpoid.dev/go-skel/app/model"

	response "go.megpoid.dev/go-skel/pkg/response"
)

// MockProfileRepo is an autogenerated mock type for the ProfileRepo type
type MockProfileRepo struct {
	mock.Mock
}

type MockProfileRepo_Expecter struct {
	mock *mock.Mock
}

func (_m *MockProfileRepo) EXPECT() *MockProfileRepo_Expecter {
	return &MockProfileRepo_Expecter{mock: &_m.Mock}
}

// Delete provides a mock function with given fields: ctx, id
func (_m *MockProfileRepo) Delete(ctx context.Context, id int64) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockProfileRepo_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type MockProfileRepo_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - ctx context.Context
//   - id int64
func (_e *MockProfileRepo_Expecter) Delete(ctx interface{}, id interface{}) *MockProfileRepo_Delete_Call {
	return &MockProfileRepo_Delete_Call{Call: _e.mock.On("Delete", ctx, id)}
}

func (_c *MockProfileRepo_Delete_Call) Run(run func(ctx context.Context, id int64)) *MockProfileRepo_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64))
	})
	return _c
}

func (_c *MockProfileRepo_Delete_Call) Return(_a0 error) *MockProfileRepo_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockProfileRepo_Delete_Call) RunAndReturn(run func(context.Context, int64) error) *MockProfileRepo_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteBy provides a mock function with given fields: ctx, expr
func (_m *MockProfileRepo) DeleteBy(ctx context.Context, expr exp.Ex) (int64, error) {
	ret := _m.Called(ctx, expr)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, exp.Ex) (int64, error)); ok {
		return rf(ctx, expr)
	}
	if rf, ok := ret.Get(0).(func(context.Context, exp.Ex) int64); ok {
		r0 = rf(ctx, expr)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, exp.Ex) error); ok {
		r1 = rf(ctx, expr)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockProfileRepo_DeleteBy_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteBy'
type MockProfileRepo_DeleteBy_Call struct {
	*mock.Call
}

// DeleteBy is a helper method to define mock.On call
//   - ctx context.Context
//   - expr exp.Ex
func (_e *MockProfileRepo_Expecter) DeleteBy(ctx interface{}, expr interface{}) *MockProfileRepo_DeleteBy_Call {
	return &MockProfileRepo_DeleteBy_Call{Call: _e.mock.On("DeleteBy", ctx, expr)}
}

func (_c *MockProfileRepo_DeleteBy_Call) Run(run func(ctx context.Context, expr exp.Ex)) *MockProfileRepo_DeleteBy_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(exp.Ex))
	})
	return _c
}

func (_c *MockProfileRepo_DeleteBy_Call) Return(_a0 int64, _a1 error) *MockProfileRepo_DeleteBy_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockProfileRepo_DeleteBy_Call) RunAndReturn(run func(context.Context, exp.Ex) (int64, error)) *MockProfileRepo_DeleteBy_Call {
	_c.Call.Return(run)
	return _c
}

// Exists provides a mock function with given fields: ctx, expr
func (_m *MockProfileRepo) Exists(ctx context.Context, expr exp.Expression) (bool, error) {
	ret := _m.Called(ctx, expr)

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, exp.Expression) (bool, error)); ok {
		return rf(ctx, expr)
	}
	if rf, ok := ret.Get(0).(func(context.Context, exp.Expression) bool); ok {
		r0 = rf(ctx, expr)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, exp.Expression) error); ok {
		r1 = rf(ctx, expr)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockProfileRepo_Exists_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Exists'
type MockProfileRepo_Exists_Call struct {
	*mock.Call
}

// Exists is a helper method to define mock.On call
//   - ctx context.Context
//   - expr exp.Expression
func (_e *MockProfileRepo_Expecter) Exists(ctx interface{}, expr interface{}) *MockProfileRepo_Exists_Call {
	return &MockProfileRepo_Exists_Call{Call: _e.mock.On("Exists", ctx, expr)}
}

func (_c *MockProfileRepo_Exists_Call) Run(run func(ctx context.Context, expr exp.Expression)) *MockProfileRepo_Exists_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(exp.Expression))
	})
	return _c
}

func (_c *MockProfileRepo_Exists_Call) Return(_a0 bool, _a1 error) *MockProfileRepo_Exists_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockProfileRepo_Exists_Call) RunAndReturn(run func(context.Context, exp.Expression) (bool, error)) *MockProfileRepo_Exists_Call {
	_c.Call.Return(run)
	return _c
}

// Find provides a mock function with given fields: ctx, dest, id
func (_m *MockProfileRepo) Find(ctx context.Context, dest *model.Profile, id int64) error {
	ret := _m.Called(ctx, dest, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.Profile, int64) error); ok {
		r0 = rf(ctx, dest, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockProfileRepo_Find_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Find'
type MockProfileRepo_Find_Call struct {
	*mock.Call
}

// Find is a helper method to define mock.On call
//   - ctx context.Context
//   - dest *model.Profile
//   - id int64
func (_e *MockProfileRepo_Expecter) Find(ctx interface{}, dest interface{}, id interface{}) *MockProfileRepo_Find_Call {
	return &MockProfileRepo_Find_Call{Call: _e.mock.On("Find", ctx, dest, id)}
}

func (_c *MockProfileRepo_Find_Call) Run(run func(ctx context.Context, dest *model.Profile, id int64)) *MockProfileRepo_Find_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*model.Profile), args[2].(int64))
	})
	return _c
}

func (_c *MockProfileRepo_Find_Call) Return(_a0 error) *MockProfileRepo_Find_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockProfileRepo_Find_Call) RunAndReturn(run func(context.Context, *model.Profile, int64) error) *MockProfileRepo_Find_Call {
	_c.Call.Return(run)
	return _c
}

// First provides a mock function with given fields: ctx, expr, order
func (_m *MockProfileRepo) First(ctx context.Context, expr exp.Expression, order ...exp.OrderedExpression) (*model.Profile, error) {
	_va := make([]interface{}, len(order))
	for _i := range order {
		_va[_i] = order[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, expr)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *model.Profile
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, exp.Expression, ...exp.OrderedExpression) (*model.Profile, error)); ok {
		return rf(ctx, expr, order...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, exp.Expression, ...exp.OrderedExpression) *model.Profile); ok {
		r0 = rf(ctx, expr, order...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Profile)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, exp.Expression, ...exp.OrderedExpression) error); ok {
		r1 = rf(ctx, expr, order...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockProfileRepo_First_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'First'
type MockProfileRepo_First_Call struct {
	*mock.Call
}

// First is a helper method to define mock.On call
//   - ctx context.Context
//   - expr exp.Expression
//   - order ...exp.OrderedExpression
func (_e *MockProfileRepo_Expecter) First(ctx interface{}, expr interface{}, order ...interface{}) *MockProfileRepo_First_Call {
	return &MockProfileRepo_First_Call{Call: _e.mock.On("First",
		append([]interface{}{ctx, expr}, order...)...)}
}

func (_c *MockProfileRepo_First_Call) Run(run func(ctx context.Context, expr exp.Expression, order ...exp.OrderedExpression)) *MockProfileRepo_First_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]exp.OrderedExpression, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(exp.OrderedExpression)
			}
		}
		run(args[0].(context.Context), args[1].(exp.Expression), variadicArgs...)
	})
	return _c
}

func (_c *MockProfileRepo_First_Call) Return(_a0 *model.Profile, _a1 error) *MockProfileRepo_First_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockProfileRepo_First_Call) RunAndReturn(run func(context.Context, exp.Expression, ...exp.OrderedExpression) (*model.Profile, error)) *MockProfileRepo_First_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: ctx, id
func (_m *MockProfileRepo) Get(ctx context.Context, id int64) (*model.Profile, error) {
	ret := _m.Called(ctx, id)

	var r0 *model.Profile
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (*model.Profile, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) *model.Profile); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Profile)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockProfileRepo_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type MockProfileRepo_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - id int64
func (_e *MockProfileRepo_Expecter) Get(ctx interface{}, id interface{}) *MockProfileRepo_Get_Call {
	return &MockProfileRepo_Get_Call{Call: _e.mock.On("Get", ctx, id)}
}

func (_c *MockProfileRepo_Get_Call) Run(run func(ctx context.Context, id int64)) *MockProfileRepo_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64))
	})
	return _c
}

func (_c *MockProfileRepo_Get_Call) Return(_a0 *model.Profile, _a1 error) *MockProfileRepo_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockProfileRepo_Get_Call) RunAndReturn(run func(context.Context, int64) (*model.Profile, error)) *MockProfileRepo_Get_Call {
	_c.Call.Return(run)
	return _c
}

// GetBy provides a mock function with given fields: ctx, expr
func (_m *MockProfileRepo) GetBy(ctx context.Context, expr exp.Expression) (*model.Profile, error) {
	ret := _m.Called(ctx, expr)

	var r0 *model.Profile
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, exp.Expression) (*model.Profile, error)); ok {
		return rf(ctx, expr)
	}
	if rf, ok := ret.Get(0).(func(context.Context, exp.Expression) *model.Profile); ok {
		r0 = rf(ctx, expr)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Profile)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, exp.Expression) error); ok {
		r1 = rf(ctx, expr)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockProfileRepo_GetBy_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetBy'
type MockProfileRepo_GetBy_Call struct {
	*mock.Call
}

// GetBy is a helper method to define mock.On call
//   - ctx context.Context
//   - expr exp.Expression
func (_e *MockProfileRepo_Expecter) GetBy(ctx interface{}, expr interface{}) *MockProfileRepo_GetBy_Call {
	return &MockProfileRepo_GetBy_Call{Call: _e.mock.On("GetBy", ctx, expr)}
}

func (_c *MockProfileRepo_GetBy_Call) Run(run func(ctx context.Context, expr exp.Expression)) *MockProfileRepo_GetBy_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(exp.Expression))
	})
	return _c
}

func (_c *MockProfileRepo_GetBy_Call) Return(_a0 *model.Profile, _a1 error) *MockProfileRepo_GetBy_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockProfileRepo_GetBy_Call) RunAndReturn(run func(context.Context, exp.Expression) (*model.Profile, error)) *MockProfileRepo_GetBy_Call {
	_c.Call.Return(run)
	return _c
}

// GetByEmail provides a mock function with given fields: ctx, email
func (_m *MockProfileRepo) GetByEmail(ctx context.Context, email string) (*model.Profile, error) {
	ret := _m.Called(ctx, email)

	var r0 *model.Profile
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*model.Profile, error)); ok {
		return rf(ctx, email)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.Profile); ok {
		r0 = rf(ctx, email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Profile)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockProfileRepo_GetByEmail_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByEmail'
type MockProfileRepo_GetByEmail_Call struct {
	*mock.Call
}

// GetByEmail is a helper method to define mock.On call
//   - ctx context.Context
//   - email string
func (_e *MockProfileRepo_Expecter) GetByEmail(ctx interface{}, email interface{}) *MockProfileRepo_GetByEmail_Call {
	return &MockProfileRepo_GetByEmail_Call{Call: _e.mock.On("GetByEmail", ctx, email)}
}

func (_c *MockProfileRepo_GetByEmail_Call) Run(run func(ctx context.Context, email string)) *MockProfileRepo_GetByEmail_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockProfileRepo_GetByEmail_Call) Return(_a0 *model.Profile, _a1 error) *MockProfileRepo_GetByEmail_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockProfileRepo_GetByEmail_Call) RunAndReturn(run func(context.Context, string) (*model.Profile, error)) *MockProfileRepo_GetByEmail_Call {
	_c.Call.Return(run)
	return _c
}

// GetForUpdate provides a mock function with given fields: ctx, expr, order
func (_m *MockProfileRepo) GetForUpdate(ctx context.Context, expr exp.Expression, order ...exp.OrderedExpression) (*model.Profile, error) {
	_va := make([]interface{}, len(order))
	for _i := range order {
		_va[_i] = order[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, expr)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *model.Profile
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, exp.Expression, ...exp.OrderedExpression) (*model.Profile, error)); ok {
		return rf(ctx, expr, order...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, exp.Expression, ...exp.OrderedExpression) *model.Profile); ok {
		r0 = rf(ctx, expr, order...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Profile)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, exp.Expression, ...exp.OrderedExpression) error); ok {
		r1 = rf(ctx, expr, order...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockProfileRepo_GetForUpdate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetForUpdate'
type MockProfileRepo_GetForUpdate_Call struct {
	*mock.Call
}

// GetForUpdate is a helper method to define mock.On call
//   - ctx context.Context
//   - expr exp.Expression
//   - order ...exp.OrderedExpression
func (_e *MockProfileRepo_Expecter) GetForUpdate(ctx interface{}, expr interface{}, order ...interface{}) *MockProfileRepo_GetForUpdate_Call {
	return &MockProfileRepo_GetForUpdate_Call{Call: _e.mock.On("GetForUpdate",
		append([]interface{}{ctx, expr}, order...)...)}
}

func (_c *MockProfileRepo_GetForUpdate_Call) Run(run func(ctx context.Context, expr exp.Expression, order ...exp.OrderedExpression)) *MockProfileRepo_GetForUpdate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]exp.OrderedExpression, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(exp.OrderedExpression)
			}
		}
		run(args[0].(context.Context), args[1].(exp.Expression), variadicArgs...)
	})
	return _c
}

func (_c *MockProfileRepo_GetForUpdate_Call) Return(_a0 *model.Profile, _a1 error) *MockProfileRepo_GetForUpdate_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockProfileRepo_GetForUpdate_Call) RunAndReturn(run func(context.Context, exp.Expression, ...exp.OrderedExpression) (*model.Profile, error)) *MockProfileRepo_GetForUpdate_Call {
	_c.Call.Return(run)
	return _c
}

// Insert provides a mock function with given fields: ctx, req
func (_m *MockProfileRepo) Insert(ctx context.Context, req *model.Profile) error {
	ret := _m.Called(ctx, req)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.Profile) error); ok {
		r0 = rf(ctx, req)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockProfileRepo_Insert_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Insert'
type MockProfileRepo_Insert_Call struct {
	*mock.Call
}

// Insert is a helper method to define mock.On call
//   - ctx context.Context
//   - req *model.Profile
func (_e *MockProfileRepo_Expecter) Insert(ctx interface{}, req interface{}) *MockProfileRepo_Insert_Call {
	return &MockProfileRepo_Insert_Call{Call: _e.mock.On("Insert", ctx, req)}
}

func (_c *MockProfileRepo_Insert_Call) Run(run func(ctx context.Context, req *model.Profile)) *MockProfileRepo_Insert_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*model.Profile))
	})
	return _c
}

func (_c *MockProfileRepo_Insert_Call) Return(_a0 error) *MockProfileRepo_Insert_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockProfileRepo_Insert_Call) RunAndReturn(run func(context.Context, *model.Profile) error) *MockProfileRepo_Insert_Call {
	_c.Call.Return(run)
	return _c
}

// List provides a mock function with given fields: ctx, opts
func (_m *MockProfileRepo) List(ctx context.Context, opts ...clause.FilterOption) (*response.ListResponse[*model.Profile], error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *response.ListResponse[*model.Profile]
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, ...clause.FilterOption) (*response.ListResponse[*model.Profile], error)); ok {
		return rf(ctx, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ...clause.FilterOption) *response.ListResponse[*model.Profile]); ok {
		r0 = rf(ctx, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*response.ListResponse[*model.Profile])
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, ...clause.FilterOption) error); ok {
		r1 = rf(ctx, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockProfileRepo_List_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'List'
type MockProfileRepo_List_Call struct {
	*mock.Call
}

// List is a helper method to define mock.On call
//   - ctx context.Context
//   - opts ...clause.FilterOption
func (_e *MockProfileRepo_Expecter) List(ctx interface{}, opts ...interface{}) *MockProfileRepo_List_Call {
	return &MockProfileRepo_List_Call{Call: _e.mock.On("List",
		append([]interface{}{ctx}, opts...)...)}
}

func (_c *MockProfileRepo_List_Call) Run(run func(ctx context.Context, opts ...clause.FilterOption)) *MockProfileRepo_List_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]clause.FilterOption, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(clause.FilterOption)
			}
		}
		run(args[0].(context.Context), variadicArgs...)
	})
	return _c
}

func (_c *MockProfileRepo_List_Call) Return(_a0 *response.ListResponse[*model.Profile], _a1 error) *MockProfileRepo_List_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockProfileRepo_List_Call) RunAndReturn(run func(context.Context, ...clause.FilterOption) (*response.ListResponse[*model.Profile], error)) *MockProfileRepo_List_Call {
	_c.Call.Return(run)
	return _c
}

// ListBy provides a mock function with given fields: ctx, expr, opts
func (_m *MockProfileRepo) ListBy(ctx context.Context, expr exp.Expression, opts ...clause.FilterOption) (*response.ListResponse[*model.Profile], error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, expr)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *response.ListResponse[*model.Profile]
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, exp.Expression, ...clause.FilterOption) (*response.ListResponse[*model.Profile], error)); ok {
		return rf(ctx, expr, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, exp.Expression, ...clause.FilterOption) *response.ListResponse[*model.Profile]); ok {
		r0 = rf(ctx, expr, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*response.ListResponse[*model.Profile])
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, exp.Expression, ...clause.FilterOption) error); ok {
		r1 = rf(ctx, expr, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockProfileRepo_ListBy_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListBy'
type MockProfileRepo_ListBy_Call struct {
	*mock.Call
}

// ListBy is a helper method to define mock.On call
//   - ctx context.Context
//   - expr exp.Expression
//   - opts ...clause.FilterOption
func (_e *MockProfileRepo_Expecter) ListBy(ctx interface{}, expr interface{}, opts ...interface{}) *MockProfileRepo_ListBy_Call {
	return &MockProfileRepo_ListBy_Call{Call: _e.mock.On("ListBy",
		append([]interface{}{ctx, expr}, opts...)...)}
}

func (_c *MockProfileRepo_ListBy_Call) Run(run func(ctx context.Context, expr exp.Expression, opts ...clause.FilterOption)) *MockProfileRepo_ListBy_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]clause.FilterOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(clause.FilterOption)
			}
		}
		run(args[0].(context.Context), args[1].(exp.Expression), variadicArgs...)
	})
	return _c
}

func (_c *MockProfileRepo_ListBy_Call) Return(_a0 *response.ListResponse[*model.Profile], _a1 error) *MockProfileRepo_ListBy_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockProfileRepo_ListBy_Call) RunAndReturn(run func(context.Context, exp.Expression, ...clause.FilterOption) (*response.ListResponse[*model.Profile], error)) *MockProfileRepo_ListBy_Call {
	_c.Call.Return(run)
	return _c
}

// ListByEach provides a mock function with given fields: ctx, expr, fn, opts
func (_m *MockProfileRepo) ListByEach(ctx context.Context, expr exp.Expression, fn func(*model.Profile) error, opts ...clause.FilterOption) error {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, expr, fn)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, exp.Expression, func(*model.Profile) error, ...clause.FilterOption) error); ok {
		r0 = rf(ctx, expr, fn, opts...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockProfileRepo_ListByEach_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListByEach'
type MockProfileRepo_ListByEach_Call struct {
	*mock.Call
}

// ListByEach is a helper method to define mock.On call
//   - ctx context.Context
//   - expr exp.Expression
//   - fn func(*model.Profile) error
//   - opts ...clause.FilterOption
func (_e *MockProfileRepo_Expecter) ListByEach(ctx interface{}, expr interface{}, fn interface{}, opts ...interface{}) *MockProfileRepo_ListByEach_Call {
	return &MockProfileRepo_ListByEach_Call{Call: _e.mock.On("ListByEach",
		append([]interface{}{ctx, expr, fn}, opts...)...)}
}

func (_c *MockProfileRepo_ListByEach_Call) Run(run func(ctx context.Context, expr exp.Expression, fn func(*model.Profile) error, opts ...clause.FilterOption)) *MockProfileRepo_ListByEach_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]clause.FilterOption, len(args)-3)
		for i, a := range args[3:] {
			if a != nil {
				variadicArgs[i] = a.(clause.FilterOption)
			}
		}
		run(args[0].(context.Context), args[1].(exp.Expression), args[2].(func(*model.Profile) error), variadicArgs...)
	})
	return _c
}

func (_c *MockProfileRepo_ListByEach_Call) Return(_a0 error) *MockProfileRepo_ListByEach_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockProfileRepo_ListByEach_Call) RunAndReturn(run func(context.Context, exp.Expression, func(*model.Profile) error, ...clause.FilterOption) error) *MockProfileRepo_ListByEach_Call {
	_c.Call.Return(run)
	return _c
}

// ListByIds provides a mock function with given fields: ctx, ids
func (_m *MockProfileRepo) ListByIds(ctx context.Context, ids []int64) (*response.ListResponse[*model.Profile], error) {
	ret := _m.Called(ctx, ids)

	var r0 *response.ListResponse[*model.Profile]
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []int64) (*response.ListResponse[*model.Profile], error)); ok {
		return rf(ctx, ids)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []int64) *response.ListResponse[*model.Profile]); ok {
		r0 = rf(ctx, ids)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*response.ListResponse[*model.Profile])
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []int64) error); ok {
		r1 = rf(ctx, ids)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockProfileRepo_ListByIds_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListByIds'
type MockProfileRepo_ListByIds_Call struct {
	*mock.Call
}

// ListByIds is a helper method to define mock.On call
//   - ctx context.Context
//   - ids []int64
func (_e *MockProfileRepo_Expecter) ListByIds(ctx interface{}, ids interface{}) *MockProfileRepo_ListByIds_Call {
	return &MockProfileRepo_ListByIds_Call{Call: _e.mock.On("ListByIds", ctx, ids)}
}

func (_c *MockProfileRepo_ListByIds_Call) Run(run func(ctx context.Context, ids []int64)) *MockProfileRepo_ListByIds_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]int64))
	})
	return _c
}

func (_c *MockProfileRepo_ListByIds_Call) Return(_a0 *response.ListResponse[*model.Profile], _a1 error) *MockProfileRepo_ListByIds_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockProfileRepo_ListByIds_Call) RunAndReturn(run func(context.Context, []int64) (*response.ListResponse[*model.Profile], error)) *MockProfileRepo_ListByIds_Call {
	_c.Call.Return(run)
	return _c
}

// ListEach provides a mock function with given fields: ctx, fn, opts
func (_m *MockProfileRepo) ListEach(ctx context.Context, fn func(*model.Profile) error, opts ...clause.FilterOption) error {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, fn)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, func(*model.Profile) error, ...clause.FilterOption) error); ok {
		r0 = rf(ctx, fn, opts...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockProfileRepo_ListEach_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListEach'
type MockProfileRepo_ListEach_Call struct {
	*mock.Call
}

// ListEach is a helper method to define mock.On call
//   - ctx context.Context
//   - fn func(*model.Profile) error
//   - opts ...clause.FilterOption
func (_e *MockProfileRepo_Expecter) ListEach(ctx interface{}, fn interface{}, opts ...interface{}) *MockProfileRepo_ListEach_Call {
	return &MockProfileRepo_ListEach_Call{Call: _e.mock.On("ListEach",
		append([]interface{}{ctx, fn}, opts...)...)}
}

func (_c *MockProfileRepo_ListEach_Call) Run(run func(ctx context.Context, fn func(*model.Profile) error, opts ...clause.FilterOption)) *MockProfileRepo_ListEach_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]clause.FilterOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(clause.FilterOption)
			}
		}
		run(args[0].(context.Context), args[1].(func(*model.Profile) error), variadicArgs...)
	})
	return _c
}

func (_c *MockProfileRepo_ListEach_Call) Return(_a0 error) *MockProfileRepo_ListEach_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockProfileRepo_ListEach_Call) RunAndReturn(run func(context.Context, func(*model.Profile) error, ...clause.FilterOption) error) *MockProfileRepo_ListEach_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: ctx, req
func (_m *MockProfileRepo) Update(ctx context.Context, req *model.Profile) error {
	ret := _m.Called(ctx, req)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.Profile) error); ok {
		r0 = rf(ctx, req)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockProfileRepo_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type MockProfileRepo_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - ctx context.Context
//   - req *model.Profile
func (_e *MockProfileRepo_Expecter) Update(ctx interface{}, req interface{}) *MockProfileRepo_Update_Call {
	return &MockProfileRepo_Update_Call{Call: _e.mock.On("Update", ctx, req)}
}

func (_c *MockProfileRepo_Update_Call) Run(run func(ctx context.Context, req *model.Profile)) *MockProfileRepo_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*model.Profile))
	})
	return _c
}

func (_c *MockProfileRepo_Update_Call) Return(_a0 error) *MockProfileRepo_Update_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockProfileRepo_Update_Call) RunAndReturn(run func(context.Context, *model.Profile) error) *MockProfileRepo_Update_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateMap provides a mock function with given fields: ctx, id, req
func (_m *MockProfileRepo) UpdateMap(ctx context.Context, id int64, req map[string]interface{}) error {
	ret := _m.Called(ctx, id, req)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, map[string]interface{}) error); ok {
		r0 = rf(ctx, id, req)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockProfileRepo_UpdateMap_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateMap'
type MockProfileRepo_UpdateMap_Call struct {
	*mock.Call
}

// UpdateMap is a helper method to define mock.On call
//   - ctx context.Context
//   - id int64
//   - req map[string]interface{}
func (_e *MockProfileRepo_Expecter) UpdateMap(ctx interface{}, id interface{}, req interface{}) *MockProfileRepo_UpdateMap_Call {
	return &MockProfileRepo_UpdateMap_Call{Call: _e.mock.On("UpdateMap", ctx, id, req)}
}

func (_c *MockProfileRepo_UpdateMap_Call) Run(run func(ctx context.Context, id int64, req map[string]interface{})) *MockProfileRepo_UpdateMap_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64), args[2].(map[string]interface{}))
	})
	return _c
}

func (_c *MockProfileRepo_UpdateMap_Call) Return(_a0 error) *MockProfileRepo_UpdateMap_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockProfileRepo_UpdateMap_Call) RunAndReturn(run func(context.Context, int64, map[string]interface{}) error) *MockProfileRepo_UpdateMap_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateMapBy provides a mock function with given fields: ctx, req, expr
func (_m *MockProfileRepo) UpdateMapBy(ctx context.Context, req map[string]interface{}, expr exp.Expression) (int64, error) {
	ret := _m.Called(ctx, req, expr)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, map[string]interface{}, exp.Expression) (int64, error)); ok {
		return rf(ctx, req, expr)
	}
	if rf, ok := ret.Get(0).(func(context.Context, map[string]interface{}, exp.Expression) int64); ok {
		r0 = rf(ctx, req, expr)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, map[string]interface{}, exp.Expression) error); ok {
		r1 = rf(ctx, req, expr)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockProfileRepo_UpdateMapBy_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateMapBy'
type MockProfileRepo_UpdateMapBy_Call struct {
	*mock.Call
}

// UpdateMapBy is a helper method to define mock.On call
//   - ctx context.Context
//   - req map[string]interface{}
//   - expr exp.Expression
func (_e *MockProfileRepo_Expecter) UpdateMapBy(ctx interface{}, req interface{}, expr interface{}) *MockProfileRepo_UpdateMapBy_Call {
	return &MockProfileRepo_UpdateMapBy_Call{Call: _e.mock.On("UpdateMapBy", ctx, req, expr)}
}

func (_c *MockProfileRepo_UpdateMapBy_Call) Run(run func(ctx context.Context, req map[string]interface{}, expr exp.Expression)) *MockProfileRepo_UpdateMapBy_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(map[string]interface{}), args[2].(exp.Expression))
	})
	return _c
}

func (_c *MockProfileRepo_UpdateMapBy_Call) Return(_a0 int64, _a1 error) *MockProfileRepo_UpdateMapBy_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockProfileRepo_UpdateMapBy_Call) RunAndReturn(run func(context.Context, map[string]interface{}, exp.Expression) (int64, error)) *MockProfileRepo_UpdateMapBy_Call {
	_c.Call.Return(run)
	return _c
}

// Upsert provides a mock function with given fields: ctx, req, target
func (_m *MockProfileRepo) Upsert(ctx context.Context, req *model.Profile, target string) (bool, error) {
	ret := _m.Called(ctx, req, target)

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.Profile, string) (bool, error)); ok {
		return rf(ctx, req, target)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.Profile, string) bool); ok {
		r0 = rf(ctx, req, target)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.Profile, string) error); ok {
		r1 = rf(ctx, req, target)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockProfileRepo_Upsert_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Upsert'
type MockProfileRepo_Upsert_Call struct {
	*mock.Call
}

// Upsert is a helper method to define mock.On call
//   - ctx context.Context
//   - req *model.Profile
//   - target string
func (_e *MockProfileRepo_Expecter) Upsert(ctx interface{}, req interface{}, target interface{}) *MockProfileRepo_Upsert_Call {
	return &MockProfileRepo_Upsert_Call{Call: _e.mock.On("Upsert", ctx, req, target)}
}

func (_c *MockProfileRepo_Upsert_Call) Run(run func(ctx context.Context, req *model.Profile, target string)) *MockProfileRepo_Upsert_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*model.Profile), args[2].(string))
	})
	return _c
}

func (_c *MockProfileRepo_Upsert_Call) Return(_a0 bool, _a1 error) *MockProfileRepo_Upsert_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockProfileRepo_Upsert_Call) RunAndReturn(run func(context.Context, *model.Profile, string) (bool, error)) *MockProfileRepo_Upsert_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockProfileRepo interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockProfileRepo creates a new instance of MockProfileRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockProfileRepo(t mockConstructorTestingTNewMockProfileRepo) *MockProfileRepo {
	mock := &MockProfileRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
