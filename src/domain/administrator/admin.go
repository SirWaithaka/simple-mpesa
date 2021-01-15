package administrator

import (
	"github.com/gofrs/uuid"
)

// Admin
type Admin struct {
	ID        uuid.UUID
	FirstName string
	LastName  string
	Email     string
	Password  string
}

type Repository interface {
	Add(Admin) (Admin, error)
	Delete(Admin) error
	GetByID(uuid.UUID) (Admin, error)
	GetByEmail(string) (Admin, error)
	Update(Admin) error
}
