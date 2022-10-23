package controllers

import (
	"context"
	"credit_holidays/internal/consts"
	"credit_holidays/internal/mocks"
	"credit_holidays/internal/models"
	"database/sql"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/knadh/koanf"
	"reflect"
	"testing"
	"time"
)

func TestGetBalance(t *testing.T) {
	t.Run("correct", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockDb := mocks.NewMockCreditHolidaysDB(mockCtrl)

		ctrl := NewController(koanf.New("."), mockDb)

		mockDb.EXPECT().
			GetUserById(gomock.Any(), models.User{Id: 44}).
			Return(models.User{Id: 44}, nil).
			Times(1)

		u := models.GetBalanceRequest("44")
		_, err := ctrl.GetBalance(context.Background(), u)
		if err.Err != nil {
			t.Error(
				"For", u,
				"expected", "no error",
				"got", err.Error(),
			)
		}
	})

	t.Run("no user", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockDb := mocks.NewMockCreditHolidaysDB(mockCtrl)

		ctrl := NewController(koanf.New("."), mockDb)

		mockDb.EXPECT().
			GetUserById(gomock.Any(), models.User{Id: 44}).
			Return(models.User{}, fmt.Errorf("no user")).
			Times(1)

		u := models.GetBalanceRequest("44")
		_, err := ctrl.GetBalance(context.Background(), u)
		if err.Err == nil {
			t.Error(
				"For", u,
				"expected", "error",
				"got", "no error",
			)
		}
	})

	t.Run("negative id", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockDb := mocks.NewMockCreditHolidaysDB(mockCtrl)

		ctrl := NewController(koanf.New("."), mockDb)

		u := models.GetBalanceRequest("-44")
		_, err := ctrl.GetBalance(context.Background(), u)
		if err.Err == nil {
			t.Error(
				"For", u,
				"expected", "error",
				"got", "no error",
			)
		}
	})

	t.Run("zero id", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockDb := mocks.NewMockCreditHolidaysDB(mockCtrl)

		ctrl := NewController(koanf.New("."), mockDb)

		u := models.GetBalanceRequest("0")
		_, err := ctrl.GetBalance(context.Background(), u)
		if err.Err == nil {
			t.Error(
				"For", u,
				"expected", "error",
				"got", "no error",
			)
		}
	})

}

func TestInsertUserIfNotExists(t *testing.T) {
	t.Run("user exists", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockDb := mocks.NewMockCreditHolidaysDB(mockCtrl)

		mockDb.EXPECT().
			InsertUserIfNotExists(gomock.Any(), gomock.Any(), models.User{Id: 123}).
			Return(models.User{Id: 123, Balance: 15, FrozenBalance: 15}, nil).
			Times(1)

		ctrl := NewController(koanf.New("."), mockDb)

		u := &models.User{Id: 123}
		err := ctrl.InsertUserIfNotExists(context.Background(), &sql.Tx{}, u)
		if err != nil {
			t.Error(
				"For", *u,
				"expected", "no error",
				"got", err.Error(),
			)
		}
	})
	t.Run("user not exists", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockDb := mocks.NewMockCreditHolidaysDB(mockCtrl)

		mockDb.EXPECT().
			InsertUserIfNotExists(gomock.Any(), gomock.Any(), models.User{Id: 123}).
			Return(models.User{Id: 123, Balance: 0, FrozenBalance: 0}, nil).
			Times(1)

		ctrl := NewController(koanf.New("."), mockDb)

		u := &models.User{Id: 123}
		err := ctrl.InsertUserIfNotExists(context.Background(), &sql.Tx{}, &models.User{Id: 123})
		if err != nil {
			t.Error(
				"For", *u,
				"expected", "no error",
				"got", err.Error(),
			)
		}
	})
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

func TestGetHistory(t *testing.T) {
	t.Run("correct all fields passed 1", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockDb := mocks.NewMockCreditHolidaysDB(mockCtrl)

		ctrl := NewController(koanf.New("."), mockDb)

		mockDb.EXPECT().
			GetHistoryFrame(
				gomock.Any(),
				models.HistoryFrame{
					UserId:   123,
					FromDate: time.Date(2021, time.Month(2), 10, 0, 0, 0, 0, time.UTC),
					ToDate:   time.Date(2021, time.Month(2), 28, 0, 0, 0, 0, time.UTC),
					Offset:   12,
					Limit:    54,
					OrderBy:  "created_at",
				}).
			Return(models.HistoryFrame{}, nil).
			Times(1)

		_, err := ctrl.GetHistory(context.Background(), models.GetHistoryRequest{UserId: "123", FromDate: "10-02-2021", ToDate: "28-02-2021", Offset: "12", Limit: "54", OrderBy: "created_at"})
		if err.Err != nil {
			t.Fail()
		}
	})

	t.Run("correct all fields passed 2", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockDb := mocks.NewMockCreditHolidaysDB(mockCtrl)

		ctrl := NewController(koanf.New("."), mockDb)

		mockDb.EXPECT().
			GetHistoryFrame(
				gomock.Any(),
				models.HistoryFrame{
					UserId:   123,
					FromDate: time.Date(2021, time.Month(2), 10, 0, 0, 0, 0, time.UTC),
					ToDate:   time.Date(2021, time.Month(2), 28, 0, 0, 0, 0, time.UTC),
					Offset:   12,
					Limit:    54,
					OrderBy:  "amount",
				}).
			Return(models.HistoryFrame{}, nil).
			Times(1)

		_, err := ctrl.GetHistory(context.Background(), models.GetHistoryRequest{UserId: "123", FromDate: "10-02-2021", ToDate: "28-02-2021", Offset: "12", Limit: "54", OrderBy: "amount"})
		if err.Err != nil {
			t.Fail()
		}
	})
}
