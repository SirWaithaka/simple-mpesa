package admin

import (
	"github.com/gofrs/uuid"
)

// Administrator
type Administrator struct {
	ID        uuid.UUID
	FirstName string
	LastName  string
	Email     string
	Password  string
}

type Repository interface {
	Add(Administrator) (Administrator, error)
	Delete(Administrator) error
	GetByID(uuid.UUID) (Administrator, error)
	GetByEmail(string) (Administrator, error)
	Update(Administrator) error
}
