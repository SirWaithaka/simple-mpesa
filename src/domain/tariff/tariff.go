package tariff

import (
	"simple-mpesa/src/domain/value_objects"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Charge struct {
	ID uuid.UUID

	Transaction         value_objects.TxnOperation
	SourceUserType      value_objects.UserType
	DestinationUserType value_objects.UserType
	Fee                 value_objects.Cents

	gorm.Model
}

// ValidTransaction defines a format to identify all allowable transactions between customers
// It is an array of 2 user types since only 2 customers types are allowed in a valid transaction
// Index 0 is the source while index 1 is the destination
//
// Example
// 1. Validation{value_objects.UserTypSubscriber, value_objects.UserTypMerchant}
// Describes a valid transaction (transfer) between a subscriber to a merchant
type ValidTransaction [2]value_objects.UserType

type Repository interface {
	Add(Charge) (Charge, error)
	FetchAll() ([]Charge, error)
	FindByID(uuid.UUID) (Charge, error)
	Get(operation value_objects.TxnOperation, src value_objects.UserType, dest value_objects.UserType) (Charge, error)
	Update(Charge) error
}
