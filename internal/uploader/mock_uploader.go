// Code generated by MockGen. DO NOT EDIT.
// Source: internal/uploader/uploader.go
//
// Generated by this command:
//
//	mockgen -source=internal/uploader/uploader.go -destination=internal/uploader/mock_uploader.go -package=uploader
//

// Package uploader is a generated GoMock package.
package uploader

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockUploader is a mock of Uploader interface.
type MockUploader struct {
	ctrl     *gomock.Controller
	recorder *MockUploaderMockRecorder
	isgomock struct{}
}

// MockUploaderMockRecorder is the mock recorder for MockUploader.
type MockUploaderMockRecorder struct {
	mock *MockUploader
}

// NewMockUploader creates a new mock instance.
func NewMockUploader(ctrl *gomock.Controller) *MockUploader {
	mock := &MockUploader{ctrl: ctrl}
	mock.recorder = &MockUploaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUploader) EXPECT() *MockUploaderMockRecorder {
	return m.recorder
}

// UploadToS3 mocks base method.
func (m *MockUploader) UploadToS3(originalURL string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadToS3", originalURL)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadToS3 indicates an expected call of UploadToS3.
func (mr *MockUploaderMockRecorder) UploadToS3(originalURL any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadToS3", reflect.TypeOf((*MockUploader)(nil).UploadToS3), originalURL)
}

// UploadToSFTP mocks base method.
func (m *MockUploader) UploadToSFTP(originalURL string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadToSFTP", originalURL)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadToSFTP indicates an expected call of UploadToSFTP.
func (mr *MockUploaderMockRecorder) UploadToSFTP(originalURL any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadToSFTP", reflect.TypeOf((*MockUploader)(nil).UploadToSFTP), originalURL)
}
