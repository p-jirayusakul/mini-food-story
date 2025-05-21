package usecase

import (
	"context"
	"errors"
	"food-story/kitchen-service/internal/adapter/repository"
	"food-story/pkg/exceptions"
	database "food-story/shared/database/sqlc"
	"food-story/shared/mock/database"
	"food-story/shared/mock/shared"
	shareModel "food-story/shared/model"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestUpdateOrderItemsStatus(t *testing.T) {

	orderID := int64(1)
	orderItemsID := int64(1)
	testCases := []struct {
		name          string
		payload       shareModel.OrderItemsStatus
		buildStubs    func(db *mockdb.MockStore, payload shareModel.OrderItemsStatus)
		checkResponse func(t *testing.T, customError *exceptions.CustomError)
	}{
		{
			name: "OK",
			payload: shareModel.OrderItemsStatus{
				ID:         orderItemsID,
				OrderID:    orderID,
				StatusCode: "SERVED",
			},
			buildStubs: func(db *mockdb.MockStore, payload shareModel.OrderItemsStatus) {
				db.EXPECT().IsOrderWithItemsExists(gomock.Any(), database.IsOrderWithItemsExistsParams{
					OrderID:      orderID,
					OrderItemsID: orderItemsID,
				}).Times(1).Return(true, nil)

				// update order items status
				db.EXPECT().UpdateOrderItemsStatus(gomock.Any(), database.UpdateOrderItemsStatusParams{
					StatusCode: payload.StatusCode,
					ID:         payload.ID,
				}).Times(1).Return(nil)
			},
			checkResponse: func(t *testing.T, customError *exceptions.CustomError) {
				require.Nil(t, customError)
			},
		},
		{
			name: "Order Items Not Found",
			payload: shareModel.OrderItemsStatus{
				ID:         orderItemsID,
				OrderID:    orderID,
				StatusCode: "SERVED",
			},
			buildStubs: func(db *mockdb.MockStore, payload shareModel.OrderItemsStatus) {
				db.EXPECT().IsOrderWithItemsExists(gomock.Any(), database.IsOrderWithItemsExistsParams{
					OrderID:      orderID,
					OrderItemsID: orderItemsID,
				}).Times(1).Return(false, nil)
			},
			checkResponse: func(t *testing.T, customError *exceptions.CustomError) {
				require.Error(t, customError.Errors)
				require.Equal(t, exceptions.ERRNOTFOUND, customError.Status)
			},
		},
		{
			name: "Internal Server Error",
			payload: shareModel.OrderItemsStatus{
				ID:         orderItemsID,
				OrderID:    orderID,
				StatusCode: "SERVED",
			},
			buildStubs: func(db *mockdb.MockStore, payload shareModel.OrderItemsStatus) {
				db.EXPECT().IsOrderWithItemsExists(gomock.Any(), database.IsOrderWithItemsExistsParams{
					OrderID:      orderID,
					OrderItemsID: orderItemsID,
				}).Times(1).Return(false, errors.New("something error"))
			},
			checkResponse: func(t *testing.T, customError *exceptions.CustomError) {
				require.Error(t, customError.Errors)
				require.Equal(t, exceptions.ERRREPOSITORY, customError.Status)
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

			tc.buildStubs(store, tc.payload)
			repo := repository.NewRepository(cfg, store, node)
			useCase := NewUsecase(cfg, *repo)

			ctx := context.Background()
			customError := useCase.UpdateOrderItemsStatus(ctx, tc.payload)
			tc.checkResponse(t, customError)
		})

	}
}
