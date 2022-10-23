package controllers

import (
	"context"
	"credit_holidays/internal/consts"
	"credit_holidays/internal/models"
	"database/sql"
	"fmt"
)

// logic for working with orders

func (c *Controller) AddOrder(
	ctx context.Context,
	orderInfo models.AddOrderRequest,
) (models.Order, models.HandlerError) {
	var err error

	// validate params
	err = validateOrderParams(orderInfo)
	if err != nil {
		return models.Order{}, models.CreateBadRequestError(fmt.Errorf("validation error: %w", err))
	}

	user, service, order := models.User{Id: orderInfo.UserId}, models.Service{Id: orderInfo.ServiceId}, models.Order{}

	// get service
	err = c.GetServiceInfo(ctx, &service)
	if err != nil {
		return models.Order{}, models.CreateNotFoundError(fmt.Errorf("cant get service %d: %w", service.Id, err))
	}

	// begin transaction
	ctxTm, cancel := context.WithTimeout(ctx, c.dbTm)
	defer cancel()

	var tx *sql.Tx
	tx, err = c.db.Begin(ctxTm, sql.LevelSerializable)
	if err != nil {
		return models.Order{}, models.CreateInternalError(fmt.Errorf("cant Create transaction: %w", err))
	}
	defer c.db.Rollback(tx)

	// get current user balance
	err = c.InsertUserIfNotExists(ctxTm, tx, &user)
	if err != nil {
		return models.Order{}, models.CreateInternalError(fmt.Errorf("cant get user %d: %w", user.Id, err))
	}

	// Create new order
	order.UserId = user.Id
	order.ServiceId = service.Id
	order.Amount = orderInfo.Amount
	order, err = c.db.CreateOrder(ctxTm, tx, order)
	if err != nil {
		return models.Order{}, models.CreateInternalError(fmt.Errorf("cant Create order: %w", err))
	}

	// change user balance according to order
	if service.ServiceType == consts.OperationAccrual {
		handleAccrual(&order, &user, &service)
	} else {
		err = handleWithdraw(&order, &user, &service)
	}
	if err != nil {
		return models.Order{}, models.CreateBadRequestError(fmt.Errorf("cant perform %s: %w", service.ServiceType, err))
	}

	// update user
	user, err = c.db.UpdateUser(ctxTm, tx, user)
	if err != nil {
		return models.Order{}, models.CreateInternalError(fmt.Errorf("cant update user %d: %w", user.Id, err))
	}

	// update order
	order, err = c.db.UpdateOrder(ctxTm, tx, order)
	if err != nil {
		return models.Order{}, models.CreateInternalError(fmt.Errorf("cant update order: %w", err))
	}

	// commit transaction
	c.db.Commit(tx)

	return order, models.HandlerError{}
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
	order.Status = consts.OrderDeclined
}

func (c *Controller) ChangeOrderStatus(
	ctx context.Context,
	orderInfo models.ChangeOrderRequest,
) (models.Order, models.HandlerError) {
	var err error

	// validate params
	err = validateChangeStatusParams(orderInfo)
	if err != nil {
		return models.Order{}, models.CreateBadRequestError(fmt.Errorf("validation error: %w", err))
	}

	// init transaction
	ctxTm, cancel := context.WithTimeout(ctx, c.dbTm)
	defer cancel()

	var tx *sql.Tx
	tx, err = c.db.Begin(ctxTm, sql.LevelSerializable)
	if err != nil {
		return models.Order{}, models.CreateInternalError(fmt.Errorf("cant Create transaction: %w", err))
	}
	defer c.db.Rollback(tx)

	// get full info about order, user and service by order_id
	order, user, service := models.Order{Id: orderInfo.OrderId}, models.User{}, models.Service{}
	order, user, service, err = c.db.GetFullOrderInfo(ctx, tx, order, user, service)
	if err != nil {
		return models.Order{}, models.CreateInternalError(fmt.Errorf("cant get full order info %d: %w", order.Id, err))
	}

	if order.Status != consts.OrderInProgress {
		return models.Order{}, models.CreateBadRequestError(fmt.Errorf(
			"order with id %d is already ended with status %s", order.Id, order.Status))
	}

	// manipulate with balance
	if orderInfo.Action == consts.OrderProof {
		acceptOrder(&order, &user, &service)
	} else {
		declineOrder(&order, &user, &service)
	}

	// update user
	user, err = c.db.UpdateUser(ctxTm, tx, user)
	if err != nil {
		return models.Order{}, models.CreateInternalError(fmt.Errorf("cant update user %d: %w", user.Id, err))
	}

	// update order
	order, err = c.db.UpdateOrder(ctxTm, tx, order)
	if err != nil {
		return models.Order{}, models.CreateInternalError(fmt.Errorf("cant update order %d: %w", order.Id, err))
	}

	// commit transaction
	c.db.Commit(tx)

	return order, models.HandlerError{}
}
