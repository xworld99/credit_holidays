package controllers

import (
	"context"
	"credit_holidays/internal/models"
	"fmt"
	"net/http"
)

func (c *Controller) GetServicesList(ctx context.Context) ([]models.Service, models.InternalError) {
	var err models.InternalError
	var res []models.Service

	ctxTm, cancel := context.WithTimeout(ctx, c.dbTm)
	defer cancel()

	res, err.Err = c.db.GetServicesList(ctxTm)
	if err.Err != nil {
		err.Type = http.StatusInternalServerError
		err.Err = fmt.Errorf("cant get services list: %w", err.Err)
		return nil, err
	}

	return res, err
}

func (c *Controller) GenerateReport(ctx context.Context, request models.GenerateReportRequest) (string, models.InternalError) {
	var err models.InternalError
	var csvdata models.CSVData

	csvdata, err.Err = validateSaveReportParams(request)
	if err.Err != nil {
		err.Type = http.StatusBadRequest
		err.Err = fmt.Errorf("validation error: %w", err.Err)
		return "", err
	}

	ctxTm, cancel := context.WithTimeout(ctx, c.dbTm)
	defer cancel()

	var order models.Order
	order, err.Err = c.db.GetLastOrderMonth(ctxTm, csvdata.Period)
	if err.Err != nil {
		err.Type = http.StatusNotFound
		err.Err = fmt.Errorf("cant get order: %w", err.Err)
		return "", err
	}

	filepath := createReportPath(csvdata.Period, order.CreatedAt)
	if fileAlreadyExists(c.staticPath, filepath) {
		return filepath, err
	}

	csvdata, err.Err = c.db.FormReport(ctxTm, csvdata)
	if err.Err != nil {
		err.Type = http.StatusInternalServerError
		err.Err = fmt.Errorf("cant form csvfile content: %w", err.Err)
		return "", err
	}

	err.Err = saveReport(c.staticPath, filepath, csvdata)
	if err.Err != nil {
		err.Type = http.StatusInternalServerError
		err.Err = fmt.Errorf("cant save csvfile: %w", err.Err)
		return "", err
	}

	deleteUnnecessaryReports(c.staticPath, filepath, csvdata.Period)

	return filepath, err
}
