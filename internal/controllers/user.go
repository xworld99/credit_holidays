package controllers

import (
	"context"
	"credit_holidays/internal/consts"
	"credit_holidays/internal/models"
	"database/sql"
	"fmt"
	"net/http"
)

func (c *Controller) GetBalance(
	ctx context.Context,
	userId models.GetBalanceRequest,
) (models.User, models.InternalError) {
	var err models.InternalError

	var id int64
	id, err.Err = validateId(string(userId))
	if err.Err != nil {
		err.Type = http.StatusBadRequest
		err.Err = fmt.Errorf("validation error: %w", err.Err)
		return models.User{}, err
	}

	ctxTm, cancel := context.WithTimeout(ctx, c.dbTm)
	defer cancel()

	user := models.User{Id: id}
	user, err.Err = c.db.GetUserById(ctxTm, user)
	if err.Err != nil {
		err.Type = http.StatusInternalServerError
		err.Err = fmt.Errorf("cant get user by id: %w", err.Err)
		return models.User{}, err
	}

	return user, err
}

func (c *Controller) getCreateUser(ctx context.Context, tx *sql.Tx, user *models.User) error {
	ctxTm, cancel := context.WithTimeout(ctx, c.dbTm)
	defer cancel()

	var err error
	*user, err = c.db.GetCreateUser(ctxTm, tx, *user)
	if err != nil {
		return err
	}

	return nil
}

func handleAccrual(order *models.Order, user *models.User, service *models.Service) {
	if service.ConfNeeded {
		user.FrozenBalance += order.Amount
		order.Status = consts.OrderInProgress
		order.ProofedAtStr = consts.ProofedNull
	} else {
		user.Balance += order.Amount
		order.Status = consts.OrderSuccess
		order.ProofedAtStr = consts.ProofedNow
	}
}

func handleWithdraw(order *models.Order, user *models.User, service *models.Service) error {
	var err error

	if service.ConfNeeded {
		if user.Balance-order.Amount >= 0 {
			user.FrozenBalance += order.Amount
			user.Balance -= order.Amount
			order.Status = consts.OrderInProgress
			order.ProofedAtStr = consts.ProofedNull
		} else {
			order.Status = consts.OrderDeclined
			order.ProofedAtStr = consts.ProofedNow
			err = fmt.Errorf("not enough money on balance")
		}
	} else {
		if user.Balance-order.Amount >= 0 {
			user.Balance -= order.Amount
			order.Status = consts.OrderSuccess
			order.ProofedAtStr = consts.ProofedNow
		} else {
			order.Status = consts.OrderDeclined
			order.ProofedAtStr = consts.ProofedNow
			err = fmt.Errorf("not enough money on balance")
		}
	}

	return err
}

func (c *Controller) GetHistory(
	ctx context.Context,
	request models.GetHistoryRequest,
) (models.HistoryFrame, models.InternalError) {
	var err models.InternalError
	var history models.HistoryFrame

	history, err.Err = validateGetHistoryParams(request)
	if err.Err != nil {
		err.Type = http.StatusBadRequest
		err.Err = fmt.Errorf("validation error: %w", err.Err)
		return models.HistoryFrame{}, err
	}

	ctxTm, cancel := context.WithTimeout(ctx, c.dbTm)
	defer cancel()

	history, err.Err = c.db.GetHistoryFrame(ctxTm, history)
	if err.Err != nil {
		err.Type = http.StatusInternalServerError
		err.Err = fmt.Errorf("cant get history frame: %w", err.Err)
		return models.HistoryFrame{}, err
	}

	return history, err
}
