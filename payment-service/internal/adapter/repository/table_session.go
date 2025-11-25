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
