package usecase

import (
	"github.com/google/uuid"
)

func (i *Implement) IsSessionValid(sessionID uuid.UUID) error {
	return i.cache.IsCachedTableExist(sessionID)
}
