package usecase

import (
	"context"
	"errors"
	"food-story/pkg/exceptions"
	database "food-story/shared/database/sqlc"
	"food-story/shared/mock/cache/table"
	"food-story/shared/mock/database"
	"food-story/shared/mock/shared"
	"food-story/table-service/internal/adapter/repository"
	"food-story/table-service/internal/domain"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestListTableStatus(t *testing.T) {

	mockupData := []*database.ListTableStatusRow{
		{
			ID:     1,
			Name:   "name1",
			NameEn: "nameEN1",
			Code:   "code1",
		},
		{
			ID:     2,
			Name:   "name2",
			NameEn: "nameEN2",
			Code:   "code2",
		},
	}

	mockupResult := []*domain.Status{
		{
			ID:     1,
			Name:   "name1",
			NameEn: "nameEN1",
			Code:   "code1",
		},
		{
			ID:     2,
			Name:   "name2",
			NameEn: "nameEN2",
			Code:   "code2",
		},
	}

	testCases := []struct {
		name          string
		buildStubs    func(db *mockdb.MockStore)
		checkResponse func(t *testing.T, result []*domain.Status, customError *exceptions.CustomError)
	}{
		{
			name: "OK",
			buildStubs: func(db *mockdb.MockStore) {
				db.EXPECT().ListTableStatus(gomock.Any()).Times(1).Return(mockupData, nil)
			},
			checkResponse: func(t *testing.T, result []*domain.Status, customError *exceptions.CustomError) {
				require.Nil(t, customError)
				require.Equal(t, mockupResult, result)
			},
		},
		{
			name: "Internal Server Error",
			buildStubs: func(db *mockdb.MockStore) {
				db.EXPECT().ListTableStatus(gomock.Any()).Times(1).Return(nil, errors.New("something error"))
			},
			checkResponse: func(t *testing.T, result []*domain.Status, customError *exceptions.CustomError) {
				require.Error(t, customError.Errors)
				require.Equal(t, exceptions.ERRREPOSITORY, customError.Status)
				require.Nil(t, result)
			},
		},
		{
			name: "Data Not Found",
			buildStubs: func(db *mockdb.MockStore) {
				db.EXPECT().ListTableStatus(gomock.Any()).Times(1).Return(nil, nil)
			},
			checkResponse: func(t *testing.T, result []*domain.Status, customError *exceptions.CustomError) {
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
			result, customError := useCase.ListTableStatus(ctx)
			tc.checkResponse(t, result, customError)
		})

	}
}
