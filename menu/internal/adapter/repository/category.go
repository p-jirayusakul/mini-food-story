package repository

import (
	"context"
	"fmt"
	"food-story/menu/internal/domain"
	"food-story/pkg/exceptions"
)

func (i *ProductRepoImplement) ListCategory(ctx context.Context) (result []*domain.Category, customError *exceptions.CustomError) {
	data, err := i.repository.ListCategory(ctx)
	if err != nil {
		return nil, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to check product exists: %w", err),
		}
	}
	if data == nil {
		return nil, nil
	}

	result = make([]*domain.Category, len(data))
	for index, v := range data {
		result[index] = &domain.Category{
			ID:     v.ID,
			Name:   v.Name,
			NameEn: v.NameEN,
		}
	}
	return result, nil
}
