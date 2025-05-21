package usecase

import (
	"context"
	"errors"
	"food-story/order-service/internal/adapter/repository"
	"food-story/order-service/internal/domain"
	"food-story/pkg/exceptions"
	"food-story/pkg/utils"
	database "food-story/shared/database/sqlc"
	"food-story/shared/mock/cache/order"
	"food-story/shared/mock/database"
	"food-story/shared/mock/queue/order"
	"food-story/shared/mock/shared"
	shareModel "food-story/shared/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"sort"
	"testing"
	"time"
)

func TestCreateOrder(t *testing.T) {
	cfg := mockshared.MockupConfig()

	sessionID := uuid.New()
	mockupOrderItems := mockshared.MockOrderItems()

	//var oderID *string
	//orderIDString := "1"
	//oderID = &orderIDString

	orderIDInt := int64(1)

	currentLocation, _ := time.LoadLocation(cfg.TimeZone)
	mockTimeNow, _ := time.Parse("2006-01-02", mockupOrderItems[0].CreatedAt)
	currentTime := mockTimeNow.In(currentLocation)

	mockupOrderNumber := "FS-00010101-0001"

	mockupCurrentTableSession := shareModel.CurrentTableSession{
		SessionID:   sessionID,
		TableID:     int64(1),
		TableNumber: int32(1),
		Status:      "active",
		StartedAt:   time.Now(),
		//OrderID:     oderID,
	}

	var mockupData []*database.GetOrderWithItemsRow
	var createOrderTx []database.CreateOrderItemsParams
	for _, item := range mockupOrderItems {

		createOrderTx = append(createOrderTx, database.CreateOrderItemsParams{
			ID:            item.OrderID,
			OrderID:       orderIDInt,
			ProductID:     item.ProductID,
			StatusID:      int64(1),
			ProductName:   item.ProductName,
			ProductNameEn: item.ProductNameEN,
			Price:         utils.Float64ToPgNumeric(item.Price),
			Quantity:      item.Quantity,
			Note:          utils.StringPtrToPgText(item.Note),
			CreatedAt:     pgtype.Timestamptz{Time: mockTimeNow, Valid: true},
		})

		layout := time.RFC3339
		convertedTime, _ := time.Parse(layout, item.CreatedAt)
		mockupData = append(mockupData, &database.GetOrderWithItemsRow{
			OrderID:       item.OrderID,
			OrderNumber:   item.OrderNumber,
			ID:            item.ID,
			ProductID:     item.ProductID,
			ProductName:   item.ProductName,
			ProductNameEN: item.ProductNameEN,
			Quantity:      item.Quantity,
			Price:         utils.Float64ToPgNumeric(item.Price),
			StatusID:      item.StatusID,
			StatusName:    item.StatusName,
			StatusNameEN:  item.StatusNameEN,
			StatusCode:    item.StatusCode,
			Note:          utils.StringPtrToPgText(item.Note),
			CreatedAt:     pgtype.Timestamptz{Time: convertedTime, Valid: true},
			TableNumber:   item.TableNumber,
		})
	}

	testCases := []struct {
		name          string
		buildStubs    func(db *mockdb.MockStore, cache *mockcache.MockRedisTableCacheInterface, queue *mockqueue.MockQueueProducerInterface, snowflake *mockshared.MockSnowflakeInterface)
		checkResponse func(t *testing.T, result int64, customError *exceptions.CustomError)
	}{
		{
			name: "OK",
			buildStubs: func(db *mockdb.MockStore, cache *mockcache.MockRedisTableCacheInterface, queue *mockqueue.MockQueueProducerInterface, snowflake *mockshared.MockSnowflakeInterface) {
				cache.EXPECT().GetCachedTable(sessionID).AnyTimes().Return(&mockupCurrentTableSession, nil)

				db.EXPECT().GetOrderStatusPreparing(gomock.Any()).Times(1).Return(int64(1), nil)
				db.EXPECT().GetTimeNow(gomock.Any()).AnyTimes().Return(pgtype.Timestamptz{Time: mockTimeNow, Valid: true}, nil)
				db.EXPECT().GetOrderSequence(gomock.Any(), pgtype.Date{Time: currentTime, Valid: true}).Times(1).Return(int32(1), nil)

				for _, v := range mockupData {
					db.EXPECT().GetProductByID(gomock.Any(), v.ProductID).Times(1).Return(&database.GetProductByIDRow{
						ID:             int64(1),
						Name:           v.ProductName,
						NameEn:         v.ProductNameEN,
						Categories:     int64(1),
						CategoryName:   "ขนม",
						CategoryNameEN: "Meat",
						Description:    pgtype.Text{String: "description", Valid: true},
						Price:          v.Price,
						IsAvailable:    true,
						//ImageUrl
					}, nil)
				}

				snowflake.EXPECT().Generate().AnyTimes().Return(orderIDInt)
				db.EXPECT().TXCreateOrder(gomock.Any(), database.TXCreateOrderParams{
					CreateOrderItems: createOrderTx,
					CreateOrder: database.CreateOrderParams{
						ID:          orderIDInt,
						OrderNumber: mockupOrderNumber,
						SessionID:   utils.UUIDToPgUUID(sessionID),
						TableID:     int64(1),
					},
				}).Times(1).Return(orderIDInt, nil)

				cache.EXPECT().UpdateOrderID(sessionID, orderIDInt).Times(1).Return(nil)

				db.EXPECT().IsOrderExist(gomock.Any(), orderIDInt).Times(1).Return(true, nil)
				db.EXPECT().GetOrderWithItems(gomock.Any(), orderIDInt).Times(1).Return(mockupData, nil)

				var orderItems = mockshared.MockOrderItemsPt()
				sort.Slice(orderItems, func(i, j int) bool {
					return orderItems[i].ID < orderItems[j].ID
				})

				for _, v := range orderItems {
					queue.EXPECT().PublishOrder(*v).Times(1).Return(nil)
				}

			},
			checkResponse: func(t *testing.T, result int64, customError *exceptions.CustomError) {
				require.Nil(t, customError)
				require.Equal(t, orderIDInt, result)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			redis := mockcache.NewMockRedisTableCacheInterface(ctrl)
			node := mockshared.NewMockSnowflakeInterface(ctrl)
			queue := mockqueue.NewMockQueueProducerInterface(ctrl)

			tc.buildStubs(store, redis, queue, node)
			repo := repository.NewRepository(cfg, store, node)
			useCase := NewUsecase(cfg, *repo, redis, queue)

			ctx := context.Background()
			result, customError := useCase.CreateOrder(ctx, sessionID, mockupOrderItems)
			tc.checkResponse(t, result, customError)
		})

	}
}

