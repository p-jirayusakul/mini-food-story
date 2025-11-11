package repository

import (
	"context"
	"errors"
	"fmt"
	"food-story/pkg/exceptions"
	"food-story/pkg/utils"
	"food-story/table-service/internal/domain"
)

func (i *Implement) ListSessionExtensionReason(ctx context.Context) (result []*domain.ListSessionExtensionReason, customError *exceptions.CustomError) {
	data, err := i.repository.ListSessionExtensionReason(ctx)
	if err != nil {
		return nil, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: fmt.Errorf("failed to fetch session extension reason: %w", err),
		}
	}

	if data == nil {
		return nil, &exceptions.CustomError{
			Status: exceptions.ERRREPOSITORY,
			Errors: errors.New("session extension reason not found"),
		}
	}

	result = make([]*domain.ListSessionExtensionReason, len(data))
	for index, v := range data {
		result[index] = &domain.ListSessionExtensionReason{
			ID:       v.ID,
			Code:     v.Code,
			Name:     v.Name,
			NameEN:   v.NameEN,
			Category: utils.PgTextToStringPtr(v.Category),
			ModeCode: utils.PgTextToStringPtr(v.ModeCode),
		}
	}

	return result, nil
}
