package repository

import (
	"context"
	"food-story/menu-service/internal/domain"
	"food-story/pkg/exceptions"
	"food-story/pkg/utils"
)

func (i *Implement) ListCategory(ctx context.Context) (result []*domain.Category, err error) {

	const _errMsg = "failed to fetch category"

	data, err := i.repository.ListCategory(ctx)
	if err != nil {
		return nil, exceptions.Errorf(exceptions.CodeRepository, _errMsg, err)
	}

	if data == nil {
		return nil, exceptions.Errorf(exceptions.CodeRepository, _errMsg, err)
	}

	result = make([]*domain.Category, len(data))
	for index, v := range data {
		result[index] = &domain.Category{
			ID:     v.ID,
			Name:   v.Name,
			NameEn: v.NameEN,
			Icon:   utils.PgTextToStringPtr(v.Icon),
		}
	}
	return result, nil
}
