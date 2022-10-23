package controllers

import (
	"context"
	"credit_holidays/internal/models"
	"fmt"
)

func (c *Controller) GenerateReport(ctx context.Context, request models.GenerateReportRequest) (string, models.HandlerError) {
	var err error
	var csvdata models.CSVData

	csvdata, err = validateSaveReportParams(request)
	if err != nil {
		return "", models.CreateBadRequestError(fmt.Errorf("validation error: %w", err))
	}

	ctxTm, cancel := context.WithTimeout(ctx, c.dbTm)
	defer cancel()

	var order models.Order
	order, err = c.db.GetLastOrderMonth(ctxTm, csvdata.Period)
	if err != nil {
		return "", models.CreateNotFoundError(fmt.Errorf("cant get order: %w", err))
	}

	filepath := createReportPath(csvdata.Period, *order.CreatedAt)
	if fileAlreadyExists(c.staticPath, filepath) {
		return filepathToLink(filepath), models.HandlerError{}
	}

	csvdata, err = c.db.FormReport(ctxTm, csvdata)
	if err != nil {
		return "", models.CreateInternalError(fmt.Errorf("cant form csvfile content: %w", err))
	}

	err = saveReport(c.staticPath, filepath, csvdata)
	if err != nil {
		return "", models.CreateInternalError(fmt.Errorf("cant save csvfile: %w", err))
	}

	deleteUnnecessaryReports(c.staticPath, filepath, csvdata.Period)

	return filepathToLink(filepath), models.HandlerError{}
}
