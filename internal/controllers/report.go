package controllers

import (
	"context"
	"credit_holidays/internal/models"
	"fmt"
)

func (c *Controller) GetServicesList(ctx context.Context) ([]models.Service, error) {
	ctxTm, cancel := context.WithTimeout(ctx, c.selectTm)
	defer cancel()

	res, err := c.db.GetServicesList(ctxTm)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Controller) GetHistory(ctx context.Context, request models.GetHistoryRequest) (models.HistoryFrame, error) {
	var err error
	history, err := validateGetHistoryParams(request)
	if err != nil {
		return models.HistoryFrame{}, nil
	}

	ctxTm, cancel := context.WithTimeout(ctx, c.selectTm)
	defer cancel()

	history, err = c.db.GetHistoryFrame(ctxTm, history)
	if err != nil {
		return models.HistoryFrame{}, err
	}

	return history, nil
}

func (c *Controller) SaveReport(ctx context.Context, request models.SaveReportRequest) (string, error) {
	var err error

	csvdata, err := validateSaveReportParams(request)
	if err != nil {
		return "", err
	}

	ctxTm1, cancel1 := context.WithTimeout(ctx, c.selectTm)
	defer cancel1()
	order, err := c.db.GetLastOrderMonth(ctxTm1, csvdata.Period)
	if err != nil {
		return "", err
	}

	filepath := createReportPath(csvdata.Period, order.CreatedAt)
	if fileAlreadyExists(c.staticPath, filepath) {
		return filepath, nil
	}

	ctxTm, cancel := context.WithTimeout(ctx, c.selectTm)
	defer cancel()

	csvdata, err = c.db.FormReport(ctxTm, csvdata)
	if err != nil {
		return "", err
	}

	err = saveReport(c.staticPath, filepath, csvdata)
	if err != nil {
		return "", fmt.Errorf("internal server error")
	}

	deleteUnnecessaryReport(c.staticPath, filepath, csvdata.Period)

	return filepath, nil
}
