package controllers

import (
	"context"
	"credit_holidays/internal/consts"
	"credit_holidays/internal/models"
	"database/sql"
	"fmt"
)

func (c *Controller) GetBalance(
	ctx context.Context,
	userId models.GetBalanceRequest,
) (models.User, models.HandlerError) {
	var err error

	var id int64
	id, err = validateId(string(userId))
	if err != nil {
		return models.User{}, models.CreateBadRequestError(fmt.Errorf("validation error: %w", err))
	}

	ctxTm, cancel := context.WithTimeout(ctx, c.dbTm)
	defer cancel()

	user := models.User{Id: id}
	user, err = c.db.GetUserById(ctxTm, user)
	if err != nil {
		return models.User{}, models.CreateNotFoundError(fmt.Errorf("cant get user by id: %w", err))
	}

	return user, models.HandlerError{}
}

func (c *Controller) InsertUserIfNotExists(ctx context.Context, tx *sql.Tx, user *models.User) error {
	ctxTm, cancel := context.WithTimeout(ctx, c.dbTm)
	defer cancel()

	var err error
	*user, err = c.db.InsertUserIfNotExists(ctxTm, tx, *user)
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
) (models.HistoryFrame, models.HandlerError) {
	var err error
	var history models.HistoryFrame

	history, err = validateGetHistoryParams(request)
	if err != nil {
		return models.HistoryFrame{}, models.CreateBadRequestError(fmt.Errorf("validation error: %w", err))
	}

	ctxTm, cancel := context.WithTimeout(ctx, c.dbTm)
	defer cancel()

	history, err = c.db.GetHistoryFrame(ctxTm, history)
	if err != nil {
		return models.HistoryFrame{}, models.CreateNotFoundError(fmt.Errorf("cant get history frame: %w", err))
	}

	return history, models.HandlerError{}
}
