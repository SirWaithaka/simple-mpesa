package subscriber

import (
	"github.com/gofrs/uuid"
)

// Subscriber
type Subscriber struct {
	ID    uuid.UUID
	Email string // email is used as account number

	FirstName   string
	LastName    string
	PhoneNumber string
	PassportNo  string
	Password    string
}

type Repository interface {
	Add(Subscriber) (Subscriber, error)
	Delete(Subscriber) error
	FetchAll() ([]Subscriber, error)
	FindByID(uuid.UUID) (Subscriber, error)
	FindByEmail(string) (Subscriber, error)
	Update(Subscriber) error
}
