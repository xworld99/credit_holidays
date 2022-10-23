package handlers

import (
	"bytes"
	"credit_holidays/internal/consts"
	"credit_holidays/internal/mocks"
	"credit_holidays/internal/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddOrder(t *testing.T) {
	t.Run("correct", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockController := mocks.NewMockCreditHolidaysController(mockCtrl)

		handler := NewHandler(mockController)

		mockController.EXPECT().
			AddOrder(gomock.Any(), models.AddOrderRequest{UserId: 123, ServiceId: 1, Amount: 400}).
			Return(models.Order{}, models.HandlerError{}).
			Times(1)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		reqBody := []byte(`{"user_id": 123, "service_id": 1, "amount": 400}`)
		c.Request, _ = http.NewRequest(http.MethodPost, "", bytes.NewBuffer(reqBody))

		handler.AddOrder(c)

		if w.Code != http.StatusOK {
			t.Error(
				"For", "",
				"expected", http.StatusOK,
				"got", w.Code,
			)
		}
	})

	t.Run("invalid order info", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockController := mocks.NewMockCreditHolidaysController(mockCtrl)

		handler := NewHandler(mockController)

		mockController.EXPECT().
			AddOrder(gomock.Any(), models.AddOrderRequest{UserId: 123, ServiceId: 1, Amount: 400}).
			Return(models.Order{}, models.HandlerError{Type: http.StatusBadRequest, Err: fmt.Errorf("no operation")}).
			Times(1)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		reqBody := []byte(`{"user_id": 123, "service_id": 1, "amount": 400}`)
		c.Request, _ = http.NewRequest(http.MethodPost, "", bytes.NewBuffer(reqBody))

		handler.AddOrder(c)

		if w.Code != http.StatusBadRequest {
			t.Error(
				"For", "",
				"expected", http.StatusBadRequest,
				"got", w.Code,
			)
		}
	})

	t.Run("no user id", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockController := mocks.NewMockCreditHolidaysController(mockCtrl)

		handler := NewHandler(mockController)

		mockController.EXPECT().
			AddOrder(gomock.Any(), models.AddOrderRequest{ServiceId: 1, Amount: 400}).
			Return(models.Order{}, models.HandlerError{Type: http.StatusBadRequest, Err: fmt.Errorf("no operation")}).
			Times(1)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		reqBody := []byte(`{"service_id": 1, "amount": 400}`)
		c.Request, _ = http.NewRequest(http.MethodPost, "", bytes.NewBuffer(reqBody))

		handler.AddOrder(c)

		if w.Code != http.StatusBadRequest {
			t.Error(
				"For", "",
				"expected", http.StatusBadRequest,
				"got", w.Code,
			)
		}
	})

	t.Run("no service id", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockController := mocks.NewMockCreditHolidaysController(mockCtrl)

		handler := NewHandler(mockController)

		mockController.EXPECT().
			AddOrder(gomock.Any(), models.AddOrderRequest{UserId: 123, Amount: 400}).
			Return(models.Order{}, models.HandlerError{Type: http.StatusBadRequest, Err: fmt.Errorf("no operation")}).
			Times(1)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		reqBody := []byte(`{"user_id": 123, "amount": 400}`)
		c.Request, _ = http.NewRequest(http.MethodPost, "", bytes.NewBuffer(reqBody))

		handler.AddOrder(c)

		if w.Code != http.StatusBadRequest {
			t.Error(
				"For", "",
				"expected", http.StatusBadRequest,
				"got", w.Code,
			)
		}
	})

	t.Run("cant create order", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockController := mocks.NewMockCreditHolidaysController(mockCtrl)

		handler := NewHandler(mockController)

		mockController.EXPECT().
			AddOrder(gomock.Any(), models.AddOrderRequest{UserId: 123, ServiceId: 1, Amount: 400}).
			Return(models.Order{}, models.HandlerError{Type: http.StatusInternalServerError, Err: fmt.Errorf("cant create order")}).
			Times(1)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		reqBody := []byte(`{"user_id": 123, "service_id": 1, "amount": 400}`)
		c.Request, _ = http.NewRequest(http.MethodPost, "", bytes.NewBuffer(reqBody))

		handler.AddOrder(c)

		if w.Code != http.StatusInternalServerError {
			t.Error(
				"For", "",
				"expected", http.StatusInternalServerError,
				"got", w.Code,
			)
		}
	})
}

/*
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
*/

func TestChangeOrderStatus(t *testing.T) {
	t.Run("correct", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockController := mocks.NewMockCreditHolidaysController(mockCtrl)

		handler := NewHandler(mockController)

		mockController.EXPECT().
			ChangeOrderStatus(gomock.Any(), models.ChangeOrderRequest{OrderId: 123, Action: consts.OrderProof}).
			Return(models.Order{}, models.HandlerError{}).
			Times(1)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		reqBody := []byte(`{"order_id": 123, "action": "proof"}`)
		c.Request, _ = http.NewRequest(http.MethodPost, "", bytes.NewBuffer(reqBody))

		handler.ChangeOrderStatus(c)

		if w.Code != http.StatusOK {
			t.Error(
				"For", "",
				"expected", http.StatusOK,
				"got", w.Code,
			)
		}
	})

	t.Run("invalid order info", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockController := mocks.NewMockCreditHolidaysController(mockCtrl)

		handler := NewHandler(mockController)

		mockController.EXPECT().
			ChangeOrderStatus(gomock.Any(), models.ChangeOrderRequest{OrderId: 123, Action: consts.OrderProof}).
			Return(models.Order{}, models.HandlerError{Type: http.StatusBadRequest, Err: fmt.Errorf("no operation")}).
			Times(1)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		reqBody := []byte(`{"order_id": 123, "action": "proof"}"`)
		c.Request, _ = http.NewRequest(http.MethodPost, "", bytes.NewBuffer(reqBody))

		handler.ChangeOrderStatus(c)

		if w.Code != http.StatusBadRequest {
			t.Error(
				"For", "",
				"expected", http.StatusBadRequest,
				"got", w.Code,
			)
		}
	})

	t.Run("invalid request", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockController := mocks.NewMockCreditHolidaysController(mockCtrl)

		handler := NewHandler(mockController)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		reqBody := []byte(`{ount": 400}`)
		c.Request, _ = http.NewRequest(http.MethodPost, "", bytes.NewBuffer(reqBody))

		handler.ChangeOrderStatus(c)

		if w.Code != http.StatusBadRequest {
			t.Error(
				"For", "",
				"expected", http.StatusBadRequest,
				"got", w.Code,
			)
		}
	})
}
