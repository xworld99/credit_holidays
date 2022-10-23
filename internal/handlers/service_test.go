package handlers

import (
	"credit_holidays/internal/mocks"
	"credit_holidays/internal/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetServicesList(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("correct", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockController := mocks.NewMockCreditHolidaysController(mockCtrl)

		handler := NewHandler(mockController)

		mockController.EXPECT().
			GetServicesList(gomock.Any()).
			Return([]models.Service{}, models.HandlerError{Err: nil}).
			Times(1)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		handler.GetServicesList(c)

		if w.Code != http.StatusOK {
			t.Error(
				"For", "empty context",
				"expected", http.StatusOK,
				"got", w.Code,
			)
		}
	})

	t.Run("incorrect", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockController := mocks.NewMockCreditHolidaysController(mockCtrl)

		handler := NewHandler(mockController)

		mockController.EXPECT().
			GetServicesList(gomock.Any()).
			Return([]models.Service{}, models.HandlerError{Type: http.StatusInternalServerError, Err: fmt.Errorf("no services")}).
			Times(1)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		handler.GetServicesList(c)

		if w.Code != http.StatusInternalServerError {
			t.Error(
				"For", "empty response",
				"expected", http.StatusInternalServerError,
				"got", w.Code,
			)
		}
	})
}
