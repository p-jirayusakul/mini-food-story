package usecase

import (
	"food-story/pkg/exceptions"
	"github.com/google/uuid"
)

func (i *Implement) IsSessionValid(sessionID uuid.UUID) *exceptions.CustomError {
	err := i.cache.IsCachedTableExist(sessionID)
	if err != nil {
		return &exceptions.CustomError{
			Status: exceptions.ERRAUTHORIZED,
			Errors: exceptions.ErrSessionExpired,
		}
	}
	return nil
}
