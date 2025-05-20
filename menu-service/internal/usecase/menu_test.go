package usecase

import (
	"context"
	"errors"
	"food-story/menu-service/internal/adapter/repository"
	"food-story/menu-service/internal/domain"
	"food-story/pkg/exceptions"
	"food-story/shared/config"
	database "food-story/shared/database/sqlc"
	mockdb "food-story/shared/mock/database"
	mocksnowflake "food-story/shared/mock/snowflake"

	mockcache "food-story/shared/mock/menu/cache"

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
				require.Equal(t, exceptions.ERRNOTFOUND, customError.Status)
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
			node := mocksnowflake.NewMockSnowflakeInterface(ctrl)

			tc.buildStubs(store)
			repo := repository.NewRepository(cfg, store, node)
			useCase := NewUsecase(cfg, *repo, redis)

			ctx := context.Background()
			result, customError := useCase.ListCategory(ctx)
			tc.checkResponse(t, result, customError)
		})

	}
}
