package controllers

import (
	"context"
	"credit_holidays/internal/models"
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
