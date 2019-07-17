// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// MigrationQueries is an autogenerated mock type for the MigrationQueries type
type MigrationQueries struct {
	mock.Mock
}

// CountMigrationTableSQL provides a mock function with given fields:
func (_m *MigrationQueries) CountMigrationTableSQL() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// CountSchemaSQL provides a mock function with given fields:
func (_m *MigrationQueries) CountSchemaSQL() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// CreateMigrationTableSQL provides a mock function with given fields: schema
func (_m *MigrationQueries) CreateMigrationTableSQL(schema string) string {
	ret := _m.Called(schema)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(schema)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// CreateSchemaSQL provides a mock function with given fields: schema
func (_m *MigrationQueries) CreateSchemaSQL(schema string) string {
	ret := _m.Called(schema)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(schema)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// DeleteMigrationLogSQL provides a mock function with given fields: schema
func (_m *MigrationQueries) DeleteMigrationLogSQL(schema string) string {
	ret := _m.Called(schema)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(schema)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetLatestBatchSQL provides a mock function with given fields: schema
func (_m *MigrationQueries) GetLatestBatchSQL(schema string) string {
	ret := _m.Called(schema)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(schema)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetMaxSequenceSQL provides a mock function with given fields: schema
func (_m *MigrationQueries) GetMaxSequenceSQL(schema string) string {
	ret := _m.Called(schema)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(schema)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetSequenceByBatchSQL provides a mock function with given fields: schema
func (_m *MigrationQueries) GetSequenceByBatchSQL(schema string) string {
	ret := _m.Called(schema)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(schema)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// InsertMigrationLogSQL provides a mock function with given fields: schema
func (_m *MigrationQueries) InsertMigrationLogSQL(schema string) string {
	ret := _m.Called(schema)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(schema)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// SetSchemaSQL provides a mock function with given fields: schema
func (_m *MigrationQueries) SetSchemaSQL(schema string) string {
	ret := _m.Called(schema)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(schema)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}
