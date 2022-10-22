package controllers

import (
	"credit_holidays/internal/consts"
	"credit_holidays/internal/models"
	"fmt"
	"testing"
	"time"
)

func TestParseTime(t *testing.T) {
	t.Run("correct", func(t *testing.T) {
		_time := "22-10-2011"

		parsedTime, err := parseTime(_time)
		if err != nil {
			t.Error(
				"For", _time,
				"expected", _time,
				"got", err.Error(),
			)
		}

		if parsedTime.Format("02-01-2006") != _time {
			t.Error(
				"For", _time,
				"expected", _time,
				"got", parsedTime.Format("02-01-2006"),
			)
		}
	})

	t.Run("incorrect", func(t *testing.T) {
		times := []string{"2022-12-12", "11:00 12-12-2022", "10-2011", "2011", "2011-11", "99-99-9999"}

		for _, _time := range times {
			parsedTime, err := parseTime(_time)
			if err == nil {
				t.Error(
					"For", _time,
					"expected", "error",
					"got", parsedTime,
				)
			}
		}
	})
}

func TestValidateId(t *testing.T) {
	t.Run("correct", func(t *testing.T) {
		id := "1123"

		parsedId, err := validateId(id)

		if err != nil {
			t.Error(
				"For", id,
				"expected", 1123,
				"got", err.Error(),
			)
		}

		if parsedId != 1123 {
			t.Error(
				"For", id,
				"expected", 1123,
				"got", parsedId,
			)
		}
	})

	t.Run("incorrect", func(t *testing.T) {
		ids := []string{"-123", "0", "asdasd", "123a"}

		for _, id := range ids {
			parsedId, err := parseTime(id)
			if err == nil {
				t.Error(
					"For", id,
					"expected", "error",
					"got", parsedId,
				)
			}
		}
	})
}

func TestValidateOrderParams(t *testing.T) {
	t.Run("correct", func(t *testing.T) {
		orderReq := models.AddOrderRequest{UserId: 123, ServiceId: 123, Amount: 10000}

		if err := validateOrderParams(orderReq); err != nil {
			t.Error(
				"For", orderReq,
				"expected", "nil error",
				"got", err.Error(),
			)
		}
	})

	t.Run("incorrect", func(t *testing.T) {
		reqs := []models.AddOrderRequest{
			{UserId: -123, ServiceId: 123, Amount: 1333},
			{UserId: 123, ServiceId: -123, Amount: 1333},
			{UserId: 123, ServiceId: 123, Amount: -1333},
			{UserId: 0, ServiceId: 123, Amount: 1333},
			{UserId: 123, ServiceId: 0, Amount: 1333},
			{UserId: 123, ServiceId: 123, Amount: 0},
			{UserId: 0, ServiceId: 0, Amount: 0},
		}

		for _, r := range reqs {
			if err := validateOrderParams(r); err == nil {
				t.Error(
					"For", r,
					"expected", "error",
					"got", "valid order request",
				)
			}
		}
	})
}

func TestValidateChangeStatusParams(t *testing.T) {
	t.Run("correct", func(t *testing.T) {
		reqs := []models.ChangeOrderRequest{
			{OrderId: 123, Action: consts.OrderProof},
			{OrderId: 123, Action: consts.OrderDecline},
		}

		for _, r := range reqs {
			if err := validateChangeStatusParams(r); err != nil {
				t.Error(
					"For", r,
					"expected", "nil error",
					"got", err.Error(),
				)
			}
		}
	})

	t.Run("incorect", func(t *testing.T) {
		reqs := []models.ChangeOrderRequest{
			{OrderId: 123, Action: "asdasdasdasd"},
			{OrderId: -123, Action: consts.OperationWithdraw},
			{OrderId: 0, Action: consts.OperationAccrual},
			{OrderId: -123, Action: "asdasd"},
		}

		for _, r := range reqs {
			if err := validateChangeStatusParams(r); err == nil {
				t.Error(
					"For", r,
					"expected", "error",
					"got", "valid request",
				)
			}
		}
	})
}

