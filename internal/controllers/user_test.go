package controllers

import (
	"credit_holidays/internal/consts"
	"credit_holidays/internal/models"
	"reflect"
	"testing"
)

func TestGetBalance(t *testing.T) {

}

func TestInsertUserIfNotExists(t *testing.T) {

}

func TestHandleAccrual(t *testing.T) {
	type args struct {
		order   *models.Order
		user    *models.User
		service *models.Service
	}
	tests := []struct {
		name string
		args args
		res  args
	}{
		{
			name: "correct confirmation needed",
			args: args{&models.Order{Amount: 300}, &models.User{Balance: 400, FrozenBalance: 500}, &models.Service{ConfNeeded: true}},
			res:  args{&models.Order{Amount: 300, ProofedAtStr: consts.ProofedNull, Status: consts.OrderInProgress}, &models.User{Balance: 400, FrozenBalance: 800}, &models.Service{ConfNeeded: true}},
		},
		{
			name: "correct confirmation dont needed",
			args: args{&models.Order{Amount: 300}, &models.User{Balance: 400, FrozenBalance: 500}, &models.Service{ConfNeeded: false}},
			res:  args{&models.Order{Amount: 300, ProofedAtStr: consts.ProofedNow, Status: consts.OrderSuccess}, &models.User{Balance: 700, FrozenBalance: 500}, &models.Service{ConfNeeded: false}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handleAccrual(tt.args.order, tt.args.user, tt.args.service)
			res := reflect.DeepEqual(*tt.args.order, *tt.res.order)
			if !res {
				t.Error(
					"For", *tt.args.order,
					"expected", *tt.res.order,
					"got", tt.args.order,
				)
			}
			res = reflect.DeepEqual(*tt.args.user, *tt.res.user)
			if !res {
				t.Error(
					"For", *tt.args.user,
					"expected", *tt.res.user,
					"got", tt.args.user,
				)
			}
			res = reflect.DeepEqual(*tt.args.service, *tt.res.service)
			if !res {
				t.Error(
					"For", *tt.args.service,
					"expected", *tt.res.service,
					"got", tt.args.service,
				)
			}
		})
	}
}

func TestHandleWithdraw(t *testing.T) {
	type args struct {
		order   *models.Order
		user    *models.User
		service *models.Service
	}
	tests := []struct {
		name string
		args args
		res  args
	}{
		{
			name: "correct confirmation needed",
			args: args{&models.Order{Amount: 300}, &models.User{Balance: 400, FrozenBalance: 500}, &models.Service{ConfNeeded: true}},
			res:  args{&models.Order{Amount: 300, ProofedAtStr: consts.ProofedNull, Status: consts.OrderInProgress}, &models.User{Balance: 100, FrozenBalance: 800}, &models.Service{ConfNeeded: true}},
		},
		{
			name: "incorrect confirmation needed",
			args: args{&models.Order{Amount: 300}, &models.User{Balance: 200, FrozenBalance: 500}, &models.Service{ConfNeeded: true}},
			res:  args{&models.Order{Amount: 300, ProofedAtStr: consts.ProofedNow, Status: consts.OrderDeclined}, &models.User{Balance: 200, FrozenBalance: 500}, &models.Service{ConfNeeded: true}},
		},
		{
			name: "correct confirmation dont needed",
			args: args{&models.Order{Amount: 300}, &models.User{Balance: 400, FrozenBalance: 500}, &models.Service{ConfNeeded: false}},
			res:  args{&models.Order{Amount: 300, ProofedAtStr: consts.ProofedNow, Status: consts.OrderSuccess}, &models.User{Balance: 100, FrozenBalance: 500}, &models.Service{ConfNeeded: false}},
		},
		{
			name: "incorrect confirmation dont needed",
			args: args{&models.Order{Amount: 300}, &models.User{Balance: 200, FrozenBalance: 500}, &models.Service{ConfNeeded: false}},
			res:  args{&models.Order{Amount: 300, ProofedAtStr: consts.ProofedNow, Status: consts.OrderDeclined}, &models.User{Balance: 200, FrozenBalance: 500}, &models.Service{ConfNeeded: false}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handleWithdraw(tt.args.order, tt.args.user, tt.args.service)
			res := reflect.DeepEqual(*tt.args.order, *tt.res.order)
			if !res {
				t.Error(
					"For", *tt.args.order,
					"expected", *tt.res.order,
					"got", tt.args.order,
				)
			}
			res = reflect.DeepEqual(*tt.args.user, *tt.res.user)
			if !res {
				t.Error(
					"For", *tt.args.user,
					"expected", *tt.res.user,
					"got", tt.args.user,
				)
			}
			res = reflect.DeepEqual(*tt.args.service, *tt.res.service)
			if !res {
				t.Error(
					"For", *tt.args.service,
					"expected", *tt.res.service,
					"got", tt.args.service,
				)
			}
		})
	}
}
