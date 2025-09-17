package usecase

import (
	"food-story/pkg/exceptions"

	"github.com/google/uuid"
)

func (i *Implement) IsSessionValid(sessionID uuid.UUID) *exceptions.CustomError {
	return i.cache.IsCachedTableExist(sessionID)
}
