// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// FileReader is an autogenerated mock type for the FileReader type
type FileReader struct {
	mock.Mock
}

// ReadDirs provides a mock function with given fields: root
func (_m *FileReader) ReadDirs(root string) ([]string, error) {
	ret := _m.Called(root)

	var r0 []string
	if rf, ok := ret.Get(0).(func(string) []string); ok {
		r0 = rf(root)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(root)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReadFileAsString provides a mock function with given fields: path
func (_m *FileReader) ReadFileAsString(path string) (string, error) {
	ret := _m.Called(path)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(path)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReadFilesWithExtension provides a mock function with given fields: root, extension
func (_m *FileReader) ReadFilesWithExtension(root string, extension string) ([]string, error) {
	ret := _m.Called(root, extension)

	var r0 []string
	if rf, ok := ret.Get(0).(func(string, string) []string); ok {
		r0 = rf(root, extension)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(root, extension)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
