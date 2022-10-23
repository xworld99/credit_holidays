// Code generated by MockGen. DO NOT EDIT.
// Source: credit_holidays/internal/controllers (interfaces: CreditHolidaysController)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	models "credit_holidays/internal/models"
	sql "database/sql"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockCreditHolidaysController is a mock of CreditHolidaysController interface.
type MockCreditHolidaysController struct {
	ctrl     *gomock.Controller
	recorder *MockCreditHolidaysControllerMockRecorder
}

// MockCreditHolidaysControllerMockRecorder is the mock recorder for MockCreditHolidaysController.
type MockCreditHolidaysControllerMockRecorder struct {
	mock *MockCreditHolidaysController
}

// NewMockCreditHolidaysController creates a new mock instance.
func NewMockCreditHolidaysController(ctrl *gomock.Controller) *MockCreditHolidaysController {
	mock := &MockCreditHolidaysController{ctrl: ctrl}
	mock.recorder = &MockCreditHolidaysControllerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCreditHolidaysController) EXPECT() *MockCreditHolidaysControllerMockRecorder {
	return m.recorder
}

// AddOrder mocks base method.
func (m *MockCreditHolidaysController) AddOrder(arg0 context.Context, arg1 models.AddOrderRequest) (models.Order, models.HandlerError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddOrder", arg0, arg1)
	ret0, _ := ret[0].(models.Order)
	ret1, _ := ret[1].(models.HandlerError)
	return ret0, ret1
}

// AddOrder indicates an expected call of AddOrder.
func (mr *MockCreditHolidaysControllerMockRecorder) AddOrder(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddOrder", reflect.TypeOf((*MockCreditHolidaysController)(nil).AddOrder), arg0, arg1)
}

// ChangeOrderStatus mocks base method.
func (m *MockCreditHolidaysController) ChangeOrderStatus(arg0 context.Context, arg1 models.ChangeOrderRequest) (models.Order, models.HandlerError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeOrderStatus", arg0, arg1)
	ret0, _ := ret[0].(models.Order)
	ret1, _ := ret[1].(models.HandlerError)
	return ret0, ret1
}

// ChangeOrderStatus indicates an expected call of ChangeOrderStatus.
func (mr *MockCreditHolidaysControllerMockRecorder) ChangeOrderStatus(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeOrderStatus", reflect.TypeOf((*MockCreditHolidaysController)(nil).ChangeOrderStatus), arg0, arg1)
}

// Close mocks base method.
func (m *MockCreditHolidaysController) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockCreditHolidaysControllerMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockCreditHolidaysController)(nil).Close))
}

// GenerateReport mocks base method.
func (m *MockCreditHolidaysController) GenerateReport(arg0 context.Context, arg1 models.GenerateReportRequest) (string, models.HandlerError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateReport", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(models.HandlerError)
	return ret0, ret1
}

// GenerateReport indicates an expected call of GenerateReport.
func (mr *MockCreditHolidaysControllerMockRecorder) GenerateReport(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateReport", reflect.TypeOf((*MockCreditHolidaysController)(nil).GenerateReport), arg0, arg1)
}

// GetBalance mocks base method.
func (m *MockCreditHolidaysController) GetBalance(arg0 context.Context, arg1 models.GetBalanceRequest) (models.User, models.HandlerError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBalance", arg0, arg1)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(models.HandlerError)
	return ret0, ret1
}

// GetBalance indicates an expected call of GetBalance.
func (mr *MockCreditHolidaysControllerMockRecorder) GetBalance(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBalance", reflect.TypeOf((*MockCreditHolidaysController)(nil).GetBalance), arg0, arg1)
}

// GetHistory mocks base method.
func (m *MockCreditHolidaysController) GetHistory(arg0 context.Context, arg1 models.GetHistoryRequest) (models.HistoryFrame, models.HandlerError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHistory", arg0, arg1)
	ret0, _ := ret[0].(models.HistoryFrame)
	ret1, _ := ret[1].(models.HandlerError)
	return ret0, ret1
}

// GetHistory indicates an expected call of GetHistory.
func (mr *MockCreditHolidaysControllerMockRecorder) GetHistory(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHistory", reflect.TypeOf((*MockCreditHolidaysController)(nil).GetHistory), arg0, arg1)
}

// GetServiceInfo mocks base method.
func (m *MockCreditHolidaysController) GetServiceInfo(arg0 context.Context, arg1 *models.Service) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetServiceInfo", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetServiceInfo indicates an expected call of GetServiceInfo.
func (mr *MockCreditHolidaysControllerMockRecorder) GetServiceInfo(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetServiceInfo", reflect.TypeOf((*MockCreditHolidaysController)(nil).GetServiceInfo), arg0, arg1)
}

// GetServicesList mocks base method.
func (m *MockCreditHolidaysController) GetServicesList(arg0 context.Context) ([]models.Service, models.HandlerError) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetServicesList", arg0)
	ret0, _ := ret[0].([]models.Service)
	ret1, _ := ret[1].(models.HandlerError)
	return ret0, ret1
}

// GetServicesList indicates an expected call of GetServicesList.
func (mr *MockCreditHolidaysControllerMockRecorder) GetServicesList(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetServicesList", reflect.TypeOf((*MockCreditHolidaysController)(nil).GetServicesList), arg0)
}

// InsertUserIfNotExists mocks base method.
func (m *MockCreditHolidaysController) InsertUserIfNotExists(arg0 context.Context, arg1 *sql.Tx, arg2 *models.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertUserIfNotExists", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertUserIfNotExists indicates an expected call of InsertUserIfNotExists.
func (mr *MockCreditHolidaysControllerMockRecorder) InsertUserIfNotExists(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertUserIfNotExists", reflect.TypeOf((*MockCreditHolidaysController)(nil).InsertUserIfNotExists), arg0, arg1, arg2)
}