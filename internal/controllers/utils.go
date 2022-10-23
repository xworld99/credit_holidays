package controllers

import (
	"credit_holidays/internal/consts"
	"credit_holidays/internal/models"
	"encoding/csv"
	"fmt"
	"os"
	fp "path/filepath"
	"strconv"
	"time"
)

func parseTime(t string) (time.Time, error) {
	layout := "02-01-2006"
	res, err := time.Parse(layout, t)
	if err != nil {
		return time.Time{}, fmt.Errorf("time isnt fits layout: '%s': %w", layout, err)
	}

	return res, nil
}

func validateId(id string) (int64, error) {
	valId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("id is an invalid int")
	}
	if valId <= 0 {
		return 0, fmt.Errorf("id should be positive")
	}
	return valId, nil
}

func validateOrderParams(request models.AddOrderRequest) error {
	if request.UserId <= 0 {
		return fmt.Errorf("user id should be > 0")
	}
	if request.ServiceId <= 0 {
		return fmt.Errorf("service id should be > 0")
	}
	if request.Amount <= 0 {
		return fmt.Errorf("amount should be > 0")
	}

	return nil
}

func validateChangeStatusParams(request models.ChangeOrderRequest) error {
	if request.OrderId <= 0 {
		return fmt.Errorf("order id should be > 0")
	}

	if _, ok := consts.OrderActions[request.Action]; !ok {
		return fmt.Errorf("unknown action: %s", request.Action)
	}

	return nil
}

func validateNonNegativeInt(i string) (int, error) {
	val, err := strconv.Atoi(i)
	if err != nil {
		return 0, fmt.Errorf("param is an invalid int: %w", err)
	}
	if val < 0 {
		return 0, fmt.Errorf("param should be non negative")
	}
	return val, nil
}

func validateGetHistoryParams(request models.GetHistoryRequest) (models.HistoryFrame, error) {
	var frame models.HistoryFrame
	var err error

	frame.UserId, err = validateId(request.UserId)
	if err != nil {
		return models.HistoryFrame{}, err
	}

	frame.FromDate, err = parseTime(request.FromDate)
	if err != nil {
		return models.HistoryFrame{}, err
	}
	frame.ToDate, err = parseTime(request.ToDate)
	if err != nil {
		return models.HistoryFrame{}, err
	}
	if frame.FromDate.After(frame.ToDate) {
		return models.HistoryFrame{}, fmt.Errorf("from date is greater then to date")
	}

	frame.Offset, err = validateNonNegativeInt(request.Offset)
	if err != nil {
		return models.HistoryFrame{}, fmt.Errorf("offset %s should be non negative int: %w", request.Offset, err)
	}

	frame.Limit, err = validateNonNegativeInt(request.Limit)
	if err != nil {
		return models.HistoryFrame{}, fmt.Errorf("limit %s should be non negative int: %w", request.Limit, err)
	}

	if _, ok := consts.SortingType[request.OrderBy]; !ok {
		return models.HistoryFrame{}, fmt.Errorf("unknown order by type")
	}
	frame.OrderBy = request.OrderBy

	return frame, nil
}

func validateSaveReportParams(request models.GenerateReportRequest) (models.CSVData, error) {
	var res models.CSVData
	var err error

	res.Period, err = parseTime("01-" + string(request))
	if err != nil {
		return models.CSVData{}, err
	}

	return res, nil
}

func createReportPath(period, lastOperation time.Time) string {
	return fmt.Sprintf("%d-%d-%d.csv", period.Month(), period.Year(), lastOperation.Unix())
}

func fileAlreadyExists(dir, name string) bool {
	if _, err := os.Stat(fmt.Sprintf("%s/%s", dir, name)); err == nil {
		return true
	}
	return false
}

func saveReport(dir, filepath string, data models.CSVData) error {
	f, err := os.Create(fmt.Sprintf("%s/%s", dir, filepath))
	defer f.Close()

	if err != nil {
		return fmt.Errorf("cant create file: %w", err)
	}

	writer := csv.NewWriter(f)
	defer writer.Flush()

	for _, r := range data.ToStringSlice() {
		if err = writer.Write(r); err != nil {
			return fmt.Errorf("cant write data in file: %w", err)
		}
	}

	return nil
}

func deleteUnnecessaryReports(dir, filepath string, period time.Time) {
	absPath := fmt.Sprintf("%s/%s", dir, filepath)
	files, _ := fp.Glob(fmt.Sprintf("%s/%d-%d-*.csv", dir, period.Month(), period.Year()))

	if files == nil {
		return
	}

	for _, f := range files {
		if f == absPath {
			continue
		}
		os.Remove(f)
	}
}

func filepathToLink(path string) string {
	return fmt.Sprintf(consts.ReportLinkPattern, consts.Host, path)
}
