package controllers

import (
	"context"
	"credit_holidays/internal/models"
	"database/sql"
	"fmt"
)

// logic for working with orders

func handleAccrual(order models.Order, user models.User, service models.Service) (models.Order, models.User, error) {
	if service.ConfNeeded {
		user.FrozenBalance += order.Amount
		order.Status = "in_progress"
		order.ProofedAtStr = "null"
	} else {
		user.Balance += order.Amount
		order.Status = "success"
		order.ProofedAtStr = "now()"
	}

	return order, user, nil
}

func handleWithdraw(order models.Order, user models.User, service models.Service) (models.Order, models.User, error) {
	var err error

	if service.ConfNeeded {
		if user.Balance-order.Amount >= 0 {
			user.FrozenBalance += order.Amount
			user.Balance -= order.Amount
			order.Status = "in_progress"
			order.ProofedAtStr = "null"
		} else {
			order.Status = "declined"
			order.ProofedAtStr = "now()"
			err = fmt.Errorf("not enough money on balance")
		}
	} else {
		if user.Balance-order.Amount >= 0 {
			user.Balance -= order.Amount
			order.Status = "success"
			order.ProofedAtStr = "now()"
		} else {
			order.Status = "declined"
			order.ProofedAtStr = "now()"
			err = fmt.Errorf("not enough money on balance")
		}
	}

	return order, user, err
}

func (c *Controller) AddOrder(ctx context.Context, userId, serviceId, amount string) (models.Order, error) {
	var err error

	// validate params
	uid, sid, am, err := validateOrderParams(userId, serviceId, amount)
	if err != nil {
		return models.Order{}, err
	}

	user, service, order := models.User{Id: uid}, models.Service{Id: sid}, models.Order{}

	// get service
	service, err = c.getServiceInfo(ctx, service)
	if err != nil {
		return models.Order{}, err
	}

	// begin transaction
	ctxTm, cancel := context.WithTimeout(ctx, c.insertTm)
	defer cancel()

	tx, err := c.db.Begin(ctxTm, sql.LevelSerializable)
	if err != nil {
		return models.Order{}, fmt.Errorf("cant init transaction: %w", err)
	}
	defer c.db.RollbackTransaction(tx)

	// get current user balance
	user, err = c.getCreateUser(ctxTm, tx, user)
	if err != nil {
		return models.Order{}, err
	}

	// create new order
	order.UserId = user.Id
	order.ServiceId = service.Id
	order.Amount = am
	order, err = c.db.CreateOrder(ctxTm, tx, order)
	if err != nil {
		return models.Order{}, err
	}

	// change user balance according to order
	if service.ServiceType == "accrual" {
		order, user, err = handleAccrual(order, user, service)
	} else {
		order, user, err = handleWithdraw(order, user, service)
	}
	if err != nil {
		return models.Order{}, err
	}

	// update user
	user, err = c.db.UpdateUser(ctxTm, tx, user)
	if err != nil {
		return models.Order{}, err
	}

	// update order
	order, err = c.db.UpdateOrder(ctxTm, tx, order)
	if err != nil {
		return models.Order{}, err
	}

	// commit transaction
	c.db.CommitTransaction(tx)

	return order, nil
}

func acceptOrder(order models.Order, user models.User, service models.Service) (models.Order, models.User, error) {
	if service.ServiceType == "accrual" {
		user.Balance += order.Amount
		user.FrozenBalance -= order.Amount
	} else {
		user.FrozenBalance -= order.Amount
	}

	order.ProofedAtStr = "now()"
	order.Status = "success"

	return order, user, nil
}

func declineOrder(order models.Order, user models.User, service models.Service) (models.Order, models.User, error) {
	if service.ServiceType == "accrual" {
		user.FrozenBalance -= order.Amount
	} else {
		user.Balance += order.Amount
		user.FrozenBalance -= order.Amount
	}

	order.ProofedAtStr = "now()"
	order.Status = "success"

	return order, user, nil

}

func (c *Controller) ChangeOrderStatus(ctx context.Context, orderId, action string) (models.Order, error) {
	var err error

	// validate params
	oid, action, err := validateChangeStatusParams(orderId, action)
	if err != nil {
		return models.Order{}, err
	}

	// init transaction
	ctxTm, cancel := context.WithTimeout(ctx, c.insertTm)
	defer cancel()

	tx, err := c.db.Begin(ctxTm, sql.LevelSerializable)
	if err != nil {
		return models.Order{}, fmt.Errorf("cant init transaction: %w", err)
	}
	defer c.db.RollbackTransaction(tx)

	// get full info about order, user and service by order_id
	order, user, service := models.Order{Id: oid}, models.User{}, models.Service{}
	order, user, service, err = c.db.GetFullOrderInfo(ctx, tx, order, user, service)
	if err != nil {
		return models.Order{}, err
	}

	if order.Status != "in_progress" {
		return models.Order{}, fmt.Errorf("order with id %d is already ended with status %s", order.Id, order.Status)
	}

	// manipulate with balance
	if action == "proof" {
		order, user, err = acceptOrder(order, user, service)
	} else {
		order, user, err = declineOrder(order, user, service)
	}
	if err != nil {
		return models.Order{}, err
	}

	// update user
	user, err = c.db.UpdateUser(ctxTm, tx, user)
	if err != nil {
		return models.Order{}, err
	}

	// update order
	order, err = c.db.UpdateOrder(ctxTm, tx, order)
	if err != nil {
		return models.Order{}, err
	}

	// commit transaction
	c.db.CommitTransaction(tx)

	return order, nil
}
