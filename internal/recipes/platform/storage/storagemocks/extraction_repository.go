// Code generated by mockery v2.44.2. DO NOT EDIT.

package storagemocks

import (
	context "context"

	domain "github.com/rubenbupe/recipe-video-parser/internal/recipes/domain"
	mock "github.com/stretchr/testify/mock"
)

// ExtractionRepository is an autogenerated mock type for the ExtractionRepository type
type ExtractionRepository struct {
	mock.Mock
}

// Exists provides a mock function with given fields: ctx, id
func (_m *ExtractionRepository) Exists(ctx context.Context, id domain.ExtractionID) (bool, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for Exists")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.ExtractionID) (bool, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.ExtractionID) bool); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.ExtractionID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: ctx, id
func (_m *ExtractionRepository) Get(ctx context.Context, id domain.ExtractionID) (*domain.Extraction, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *domain.Extraction
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.ExtractionID) (*domain.Extraction, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.ExtractionID) *domain.Extraction); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Extraction)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.ExtractionID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByUserID provides a mock function with given fields: ctx, extractionId
func (_m *ExtractionRepository) GetByUserID(ctx context.Context, extractionId domain.ExtractionUserID) ([]domain.Extraction, error) {
	ret := _m.Called(ctx, extractionId)

	if len(ret) == 0 {
		panic("no return value specified for GetByUserID")
	}

	var r0 []domain.Extraction
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.ExtractionUserID) ([]domain.Extraction, error)); ok {
		return rf(ctx, extractionId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.ExtractionUserID) []domain.Extraction); ok {
		r0 = rf(ctx, extractionId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Extraction)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.ExtractionUserID) error); ok {
		r1 = rf(ctx, extractionId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: ctx, extraction
func (_m *ExtractionRepository) Save(ctx context.Context, extraction domain.Extraction) error {
	ret := _m.Called(ctx, extraction)

	if len(ret) == 0 {
		panic("no return value specified for Save")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.Extraction) error); ok {
		r0 = rf(ctx, extraction)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewExtractionRepository creates a new instance of ExtractionRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewExtractionRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *ExtractionRepository {
	mock := &ExtractionRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
