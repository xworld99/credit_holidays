package handlers

import (
	"credit_holidays/internal/mocks"
	"credit_holidays/internal/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestGetBalance(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("correct", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockController := mocks.NewMockCreditHolidaysController(mockCtrl)

		handler := NewHandler(mockController)

		mockController.EXPECT().
			GetBalance(gomock.Any(), models.GetBalanceRequest("123")).
			Return(models.User{Id: 123}, models.HandlerError{Err: nil}).
			Times(1)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = &http.Request{}
		c.Request.URL, _ = url.Parse("?id=123")

		handler.GetBalance(c)

		if w.Code != http.StatusOK {
			t.Error(
				"For", "id=123",
				"expected", http.StatusOK,
				"got", w.Code,
			)
		}
	})

	t.Run("negative id", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockController := mocks.NewMockCreditHolidaysController(mockCtrl)

		handler := NewHandler(mockController)

		mockController.EXPECT().
			GetBalance(gomock.Any(), models.GetBalanceRequest("-123")).
			Return(models.User{}, models.HandlerError{Type: http.StatusBadRequest, Err: fmt.Errorf("negative id")}).
			Times(1)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = &http.Request{}
		c.Request.URL, _ = url.Parse("?id=-123")

		handler.GetBalance(c)

		if w.Code != http.StatusBadRequest {
			t.Error(
				"For", "id=-123",
				"expected", http.StatusBadRequest,
				"got", w.Code,
			)
		}
	})

	t.Run("zero id", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockController := mocks.NewMockCreditHolidaysController(mockCtrl)

		handler := NewHandler(mockController)

		mockController.EXPECT().
			GetBalance(gomock.Any(), models.GetBalanceRequest("0")).
			Return(models.User{}, models.HandlerError{Type: http.StatusBadRequest, Err: fmt.Errorf("negative id")}).
			Times(1)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = &http.Request{}
		c.Request.URL, _ = url.Parse("?id=0")

		handler.GetBalance(c)

		if w.Code != http.StatusBadRequest {
			t.Error(
				"For", "id=0",
				"expected", http.StatusBadRequest,
				"got", w.Code,
			)
		}
	})

	t.Run("no id", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockController := mocks.NewMockCreditHolidaysController(mockCtrl)

		handler := NewHandler(mockController)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = &http.Request{}
		c.Request.URL, _ = url.Parse("")

		handler.GetBalance(c)

		if w.Code != http.StatusBadRequest {
			t.Error(
				"For", "no id",
				"expected", http.StatusBadRequest,
				"got", w.Code,
			)
		}
	})

	t.Run("no user with such id", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockController := mocks.NewMockCreditHolidaysController(mockCtrl)

		handler := NewHandler(mockController)

		mockController.EXPECT().
			GetBalance(gomock.Any(), models.GetBalanceRequest("666")).
			Return(models.User{}, models.HandlerError{Type: http.StatusNotFound, Err: fmt.Errorf("no user")}).
			Times(1)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = &http.Request{}
		c.Request.URL, _ = url.Parse("?id=666")

		handler.GetBalance(c)

		if w.Code != http.StatusNotFound {
			t.Error(
				"For", "no user with such id",
				"expected", http.StatusNotFound,
				"got", w.Code,
			)
		}
	})
}

func TestGetUserHistory(t *testing.T) {
	t.Run("correct all fields passed", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockController := mocks.NewMockCreditHolidaysController(mockCtrl)

		handler := NewHandler(mockController)

		mockController.EXPECT().
			GetHistory(gomock.Any(), models.GetHistoryRequest{
				UserId:   "123",
				FromDate: "10-10-2001",
				ToDate:   "10-10-2022",
				Offset:   "44",
				Limit:    "55",
				OrderBy:  "created_at",
			}).
			Return(models.HistoryFrame{}, models.HandlerError{}).
			Times(1)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = &http.Request{}
		c.Request.URL, _ = url.Parse("?id=123&from_date=10-10-2001&to_date=10-10-2022&offset=44&limit=55&order_by=created_at")

		handler.GetUserHistory(c)

		if w.Code != http.StatusOK {
			t.Error(
				"For", "",
				"expected", http.StatusOK,
				"got", w.Code,
			)
		}
	})

	t.Run("correct all fields default", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockController := mocks.NewMockCreditHolidaysController(mockCtrl)

		handler := NewHandler(mockController)

		mockController.EXPECT().
			GetHistory(gomock.Any(), models.GetHistoryRequest{
				UserId:   "123",
				FromDate: "10-10-2001",
				ToDate:   "10-10-2022",
				Offset:   "0",
				Limit:    "10",
				OrderBy:  "created_at",
			}).
			Return(models.HistoryFrame{}, models.HandlerError{}).
			Times(1)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = &http.Request{}
		c.Request.URL, _ = url.Parse("?id=123&from_date=10-10-2001&to_date=10-10-2022")

		handler.GetUserHistory(c)

		if w.Code != http.StatusOK {
			t.Error(
				"For", "",
				"expected", http.StatusOK,
				"got", w.Code,
			)
		}
	})

	t.Run("from > to", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockController := mocks.NewMockCreditHolidaysController(mockCtrl)

		handler := NewHandler(mockController)

		mockController.EXPECT().
			GetHistory(gomock.Any(), models.GetHistoryRequest{
				UserId:   "123",
				FromDate: "10-10-2031",
				ToDate:   "10-10-2022",
				Offset:   "0",
				Limit:    "10",
				OrderBy:  "created_at",
			}).
			Return(models.HistoryFrame{}, models.HandlerError{Type: http.StatusBadRequest, Err: fmt.Errorf("invalid date")}).
			Times(1)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = &http.Request{}
		c.Request.URL, _ = url.Parse("?id=123&from_date=10-10-2031&to_date=10-10-2022")

		handler.GetUserHistory(c)

		if w.Code != http.StatusBadRequest {
			t.Error(
				"For", "",
				"expected", http.StatusBadRequest,
				"got", w.Code,
			)
		}
	})

	t.Run("no id", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockController := mocks.NewMockCreditHolidaysController(mockCtrl)

		handler := NewHandler(mockController)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = &http.Request{}
		c.Request.URL, _ = url.Parse("?from_date=10-10-2031&to_date=10-10-2022")

		handler.GetUserHistory(c)

		if w.Code != http.StatusBadRequest {
			t.Error(
				"For", "",
				"expected", http.StatusBadRequest,
				"got", w.Code,
			)
		}
	})

	t.Run("no from date", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockController := mocks.NewMockCreditHolidaysController(mockCtrl)

		handler := NewHandler(mockController)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = &http.Request{}
		c.Request.URL, _ = url.Parse("?id=123&to_date=10-10-2022")

		handler.GetUserHistory(c)

		if w.Code != http.StatusBadRequest {
			t.Error(
				"For", "",
				"expected", http.StatusBadRequest,
				"got", w.Code,
			)
		}
	})

	t.Run("no to date", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockController := mocks.NewMockCreditHolidaysController(mockCtrl)

		handler := NewHandler(mockController)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = &http.Request{}
		c.Request.URL, _ = url.Parse("?id=123&from_date=10-10-2022")

		handler.GetUserHistory(c)

		if w.Code != http.StatusBadRequest {
			t.Error(
				"For", "",
				"expected", http.StatusBadRequest,
				"got", w.Code,
			)
		}
	})

}
