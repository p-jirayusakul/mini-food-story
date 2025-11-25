package usecase

import (
	"context"
	"food-story/pkg/exceptions"
	"food-story/pkg/utils"
	shareModel "food-story/shared/model"

	"github.com/google/uuid"
)

func (i *Implement) GetOrderIDFromSession(sessionID uuid.UUID) (result int64, err error) {
	session, err := i.cache.GetCachedTable(sessionID)
	if err != nil {
		return 0, err
	}

	if session.OrderID == nil {
		return 0, exceptions.Error(exceptions.CodeNotFound, exceptions.ErrOrderNotFound.Error())
	}

	orderID, err := utils.StrToInt64(*session.OrderID)
	if err != nil {
		return 0, exceptions.Errorf(exceptions.CodeSystem, "failed to convert oder id", err)
	}

	return orderID, nil
}

func (i *Implement) GetCurrentTableSession(sessionID uuid.UUID) (result shareModel.CurrentTableSession, err error) {
	session, err := i.cache.GetCachedTable(sessionID)
	if err != nil {
		return shareModel.CurrentTableSession{}, err
	}

	if session == nil {
		return shareModel.CurrentTableSession{}, exceptions.Error(exceptions.CodeNotFound, exceptions.ErrSessionNotFound.Error())
	}

	if session.OrderID == nil {
		return shareModel.CurrentTableSession{}, exceptions.Error(exceptions.CodeNotFound, exceptions.ErrOrderNotFound.Error())
	}

	return *session, nil
}

func (i *Implement) GetSessionIDByTableID(ctx context.Context, tableID int64) (result uuid.UUID, err error) {
	return i.repository.GetSessionIDByTableID(ctx, tableID)
}

func (i *Implement) IsSessionValid(sessionID uuid.UUID) error {
	return i.cache.IsCachedTableExist(sessionID)
}
