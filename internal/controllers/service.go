package controllers

import (
	"context"
	"credit_holidays/internal/models"
)

func (c *Controller) getServiceInfo(ctx context.Context, service models.Service) (models.Service, error) {
	ctxTm, cancel := context.WithTimeout(ctx, c.selectTm)
	defer cancel()

	service, err := c.db.GetServiceById(ctxTm, service)
	if err != nil {
		return models.Service{}, err
	}

	return service, nil
}
