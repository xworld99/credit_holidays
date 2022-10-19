package controllers

import (
	"context"
	"credit_holidays/internal/models"
	"database/sql"
	"fmt"
)

func (c *Controller) GetBalance(ctx context.Context, userId string) (models.User, error) {
	id, err := validateId(userId)
	if err != nil {
		return models.User{}, err
	}

	ctxTm, cancel := context.WithTimeout(ctx, c.selectTm)
	defer cancel()

	user := models.User{Id: id}
	user, err = c.db.GetUserById(ctxTm, user)
	if err != nil {
		return models.User{}, err
	}

	if user.Id == 0 {
		return models.User{}, fmt.Errorf("user with such id isnt exists")
	}

	return user, nil
}

func (c *Controller) getCreateUser(ctx context.Context, tx *sql.Tx, user models.User) (models.User, error) {
	ctxTm, cancel := context.WithTimeout(ctx, c.insertTm)
	defer cancel()

	user, err := c.db.GetCreateUser(ctxTm, tx, user)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
