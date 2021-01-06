package statement

import (
	"time"

	"github.com/gofrs/uuid"
)

type Repository interface {
	Add(Statement) (Statement, error)
	GetStatements(userID uuid.UUID, from time.Time, limit uint) ([]Statement, error)
}
