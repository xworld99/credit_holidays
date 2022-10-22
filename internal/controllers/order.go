package controllers

import (
	"context"
	"credit_holidays/internal/consts"
	"credit_holidays/internal/models"
	"database/sql"
	"fmt"
	"net/http"
)

// logic for working with orders

func (c *Controller) AddOrder(
	ctx context.Context,
	orderInfo models.AddOrderRequest,
) (models.Order, models.InternalError) {
	var err models.InternalError

	// validate params
	err.Err = validateOrderParams(orderInfo)
	if err.Err != nil {
		err.Type = http.StatusBadRequest
		err.Err = fmt.Errorf("validation error: %w", err.Err)
		return models.Order{}, err
	}

	user, service, order := models.User{Id: orderInfo.UserId}, models.Service{Id: orderInfo.ServiceId}, models.Order{}

	// get service
	err.Err = c.getServiceInfo(ctx, &service)
	if err.Err != nil {
		err.Type = http.StatusNotFound
		err.Err = fmt.Errorf("cant get service %d: %w", service.Id, err.Err)
		return models.Order{}, err
	}

	// begin transaction
	ctxTm, cancel := context.WithTimeout(ctx, c.dbTm)
	defer cancel()

	var tx *sql.Tx
	tx, err.Err = c.db.Begin(ctxTm, sql.LevelSerializable)
	if err.Err != nil {
		err.Type = http.StatusInternalServerError
		err.Err = fmt.Errorf("cant create transaction: %w", err.Err)
		return models.Order{}, err
	}
	defer tx.Rollback()

	// get current user balance
	err.Err = c.getCreateUser(ctxTm, tx, &user)
	if err.Err != nil {
		err.Type = http.StatusInternalServerError
		err.Err = fmt.Errorf("cant get user %d: %w", user.Id, err.Err)
		return models.Order{}, err
	}

	// create new order
	order.UserId = user.Id
	order.ServiceId = service.Id
	order.Amount = orderInfo.Amount
	order, err.Err = c.db.CreateOrder(ctxTm, tx, order)
	if err.Err != nil {
		err.Type = http.StatusInternalServerError
		err.Err = fmt.Errorf("cant create order: %w", err.Err)
		return models.Order{}, err
	}

	// change user balance according to order
	if service.ServiceType == consts.OperationAccrual {
		handleAccrual(&order, &user, &service)
	} else {
		err.Err = handleWithdraw(&order, &user, &service)
	}
	if err.Err != nil {
		err.Type = http.StatusInternalServerError
		err.Err = fmt.Errorf("cant perform %s: %w", service.ServiceType, err.Err)
		return models.Order{}, err
	}

	// update user
	user, err.Err = c.db.UpdateUser(ctxTm, tx, user)
	if err.Err != nil {
		err.Type = http.StatusInternalServerError
		err.Err = fmt.Errorf("cant update user %d: %w", user.Id, err.Err)
		return models.Order{}, err
	}

	// update order
	order, err.Err = c.db.UpdateOrder(ctxTm, tx, order)
	if err.Err != nil {
		err.Type = http.StatusInternalServerError
		err.Err = fmt.Errorf("cant update order: %w", err.Err)
		return models.Order{}, err
	}

	// commit transaction
	tx.Commit()

	return order, err
}

func acceptOrder(order *models.Order, user *models.User, service *models.Service) {
	if service.ServiceType == consts.OperationAccrual {
		user.Balance += order.Amount
		user.FrozenBalance -= order.Amount
	} else {
		user.FrozenBalance -= order.Amount
	}

	order.ProofedAtStr = consts.ProofedNow
	order.Status = consts.OrderSuccess
}

func declineOrder(order *models.Order, user *models.User, service *models.Service) {
	if service.ServiceType == consts.OperationAccrual {
		user.FrozenBalance -= order.Amount
	} else {
		user.Balance += order.Amount
		user.FrozenBalance -= order.Amount
	}

	order.ProofedAtStr = consts.ProofedNow
	order.Status = consts.OrderSuccess
}

func (c *Controller) ChangeOrderStatus(
	ctx context.Context,
	orderInfo models.ChangeOrderRequest,
) (models.Order, models.InternalError) {
	var err models.InternalError

	// validate params
	err.Err = validateChangeStatusParams(orderInfo)
	if err.Err != nil {
		err.Type = http.StatusBadRequest
		err.Err = fmt.Errorf("validation error: %w", err.Err)
		return models.Order{}, err
	}

	// init transaction
	ctxTm, cancel := context.WithTimeout(ctx, c.dbTm)
	defer cancel()

	var tx *sql.Tx
	tx, err.Err = c.db.Begin(ctxTm, sql.LevelSerializable)
	if err.Err != nil {
		err.Type = http.StatusInternalServerError
		err.Err = fmt.Errorf("cant create transaction: %w", err.Err)
		return models.Order{}, err
	}
	defer tx.Rollback()

	// get full info about order, user and service by order_id
	order, user, service := models.Order{Id: orderInfo.OrderId}, models.User{}, models.Service{}
	order, user, service, err.Err = c.db.GetFullOrderInfo(ctx, tx, order, user, service)
	if err.Err != nil {
		err.Type = http.StatusInternalServerError
		err.Err = fmt.Errorf("cant get full order info %d: %w", order.Id, err.Err)
		return models.Order{}, err
	}

	if order.Status != consts.OrderInProgress {
		err.Err = fmt.Errorf("order with id %d is already ended with status %s", order.Id, order.Status)
		err.Type = http.StatusBadRequest
		return models.Order{}, err
	}

	// manipulate with balance
	if orderInfo.Action == consts.OrderProof {
		acceptOrder(&order, &user, &service)
	} else {
		declineOrder(&order, &user, &service)
	}

	// update user
	user, err.Err = c.db.UpdateUser(ctxTm, tx, user)
	if err.Err != nil {
		err.Type = http.StatusInternalServerError
		err.Err = fmt.Errorf("cant update user %d: %w", user.Id, err.Err)
		return models.Order{}, err
	}

	// update order
	order, err.Err = c.db.UpdateOrder(ctxTm, tx, order)
	if err.Err != nil {
		err.Type = http.StatusInternalServerError
		err.Err = fmt.Errorf("cant update order %d: %w", order.Id, err.Err)
		return models.Order{}, err
	}

	// commit transaction
	tx.Commit()

	return order, err
}
