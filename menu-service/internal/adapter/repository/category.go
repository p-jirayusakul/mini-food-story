package repository

import (
	"context"
	"errors"
	"fmt"
	"food-story/menu-service/internal/domain"
	"food-story/pkg/exceptions"
	"food-story/pkg/utils"
)

func (i *Implement) ListCategory(ctx context.Context) (result []*domain.Category, customError *exceptions.CustomError) {
	data, err := i.repository.ListCategory(ctx)
	if err != nil {
		return nil, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to check product exists: %w", err),
		}
	}

	if data == nil {
		return nil, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: errors.New("category not found"),
		}
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
