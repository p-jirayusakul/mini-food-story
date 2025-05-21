package usecase

import (
	"context"
	"errors"
	"food-story/menu-service/internal/adapter/repository"
	"food-story/menu-service/internal/domain"
	"food-story/pkg/exceptions"
	"food-story/pkg/utils"
	"food-story/shared/config"
	database "food-story/shared/database/sqlc"
	"food-story/shared/mock/cache/menu"
	"food-story/shared/mock/database"
	"food-story/shared/mock/shared"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestListTask(t *testing.T) {

	mockupData := []*database.ListCategoryRow{
		{
			ID:     1,
			Name:   "name",
			NameEN: "nameEN",
		},
		{
			ID:     2,
			Name:   "name2",
			NameEN: "nameEN2",
		},
	}

	mockupResult := []*domain.Category{
		{
			ID:     1,
			Name:   "name",
			NameEn: "nameEN",
		},
		{
			ID:     2,
			Name:   "name2",
			NameEn: "nameEN2",
		},
	}

	testCases := []struct {
		name          string
		buildStubs    func(db *mockdb.MockStore)
		checkResponse func(t *testing.T, result []*domain.Category, customError *exceptions.CustomError)
	}{
		{
			name: "OK",
			buildStubs: func(db *mockdb.MockStore) {
				db.EXPECT().ListCategory(gomock.Any()).Times(1).Return(mockupData, nil)
			},
			checkResponse: func(t *testing.T, result []*domain.Category, customError *exceptions.CustomError) {
				require.Nil(t, customError)
				require.Equal(t, mockupResult, result)
			},
		},
		{
			name: "Internal Server Error",
			buildStubs: func(db *mockdb.MockStore) {
				db.EXPECT().ListCategory(gomock.Any()).Times(1).Return(nil, errors.New("something error"))
			},
			checkResponse: func(t *testing.T, result []*domain.Category, customError *exceptions.CustomError) {
				require.Error(t, customError.Errors)
				require.Equal(t, exceptions.ERRREPOSITORY, customError.Status)
				require.Nil(t, result)
			},
		},
		{
			name: "Data Not Found",
			buildStubs: func(db *mockdb.MockStore) {
				db.EXPECT().ListCategory(gomock.Any()).Times(1).Return(nil, nil)
			},
			checkResponse: func(t *testing.T, result []*domain.Category, customError *exceptions.CustomError) {
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
			redis := mockcache.NewMockRedisTableCacheInterface(ctrl)
			node := mockshared.NewMockSnowflakeInterface(ctrl)

			tc.buildStubs(store)
			repo := repository.NewRepository(cfg, store, node)
			useCase := NewUsecase(cfg, *repo, redis)

			ctx := context.Background()
			result, customError := useCase.ListCategory(ctx)
			tc.checkResponse(t, result, customError)
		})

	}
}

func TestGetProductByID(t *testing.T) {
	mockupID := int64(1)
	mockStr := pgtype.Text{String: "xx", Valid: true}
	mockNumber := utils.Float64ToPgNumeric(100.00)
	mockupData := &database.GetProductAvailableByIDRow{
		ID:             mockupID,
		Name:           "ทดสอบ",
		NameEn:         "test",
		Categories:     mockupID,
		CategoryName:   "ทดสอบหมวดหมู่",
		CategoryNameEN: "test category",
		Description:    mockStr,
		Price:          mockNumber,
		IsAvailable:    true,
		ImageUrl:       mockStr,
	}

	mockupResult := &domain.Product{
		ID:             mockupID,
		Name:           "ทดสอบ",
		NameEN:         "test",
		CategoryID:     mockupID,
		CategoryName:   "ทดสอบหมวดหมู่",
		CategoryNameEN: "test category",
		Description:    utils.PgTextToStringPtr(mockupData.Description),
		Price:          utils.PgNumericToFloat64(mockNumber),
		IsAvailable:    true,
		ImageURL:       utils.PgTextToStringPtr(mockupData.ImageUrl),
	}

	testCases := []struct {
		name          string
		buildStubs    func(db *mockdb.MockStore)
		checkResponse func(t *testing.T, result *domain.Product, customError *exceptions.CustomError)
	}{
		{
			name: "OK",
			buildStubs: func(db *mockdb.MockStore) {
				db.EXPECT().GetProductAvailableByID(gomock.Any(), mockupID).Times(1).Return(mockupData, nil)
			},
			checkResponse: func(t *testing.T, result *domain.Product, customError *exceptions.CustomError) {
				require.Nil(t, customError)
				require.Equal(t, mockupResult, result)
			},
		},
		{
			name: "Internal Server Error",
			buildStubs: func(db *mockdb.MockStore) {
				db.EXPECT().GetProductAvailableByID(gomock.Any(), mockupID).Times(1).Return(nil, errors.New("something error"))
			},
			checkResponse: func(t *testing.T, result *domain.Product, customError *exceptions.CustomError) {
				require.Error(t, customError.Errors)
				require.Equal(t, exceptions.ERRREPOSITORY, customError.Status)
				require.Nil(t, result)
			},
		},
		{
			name: "Data Not Found",
			buildStubs: func(db *mockdb.MockStore) {
				db.EXPECT().GetProductAvailableByID(gomock.Any(), mockupID).Times(1).Return(nil, nil)
			},
			checkResponse: func(t *testing.T, result *domain.Product, customError *exceptions.CustomError) {
				require.Error(t, customError.Errors)
				require.Equal(t, exceptions.ERRREPOSITORY, customError.Status)
				require.Nil(t, result)
			},
		},
	}

	cfg := config.Config{}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			redis := mockcache.NewMockRedisTableCacheInterface(ctrl)
			node := mockshared.NewMockSnowflakeInterface(ctrl)

			tc.buildStubs(store)
			repo := repository.NewRepository(cfg, store, node)
			useCase := NewUsecase(cfg, *repo, redis)

			ctx := context.Background()
			result, customError := useCase.GetProductByID(ctx, mockupID)
			tc.checkResponse(t, result, customError)
		})

	}
}
