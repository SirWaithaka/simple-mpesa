package tariff

import (
	"simple-mpesa/app/models"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Charge struct {
	ID uuid.UUID

	Transaction         models.TxnOperation `gorm:"uniqueIndex:idx_unique_tx_identity"`
	SourceUserType      models.UserType     `gorm:"uniqueIndex:idx_unique_tx_identity"`
	DestinationUserType models.UserType     `gorm:"uniqueIndex:idx_unique_tx_identity"`
	Fee                 models.Cents

	gorm.Model
}

func (t *Charge) BeforeCreate(tx *gorm.DB) error {
	t.ID, _ = uuid.NewV4()
	return nil
}

// ValidTransaction defines a format to identify all allowable transactions between customers
// It is an array of 2 user types since only 2 customers types are allowed in a valid transaction
// Index 0 is the source while index 1 is the destination
//
// Example
// 1. Validation{models.UserTypSubscriber, models.UserTypMerchant}
// Describes a valid transaction (transfer) between a subscriber to a merchant
type ValidTransaction [2]models.UserType