func TestGetOrderByID(t *testing.T) {

	sessionID := uuid.New()

	var oderID *string
	orderIDString := "1"
	orderIDInt := int64(1)
	oderID = &orderIDString

	mockupCurrentTableSession := shareModel.CurrentTableSession{
		SessionID:   sessionID,
		TableID:     int64(1),
		TableNumber: int32(1),
		Status:      "active",
		StartedAt:   time.Now(),
		OrderID:     oderID,
	}

	mockupData := database.GetOrderByIDRow{
		ID:           orderIDInt,
		SessionID:    utils.UUIDToPgUUID(sessionID),
		TableID:      mockupCurrentTableSession.TableID,
		TableNumber:  mockupCurrentTableSession.TableNumber,
		StatusID:     int64(1),
		StatusName:   "กำหลังเตรียม",
		StatusNameEN: "Preparing",
		StatusCode:   "PREPARING",
	}

	mockupResult := domain.Order{
		ID:           orderIDInt,
		TableID:      mockupData.TableID,
		TableNumber:  mockupData.TableNumber,
		StatusID:     mockupData.StatusID,
		StatusName:   mockupData.StatusName,
		StatusNameEN: mockupData.StatusNameEN,
		StatusCode:   mockupData.StatusCode,
	}

	testCases := []struct {
		name          string
		buildStubs    func(db *mockdb.MockStore, cache *mockcache.MockRedisTableCacheInterface)
		checkResponse func(t *testing.T, result *domain.Order, customError *exceptions.CustomError)
	}{
		{
			name: "OK",
			buildStubs: func(db *mockdb.MockStore, cache *mockcache.MockRedisTableCacheInterface) {
				cache.EXPECT().GetCachedTable(sessionID).Times(1).Return(&mockupCurrentTableSession, nil)
				db.EXPECT().GetOrderByID(gomock.Any(), orderIDInt).Times(1).Return(&mockupData, nil)
			},
			checkResponse: func(t *testing.T, result *domain.Order, customError *exceptions.CustomError) {
				require.Nil(t, customError)
				require.Equal(t, &mockupResult, result)
			},
		},
		{
			name: "Internal Server Error",
			buildStubs: func(db *mockdb.MockStore, cache *mockcache.MockRedisTableCacheInterface) {
				cache.EXPECT().GetCachedTable(sessionID).Times(1).Return(&mockupCurrentTableSession, nil)
				db.EXPECT().GetOrderByID(gomock.Any(), orderIDInt).Times(1).Return(nil, errors.New("internal server error"))
			},
			checkResponse: func(t *testing.T, result *domain.Order, customError *exceptions.CustomError) {
				require.Error(t, customError.Errors)
				require.Equal(t, exceptions.ERRREPOSITORY, customError.Status)
				require.Nil(t, result)
			},
		},
		{
			name: "Data Not Found",
			buildStubs: func(db *mockdb.MockStore, cache *mockcache.MockRedisTableCacheInterface) {
				cache.EXPECT().GetCachedTable(sessionID).Times(1).Return(&mockupCurrentTableSession, nil)
				db.EXPECT().GetOrderByID(gomock.Any(), orderIDInt).Times(1).Return(nil, exceptions.ErrRowDatabaseNotFound)
			},
			checkResponse: func(t *testing.T, result *domain.Order, customError *exceptions.CustomError) {
				require.Error(t, customError.Errors)
				require.Equal(t, exceptions.ERRNOTFOUND, customError.Status)
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
			redis := mockcache.NewMockRedisTableCacheInterface(ctrl)
			node := mockshared.NewMockSnowflakeInterface(ctrl)
			queue := mockqueue.NewMockQueueProducerInterface(ctrl)

			tc.buildStubs(store, redis)
			repo := repository.NewRepository(cfg, store, node)
			useCase := NewUsecase(cfg, *repo, redis, queue)

			ctx := context.Background()
			result, customError := useCase.GetOrderByID(ctx, sessionID)
			tc.checkResponse(t, result, customError)
		})

	}
}
