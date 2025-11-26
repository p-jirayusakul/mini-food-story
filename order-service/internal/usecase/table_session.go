package usecase

import (
	"context"
	"food-story/order-service/internal/domain"
	"food-story/pkg/exceptions"
	"food-story/pkg/utils"

	"github.com/google/uuid"
)

func (i *Implement) GetOrderIDFromSession(sessionID uuid.UUID) (result int64, err error) {
	session, err := i.cache.GetCachedTable(sessionID)
	if err != nil {
		return 0, err
	}

	if session.OrderID == nil {
		return 0, exceptions.ErrorIDNotFound(exceptions.CodeOrderNotFound, 0)
	}

	orderID, err := utils.StrToInt64(*session.OrderID)
	if err != nil {
		return 0, exceptions.Errorf(exceptions.CodeSystem, "failed to convert oder id", err)
	}

	return orderID, nil
}

func (i *Implement) GetCurrentTableSession(sessionID uuid.UUID) (result domain.CurrentTableSession, err error) {
	session, err := i.cache.GetCachedTable(sessionID)
	if err != nil {
		return domain.CurrentTableSession{}, err
	}

	if session == nil {
		return domain.CurrentTableSession{}, exceptions.ErrorSessionNotFound()
	}

	if session.OrderID == nil {
		return domain.CurrentTableSession{}, exceptions.ErrorIDNotFound(exceptions.CodeOrderNotFound, 0)
	}

	return *session, nil
}

func (i *Implement) GetSessionIDByTableID(ctx context.Context, tableID int64) (result uuid.UUID, err error) {
	return i.repository.GetSessionIDByTableID(ctx, tableID)
}

func (i *Implement) IsSessionValid(sessionID uuid.UUID) error {
	return i.cache.IsCachedTableExist(sessionID)
}
