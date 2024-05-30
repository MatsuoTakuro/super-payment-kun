// Code generated by MockGen. DO NOT EDIT.
// Source: service.go
//
// Generated by this command:
//
//	mockgen -source=service.go -destination=service_mock.go -package=handler
//

// Package handler is a generated GoMock package.
package handler

import (
	context "context"
	reflect "reflect"
	model "super-payment-kun/internal/model"
	time "time"

	gomock "go.uber.org/mock/gomock"
)

// MockCreateInvoiceSvc is a mock of CreateInvoiceSvc interface.
type MockCreateInvoiceSvc struct {
	ctrl     *gomock.Controller
	recorder *MockCreateInvoiceSvcMockRecorder
}

// MockCreateInvoiceSvcMockRecorder is the mock recorder for MockCreateInvoiceSvc.
type MockCreateInvoiceSvcMockRecorder struct {
	mock *MockCreateInvoiceSvc
}

// NewMockCreateInvoiceSvc creates a new mock instance.
func NewMockCreateInvoiceSvc(ctrl *gomock.Controller) *MockCreateInvoiceSvc {
	mock := &MockCreateInvoiceSvc{ctrl: ctrl}
	mock.recorder = &MockCreateInvoiceSvcMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCreateInvoiceSvc) EXPECT() *MockCreateInvoiceSvcMockRecorder {
	return m.recorder
}

// CreateInvoice mocks base method.
func (m *MockCreateInvoiceSvc) CreateInvoice(ctx context.Context, user model.User, vendor model.Vendor, vendorBankAcount model.VendorBankAccount, paymentAmount float64, dueDate time.Time) (*model.Invoice, *model.Vendor, *model.VendorBankAccount, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateInvoice", ctx, user, vendor, vendorBankAcount, paymentAmount, dueDate)
	ret0, _ := ret[0].(*model.Invoice)
	ret1, _ := ret[1].(*model.Vendor)
	ret2, _ := ret[2].(*model.VendorBankAccount)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// CreateInvoice indicates an expected call of CreateInvoice.
func (mr *MockCreateInvoiceSvcMockRecorder) CreateInvoice(ctx, user, vendor, vendorBankAcount, paymentAmount, dueDate any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateInvoice", reflect.TypeOf((*MockCreateInvoiceSvc)(nil).CreateInvoice), ctx, user, vendor, vendorBankAcount, paymentAmount, dueDate)
}

// MockGetInvoicesSvc is a mock of GetInvoicesSvc interface.
type MockGetInvoicesSvc struct {
	ctrl     *gomock.Controller
	recorder *MockGetInvoicesSvcMockRecorder
}

// MockGetInvoicesSvcMockRecorder is the mock recorder for MockGetInvoicesSvc.
type MockGetInvoicesSvcMockRecorder struct {
	mock *MockGetInvoicesSvc
}

// NewMockGetInvoicesSvc creates a new mock instance.
func NewMockGetInvoicesSvc(ctrl *gomock.Controller) *MockGetInvoicesSvc {
	mock := &MockGetInvoicesSvc{ctrl: ctrl}
	mock.recorder = &MockGetInvoicesSvcMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGetInvoicesSvc) EXPECT() *MockGetInvoicesSvcMockRecorder {
	return m.recorder
}

// GetInvoices mocks base method.
func (m *MockGetInvoicesSvc) GetInvoices(ctx context.Context, user model.User, fromDueDate, toDueDate time.Time, limit int, lastID string, lastDueDate time.Time, direction string) ([]model.Invoice, string, time.Time, bool, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInvoices", ctx, user, fromDueDate, toDueDate, limit, lastID, lastDueDate, direction)
	ret0, _ := ret[0].([]model.Invoice)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(time.Time)
	ret3, _ := ret[3].(bool)
	ret4, _ := ret[4].(bool)
	ret5, _ := ret[5].(error)
	return ret0, ret1, ret2, ret3, ret4, ret5
}

// GetInvoices indicates an expected call of GetInvoices.
func (mr *MockGetInvoicesSvcMockRecorder) GetInvoices(ctx, user, fromDueDate, toDueDate, limit, lastID, lastDueDate, direction any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInvoices", reflect.TypeOf((*MockGetInvoicesSvc)(nil).GetInvoices), ctx, user, fromDueDate, toDueDate, limit, lastID, lastDueDate, direction)
}
