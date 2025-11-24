package repository

import (
	"context"
	"food-story/pkg/exceptions"
	"food-story/pkg/utils"
	"food-story/table-service/internal/domain"
)

func (i *Implement) ListSessionExtensionReason(ctx context.Context) (result []*domain.ListSessionExtensionReason, err error) {
	const _errMessage = "failed to fetch session extension reason"

	data, err := i.repository.ListSessionExtensionReason(ctx)
	if err != nil {
		return nil, exceptions.Errorf(exceptions.CodeRepository, _errMessage, err)
	}

	if data == nil {
		return nil, exceptions.Errorf(exceptions.CodeRepository, _errMessage, err)
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
