package controllers

import (
	"credit_holidays/internal/consts"
	"fmt"
	"strconv"
)

func validateId(id string) (int64, error) {
	valId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("id is an invalid int")
	}
	if valId <= 0 {
		return 0, fmt.Errorf("id should be positive")
	}
	return valId, nil
}

func validateMoney(amount string) (float64, error) {
	valAmount, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return 0, fmt.Errorf("amount is an invalid float")
	}
	if valAmount <= 0 {
		return 0, fmt.Errorf("negative amount isnt allowed")
	}
	return valAmount, nil
}

func validateOrderParams(userId, serviceId, amount string) (int64, int64, float64, error) {
	var uid, sid int64
	var err error
	var am float64

	uid, err = validateId(userId)
	if err != nil {
		return 0, 0, 0, err
	}
	sid, err = validateId(serviceId)
	if err != nil {
		return 0, 0, 0, err
	}
	am, err = validateMoney(amount)
	if err != nil {
		return 0, 0, 0, err
	}

	return uid, sid, am, nil
}

func validateChangeStatusParams(orderId, action string) (int64, string, error) {
	var oid int64
	var err error

	oid, err = validateId(orderId)
	if err != nil {
		return 0, "", err
	}

	if _, ok := consts.OrderActions[action]; !ok {
		return 0, "", fmt.Errorf("unknown action: %s", action)
	}

	return oid, action, nil

}
