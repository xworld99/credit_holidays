package controllers

import (
	"context"
	"credit_holidays/internal/models"
	"fmt"
)

func (c *Controller) GetServiceInfo(ctx context.Context, service *models.Service) error {
	ctxTm, cancel := context.WithTimeout(ctx, c.dbTm)
	defer cancel()

	var err error
	*service, err = c.db.GetServiceById(ctxTm, *service)
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) GetServicesList(ctx context.Context) ([]models.Service, models.HandlerError) {
	var err error
	var res []models.Service

	ctxTm, cancel := context.WithTimeout(ctx, c.dbTm)
	defer cancel()

	res, err = c.db.GetServicesList(ctxTm)
	if err != nil {
		return nil, models.CreateNotFoundError(fmt.Errorf("cant get services list: %w", err))
	}

	return res, models.HandlerError{}
}
