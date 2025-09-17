package usecase

import (
	"context"
	"errors"
	"food-story/payment-service/internal/adapter/repository"
	"food-story/payment-service/internal/domain"
	"food-story/pkg/exceptions"
	database "food-story/shared/database/sqlc"
	mockdb "food-story/shared/mock/database"
	mockshared "food-story/shared/mock/shared"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestListPaymentMethods(t *testing.T) {

	mockupResult := []*domain.PaymentMethod{
		{
			ID:   1,
			Name: "name",
			Code: "code",
		},
		{
			ID:   2,
			Name: "name2",
			Code: "code2",
		},
	}

	mockupData := []*database.ListPaymentMethodsRow{
		{
			ID:   1,
			Name: "name",
			Code: "code",
		},
		{
			ID:   2,
			Name: "name2",
			Code: "code2",
		},
	}

	testCases := []struct {
		name          string
		buildStubs    func(db *mockdb.MockStore)
		checkResponse func(t *testing.T, result []*domain.PaymentMethod, customError *exceptions.CustomError)
	}{
		{
			name: "OK",
			buildStubs: func(db *mockdb.MockStore) {
				db.EXPECT().ListPaymentMethods(gomock.Any()).Times(1).Return(mockupData, nil)
			},
			checkResponse: func(t *testing.T, result []*domain.PaymentMethod, customError *exceptions.CustomError) {
				require.Nil(t, customError)
				require.Equal(t, mockupResult, result)
			},
		},
		{
			name: "Internal Server Error",
			buildStubs: func(db *mockdb.MockStore) {
				db.EXPECT().ListPaymentMethods(gomock.Any()).Times(1).Return(nil, errors.New("something error"))
			},
			checkResponse: func(t *testing.T, result []*domain.PaymentMethod, customError *exceptions.CustomError) {
				require.Error(t, customError.Errors)
				require.Equal(t, exceptions.ERRREPOSITORY, customError.Status)
				require.Nil(t, result)
			},
		},
		{
			name: "Data Not Found",
			buildStubs: func(db *mockdb.MockStore) {
				db.EXPECT().ListPaymentMethods(gomock.Any()).Times(1).Return(nil, nil)
			},
			checkResponse: func(t *testing.T, result []*domain.PaymentMethod, customError *exceptions.CustomError) {
				require.Error(t, customError.Errors)
				require.Equal(t, exceptions.ERRREPOSITORY, customError.Status)
				require.Nil(t, result)
			},
		},
	}

	cfg := mockshared.MockupConfig()
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			node := mockshared.NewMockSnowflakeInterface(ctrl)

			tc.buildStubs(store)
			repo := repository.NewRepository(cfg, store, node)
			useCase := NewUsecase(cfg, *repo)

			ctx := context.Background()
			result, customError := useCase.ListPaymentMethods(ctx)
			tc.checkResponse(t, result, customError)
		})

	}
}