func TestValidateNonNegativeInt(t *testing.T) {
	t.Run("correct", func(t *testing.T) {
		ids := []string{"1123", "0"}

		for _, id := range ids {
			_, err := validateNonNegativeInt(id)

			if err != nil {
				t.Error(
					"For", id,
					"expected", id,
					"got", err.Error(),
				)
			}
		}
	})

	t.Run("incorrect", func(t *testing.T) {
		ids := []string{"-1123", "-123123"}

		for _, id := range ids {
			val, err := validateNonNegativeInt(id)

			if err == nil {
				t.Error(
					"For", id,
					"expected", "error",
					"got", val,
				)
			}
		}
	})
}

func TestValidateGetHistoryParams(t *testing.T) {
	t.Run("correct", func(t *testing.T) {
		reqs := []models.GetHistoryRequest{
			{UserId: "123", FromDate: "12-12-2022", ToDate: "14-12-2022", Offset: "3", Limit: "13", OrderBy: "created_at"},
			{UserId: "123", FromDate: "12-12-2022", ToDate: "14-12-2022", Offset: "3", Limit: "13", OrderBy: "amount"},
		}

		for _, r := range reqs {
			_, err := validateGetHistoryParams(r)
			if err != nil {
				t.Error(
					"For", r,
					"expected", "success",
					"got", err.Error(),
				)
			}
		}
	})

	t.Run("incorrect", func(t *testing.T) {
		reqs := []models.GetHistoryRequest{
			{UserId: "-123", FromDate: "12-12-2022", ToDate: "14-12-2022", Offset: "3", Limit: "13", OrderBy: "created_at"},
			{UserId: "123", FromDate: "2022-12-12", ToDate: "14-12-2022", Offset: "3", Limit: "13", OrderBy: "amount"},
			{UserId: "123", FromDate: "99-12-2033", ToDate: "14-12-2022", Offset: "3", Limit: "13", OrderBy: "amount"},
			{UserId: "123", FromDate: "12-12-2022", ToDate: "44-12-2022", Offset: "3", Limit: "13", OrderBy: "amount"},
			{UserId: "123", FromDate: "15-12-2022", ToDate: "14-12-2022", Offset: "3", Limit: "13", OrderBy: "amount"},
			{UserId: "123", FromDate: "12-12-2022", ToDate: "14-12-2022", Offset: "-3", Limit: "13", OrderBy: "amount"},
			{UserId: "123", FromDate: "12-12-2022", ToDate: "14-12-2022", Offset: "3", Limit: "13", OrderBy: "asdasda"},
		}

		for _, r := range reqs {
			f, err := validateGetHistoryParams(r)
			if err == nil {
				t.Error(
					"For", r,
					"expected", "validation error",
					"got", f,
				)
			}
		}
	})
}

func TestValidateSaveReportParams(t *testing.T) {
	t.Run("correct", func(t *testing.T) {

		month := models.GenerateReportRequest("10-2010")

		if _, err := validateSaveReportParams(month); err != nil {
			t.Error(
				"For", month,
				"expected", "correct validation",
				"got", err.Error(),
			)
		}
	})

	t.Run("incorrect", func(t *testing.T) {
		months := []models.GenerateReportRequest{
			"93-2010",
			"2010-12",
			"11-11-2010",
			"201-1203-1233",
		}

		for _, month := range months {
			if _, err := validateSaveReportParams(month); err == nil {
				t.Error(
					"For", month,
					"expected", "error",
					"got", "correct validation",
				)
			}
		}
	})
}

func TestCreateReportPath(t *testing.T) {
	period := time.Date(2021, time.Month(2), 21, 1, 10, 30, 0, time.UTC)

	ans := fmt.Sprintf("2-2021-%d.csv", period.Unix())
	res := createReportPath(period, period)

	if createReportPath(period, period) != ans {
		t.Error(
			"For", period,
			"expected", ans,
			"got", res,
		)
	}
}
