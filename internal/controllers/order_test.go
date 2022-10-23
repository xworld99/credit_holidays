package controllers

import (
	"context"
	"credit_holidays/internal/consts"
	"credit_holidays/internal/mocks"
	"credit_holidays/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/knadh/koanf"
	"reflect"
	"testing"
)

func TestAddOrder(t *testing.T) {
	t.Run("correct", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockDb := mocks.NewMockCreditHolidaysDB(mockCtrl)

		ctrl := NewController(koanf.New("."), mockDb)

		mockDb.EXPECT().
			Begin(gomock.Any(), gomock.Any()).
			Times(1)

		mockDb.EXPECT().
			InsertUserIfNotExists(gomock.Any(), gomock.Any(), models.User{Id: 123}).
			Return(models.User{}, nil).
			Times(1)

		mockDb.EXPECT().
			GetServiceById(gomock.Any(), gomock.Any()).
			Return(models.Service{}, nil).
			Times(1)

		mockDb.EXPECT().
			CreateOrder(gomock.Any(), gomock.Any(), models.Order{Amount: 100.0}).
			Return(models.Order{}, nil).
			Times(1)

		mockDb.EXPECT().
			UpdateUser(gomock.Any(), gomock.Any(), models.User{}).
			Return(models.User{}, nil).
			Times(1)

		mockDb.EXPECT().
			UpdateOrder(gomock.Any(), gomock.Any(), models.Order{Status: "success", ProofedAtStr: consts.ProofedNow}).
			Return(models.Order{}, nil).
			Times(1)

		mockDb.EXPECT().
			Commit(gomock.Any()).
			Times(1)

		mockDb.EXPECT().
			Rollback(gomock.Any()).
			Times(1)

		o := models.AddOrderRequest{UserId: 123, ServiceId: 456, Amount: 100}
		_, err := ctrl.AddOrder(context.Background(), o)
		if err.Err != nil {
			t.Error(
				"For", o,
				"expected", "no error",
				"got", err.Error(),
			)
		}
	})

}

func TestChangeOrderStatus(t *testing.T) {
	t.Run("correct proof", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockDb := mocks.NewMockCreditHolidaysDB(mockCtrl)

		ctrl := NewController(koanf.New("."), mockDb)

		mockDb.EXPECT().
			Begin(gomock.Any(), gomock.Any()).
			Times(1)

		mockDb.EXPECT().
			GetFullOrderInfo(gomock.Any(), gomock.Any(), models.Order{Id: 123}, models.User{}, models.Service{}).
			Return(
				models.Order{Id: 123, Amount: 300, Status: consts.OrderInProgress},
				models.User{Id: 1, Balance: 123, FrozenBalance: 300},
				models.Service{ServiceType: consts.OperationWithdraw},
				nil).
			Times(1)

		mockDb.EXPECT().
			UpdateUser(gomock.Any(), gomock.Any(), models.User{Id: 1, Balance: 123, FrozenBalance: 0}).
			Return(models.User{}, nil).
			Times(1)

		mockDb.EXPECT().
			UpdateOrder(
				gomock.Any(),
				gomock.Any(),
				models.Order{Id: 123, Amount: 300, Status: consts.OrderSuccess, ProofedAtStr: consts.ProofedNow}).
			Return(models.Order{}, nil).
			Times(1)

		mockDb.EXPECT().
			Commit(gomock.Any()).
			Times(1)

		mockDb.EXPECT().
			Rollback(gomock.Any()).
			Times(1)

		o := models.ChangeOrderRequest{OrderId: 123, Action: consts.OrderProof}
		_, err := ctrl.ChangeOrderStatus(context.Background(), o)
		if err.Err != nil {
			t.Error(
				"For", o,
				"expected", "no error",
				"got", err.Error(),
			)
		}
	})

	t.Run("correct decline", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockDb := mocks.NewMockCreditHolidaysDB(mockCtrl)

		ctrl := NewController(koanf.New("."), mockDb)

		mockDb.EXPECT().
			Begin(gomock.Any(), gomock.Any()).
			Times(1)

		mockDb.EXPECT().
			GetFullOrderInfo(gomock.Any(), gomock.Any(), models.Order{Id: 123}, models.User{}, models.Service{}).
			Return(
				models.Order{Id: 123, Amount: 300, Status: consts.OrderInProgress},
				models.User{Id: 1, Balance: 123, FrozenBalance: 300},
				models.Service{ServiceType: consts.OperationWithdraw},
				nil).
			Times(1)

		mockDb.EXPECT().
			UpdateUser(gomock.Any(), gomock.Any(), models.User{Id: 1, Balance: 423, FrozenBalance: 0}).
			Return(models.User{}, nil).
			Times(1)

		mockDb.EXPECT().
			UpdateOrder(
				gomock.Any(),
				gomock.Any(),
				models.Order{Id: 123, Amount: 300, Status: consts.OrderDeclined, ProofedAtStr: consts.ProofedNow}).
			Return(models.Order{}, nil).
			Times(1)

		mockDb.EXPECT().
			Commit(gomock.Any()).
			Times(1)

		mockDb.EXPECT().
			Rollback(gomock.Any()).
			Times(1)

		o := models.ChangeOrderRequest{OrderId: 123, Action: consts.OrderDecline}
		_, err := ctrl.ChangeOrderStatus(context.Background(), o)
		if err.Err != nil {
			t.Error(
				"For", o,
				"expected", "no error",
				"got", err.Error(),
			)
		}
	})
}

func TestAcceptOrder(t *testing.T) {
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
			name: "correct accrual",
			args: args{&models.Order{Amount: 300}, &models.User{Balance: 400, FrozenBalance: 500}, &models.Service{ServiceType: consts.OperationAccrual}},
			res:  args{&models.Order{Amount: 300, ProofedAtStr: consts.ProofedNow, Status: consts.OrderSuccess}, &models.User{Balance: 700, FrozenBalance: 200}, &models.Service{ServiceType: consts.OperationAccrual}},
		},
		{
			name: "correct withdraw",
			args: args{&models.Order{Amount: 300}, &models.User{Balance: 400, FrozenBalance: 500}, &models.Service{ServiceType: consts.OperationWithdraw}},
			res:  args{&models.Order{Amount: 300, ProofedAtStr: consts.ProofedNow, Status: consts.OrderSuccess}, &models.User{Balance: 400, FrozenBalance: 200}, &models.Service{ServiceType: consts.OperationWithdraw}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acceptOrder(tt.args.order, tt.args.user, tt.args.service)
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

func TestDeclineOrder(t *testing.T) {
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
			name: "correct accrual",
			args: args{&models.Order{Amount: 300}, &models.User{Balance: 400, FrozenBalance: 500}, &models.Service{ServiceType: consts.OperationAccrual}},
			res:  args{&models.Order{Amount: 300, ProofedAtStr: consts.ProofedNow, Status: consts.OrderDeclined}, &models.User{Balance: 400, FrozenBalance: 200}, &models.Service{ServiceType: consts.OperationAccrual}},
		},
		{
			name: "correct withdraw",
			args: args{&models.Order{Amount: 300}, &models.User{Balance: 400, FrozenBalance: 500}, &models.Service{ServiceType: consts.OperationWithdraw}},
			res:  args{&models.Order{Amount: 300, ProofedAtStr: consts.ProofedNow, Status: consts.OrderDeclined}, &models.User{Balance: 700, FrozenBalance: 200}, &models.Service{ServiceType: consts.OperationWithdraw}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			declineOrder(tt.args.order, tt.args.user, tt.args.service)
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
