package repository

import (
	"context"
	"food-story/menu-service/internal/domain"
	"food-story/pkg/exceptions"
	"food-story/pkg/utils"
)

func (i *Implement) ListCategory(ctx context.Context) (result []*domain.Category, err error) {

	data, err := i.repository.ListCategory(ctx)
	if err != nil {
		return nil, exceptions.Errorf(exceptions.CodeRepository, "failed to fetch category", err)
	}

	if data == nil {
		return nil, exceptions.ErrorDataNotFound()
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
