// Code generated by mockery v2.45.0. DO NOT EDIT.

package mocks

import (
	context "context"

	dynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"

	mock "github.com/stretchr/testify/mock"

	models "github.com/kmdeveloping/dynalock/internal/models"
)

// LockManagerProvider is an autogenerated mock type for the LockManagerProvider type
type LockManagerProvider struct {
	mock.Mock
}

type LockManagerProvider_Expecter struct {
	mock *mock.Mock
}

func (_m *LockManagerProvider) EXPECT() *LockManagerProvider_Expecter {
	return &LockManagerProvider_Expecter{mock: &_m.Mock}
}

// AcquireLock provides a mock function with given fields: ctx, opt
func (_m *LockManagerProvider) AcquireLock(ctx context.Context, opt *models.AcquireLockOptions) (*models.Lock, error) {
	ret := _m.Called(ctx, opt)

	if len(ret) == 0 {
		panic("no return value specified for AcquireLock")
	}

	var r0 *models.Lock
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.AcquireLockOptions) (*models.Lock, error)); ok {
		return rf(ctx, opt)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *models.AcquireLockOptions) *models.Lock); ok {
		r0 = rf(ctx, opt)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Lock)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *models.AcquireLockOptions) error); ok {
		r1 = rf(ctx, opt)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LockManagerProvider_AcquireLock_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AcquireLock'
type LockManagerProvider_AcquireLock_Call struct {
	*mock.Call
}

// AcquireLock is a helper method to define mock.On call
//   - ctx context.Context
//   - opt *models.AcquireLockOptions
func (_e *LockManagerProvider_Expecter) AcquireLock(ctx interface{}, opt interface{}) *LockManagerProvider_AcquireLock_Call {
	return &LockManagerProvider_AcquireLock_Call{Call: _e.mock.On("AcquireLock", ctx, opt)}
}

func (_c *LockManagerProvider_AcquireLock_Call) Run(run func(ctx context.Context, opt *models.AcquireLockOptions)) *LockManagerProvider_AcquireLock_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*models.AcquireLockOptions))
	})
	return _c
}

func (_c *LockManagerProvider_AcquireLock_Call) Return(_a0 *models.Lock, _a1 error) *LockManagerProvider_AcquireLock_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *LockManagerProvider_AcquireLock_Call) RunAndReturn(run func(context.Context, *models.AcquireLockOptions) (*models.Lock, error)) *LockManagerProvider_AcquireLock_Call {
	_c.Call.Return(run)
	return _c
}

// CanAcquireLock provides a mock function with given fields: ctx, key
func (_m *LockManagerProvider) CanAcquireLock(ctx context.Context, key string) error {
	ret := _m.Called(ctx, key)

	if len(ret) == 0 {
		panic("no return value specified for CanAcquireLock")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// LockManagerProvider_CanAcquireLock_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CanAcquireLock'
type LockManagerProvider_CanAcquireLock_Call struct {
	*mock.Call
}

// CanAcquireLock is a helper method to define mock.On call
//   - ctx context.Context
//   - key string
func (_e *LockManagerProvider_Expecter) CanAcquireLock(ctx interface{}, key interface{}) *LockManagerProvider_CanAcquireLock_Call {
	return &LockManagerProvider_CanAcquireLock_Call{Call: _e.mock.On("CanAcquireLock", ctx, key)}
}

func (_c *LockManagerProvider_CanAcquireLock_Call) Run(run func(ctx context.Context, key string)) *LockManagerProvider_CanAcquireLock_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *LockManagerProvider_CanAcquireLock_Call) Return(_a0 error) *LockManagerProvider_CanAcquireLock_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *LockManagerProvider_CanAcquireLock_Call) RunAndReturn(run func(context.Context, string) error) *LockManagerProvider_CanAcquireLock_Call {
	_c.Call.Return(run)
	return _c
}

// CreateLockTable provides a mock function with given fields: ctx, opt
func (_m *LockManagerProvider) CreateLockTable(ctx context.Context, opt *models.CreateDynamoDBTableOptions) (*dynamodb.CreateTableOutput, error) {
	ret := _m.Called(ctx, opt)

	if len(ret) == 0 {
		panic("no return value specified for CreateLockTable")
	}

	var r0 *dynamodb.CreateTableOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.CreateDynamoDBTableOptions) (*dynamodb.CreateTableOutput, error)); ok {
		return rf(ctx, opt)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *models.CreateDynamoDBTableOptions) *dynamodb.CreateTableOutput); ok {
		r0 = rf(ctx, opt)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dynamodb.CreateTableOutput)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *models.CreateDynamoDBTableOptions) error); ok {
		r1 = rf(ctx, opt)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LockManagerProvider_CreateLockTable_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateLockTable'
type LockManagerProvider_CreateLockTable_Call struct {
	*mock.Call
}

// CreateLockTable is a helper method to define mock.On call
//   - ctx context.Context
//   - opt *models.CreateDynamoDBTableOptions
func (_e *LockManagerProvider_Expecter) CreateLockTable(ctx interface{}, opt interface{}) *LockManagerProvider_CreateLockTable_Call {
	return &LockManagerProvider_CreateLockTable_Call{Call: _e.mock.On("CreateLockTable", ctx, opt)}
}

func (_c *LockManagerProvider_CreateLockTable_Call) Run(run func(ctx context.Context, opt *models.CreateDynamoDBTableOptions)) *LockManagerProvider_CreateLockTable_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*models.CreateDynamoDBTableOptions))
	})
	return _c
}

func (_c *LockManagerProvider_CreateLockTable_Call) Return(_a0 *dynamodb.CreateTableOutput, _a1 error) *LockManagerProvider_CreateLockTable_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *LockManagerProvider_CreateLockTable_Call) RunAndReturn(run func(context.Context, *models.CreateDynamoDBTableOptions) (*dynamodb.CreateTableOutput, error)) *LockManagerProvider_CreateLockTable_Call {
	_c.Call.Return(run)
	return _c
}

// ReleaseLock provides a mock function with given fields: ctx, key
func (_m *LockManagerProvider) ReleaseLock(ctx context.Context, key string) (bool, error) {
	ret := _m.Called(ctx, key)

	if len(ret) == 0 {
		panic("no return value specified for ReleaseLock")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (bool, error)); ok {
		return rf(ctx, key)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) bool); ok {
		r0 = rf(ctx, key)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LockManagerProvider_ReleaseLock_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ReleaseLock'
type LockManagerProvider_ReleaseLock_Call struct {
	*mock.Call
}

// ReleaseLock is a helper method to define mock.On call
//   - ctx context.Context
//   - key string
func (_e *LockManagerProvider_Expecter) ReleaseLock(ctx interface{}, key interface{}) *LockManagerProvider_ReleaseLock_Call {
	return &LockManagerProvider_ReleaseLock_Call{Call: _e.mock.On("ReleaseLock", ctx, key)}
}

func (_c *LockManagerProvider_ReleaseLock_Call) Run(run func(ctx context.Context, key string)) *LockManagerProvider_ReleaseLock_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *LockManagerProvider_ReleaseLock_Call) Return(_a0 bool, _a1 error) *LockManagerProvider_ReleaseLock_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *LockManagerProvider_ReleaseLock_Call) RunAndReturn(run func(context.Context, string) (bool, error)) *LockManagerProvider_ReleaseLock_Call {
	_c.Call.Return(run)
	return _c
}

// NewLockManagerProvider creates a new instance of LockManagerProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewLockManagerProvider(t interface {
	mock.TestingT
	Cleanup(func())
}) *LockManagerProvider {
	mock := &LockManagerProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
