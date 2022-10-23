package controllers

import (
	"context"
	"credit_holidays/internal/models"
	"fmt"
	"net/http"
)

func (c *Controller) getServiceInfo(ctx context.Context, service *models.Service) error {
	ctxTm, cancel := context.WithTimeout(ctx, c.dbTm)
	defer cancel()

	var err error
	*service, err = c.db.GetServiceById(ctxTm, *service)
	if err != nil {
		return err
	}

	return nil
}

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
