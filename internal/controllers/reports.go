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
