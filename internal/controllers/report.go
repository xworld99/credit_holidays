package controllers

import (
	"context"
	"credit_holidays/internal/models"
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
	userId, fromDate, toDate, offset, limit, err := validateGetHistoryParams(request)
	if err != nil {
		return models.HistoryFrame{}, nil
	}

	history := models.HistoryFrame{UserId: userId, FromDate: fromDate, ToDate: toDate, Offset: offset, Limit: limit}

	ctxTm, cancel := context.WithTimeout(ctx, c.selectTm)
	defer cancel()

	history, err = c.db.GetHistoryFrame(ctxTm, history)
	if err != nil {
		return models.HistoryFrame{}, err
	}

	return history, nil
}
