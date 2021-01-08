package merchant

import (
	"github.com/gofrs/uuid"
)

// Merchant
type Merchant struct {
	ID    uuid.UUID
	Email string // email is used as account number

	FirstName   string
	LastName    string
	PhoneNumber string
	PassportNo  string
	Password    string

	// a merchant is usually assigned a till number they use to accept
	// payments from other customers
	TillNumber string
}

type Repository interface {
	Add(Merchant) (Merchant, error)
	Delete(Merchant) error
	FetchAll() ([]Merchant, error)
	FindByID(uuid.UUID) (Merchant, error)
	FindByEmail(string) (Merchant, error)
	Update(Merchant) error
}
