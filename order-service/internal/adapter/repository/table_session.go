package repository

import (
	"context"
	"food-story/pkg/exceptions"
	"food-story/pkg/utils"

	"github.com/google/uuid"
)

func (i *Implement) UpdateStatusCloseTableSession(ctx context.Context, sessionID uuid.UUID) (err error) {
	err = i.repository.UpdateStatusCloseTableSession(ctx, utils.UUIDToPgUUID(sessionID))
	if err != nil {
		return exceptions.Errorf(exceptions.CodeRepository, "failed to update status close table session", err)
	}

	return nil
}

func (i *Implement) GetSessionIDByTableID(ctx context.Context, tableID int64) (result uuid.UUID, err error) {
	sessionID, err := i.repository.GetSessionIDByTableID(ctx, tableID)
	if err != nil {
		return uuid.Nil, exceptions.Errorf(exceptions.CodeRepository, "failed to get session by table id", err)
	}

	v, err := utils.PareStringToUUID(sessionID.String())
	if err != nil {
		return uuid.Nil, exceptions.Errorf(exceptions.CodeSystem, "failed to parse session id", err)
	}

	return v, nil
}
