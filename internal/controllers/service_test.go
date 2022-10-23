package controllers

import (
	"context"
	"credit_holidays/internal/mocks"
	"credit_holidays/internal/models"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/knadh/koanf"
	"testing"
)

func TestGetServiceInfo(t *testing.T) {
	t.Run("correct", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockDb := mocks.NewMockCreditHolidaysDB(mockCtrl)

		ctrl := NewController(koanf.New("."), mockDb)

		mockDb.EXPECT().
			GetServiceById(gomock.Any(), models.Service{Id: 123}).
			Return(models.Service{Id: 123}, nil).
			Times(1)

		err := ctrl.getServiceInfo(context.Background(), &models.Service{Id: 123})
		if err != nil {
			t.Fail()
		}
	})

	t.Run("no service", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockDb := mocks.NewMockCreditHolidaysDB(mockCtrl)

		ctrl := NewController(koanf.New("."), mockDb)

		mockDb.EXPECT().
			GetServiceById(gomock.Any(), models.Service{Id: 123}).
			Return(models.Service{}, fmt.Errorf("no such service")).
			Times(1)

		err := ctrl.getServiceInfo(context.Background(), &models.Service{Id: 123})
		if err == nil {
			t.Fail()
		}
	})
}

func TestGetServicesList(t *testing.T) {
	t.Run("correct", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockDb := mocks.NewMockCreditHolidaysDB(mockCtrl)

		ctrl := NewController(koanf.New("."), mockDb)

		mockDb.EXPECT().
			GetServicesList(gomock.Any()).
			Return([]models.Service{{Id: 1}, {Id: 2}}, nil).
			Times(1)

		_, err := ctrl.GetServicesList(context.Background())
		if err.Err != nil {
			t.Fail()
		}
	})

	t.Run("in correct", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockDb := mocks.NewMockCreditHolidaysDB(mockCtrl)

		ctrl := NewController(koanf.New("."), mockDb)

		mockDb.EXPECT().
			GetServicesList(gomock.Any()).
			Return([]models.Service{}, fmt.Errorf("no services")).
			Times(1)

		_, err := ctrl.GetServicesList(context.Background())
		if err.Err == nil {
			t.Fail()
		}
	})
}
