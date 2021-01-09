package repositories

import (
	"time"

	"simple-mpesa/src/value_objects"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

// Account entity definition
type Account struct {
	ID uuid.UUID

	// balance will be stored in cents
	AvailableBalance uint `gorm:"column:available_balance"`

	Status      string    `gorm:"column:status"`
	AccountType string    `gorm:"column:account_type"`
	UserID      uuid.UUID `gorm:"column:user_id;not null;unique"` // a user can only have one account

	gorm.Model
}

// Admin
type Admin struct {
	ID uuid.UUID

	FirstName string
	LastName  string
	Email     string `gorm:"not null;unique"`
	Password  string

	gorm.Model
}

// BeforeCreate hook will be used to add uuid to entity before adding to db
func (u *Admin) BeforeCreate(tx *gorm.DB) error {
	u.ID, _ = uuid.NewV4()
	return nil
}

func (Admin) TableName() string {
	return "administrators"
}

// Agent
type Agent struct {
	ID    uuid.UUID
	Email string `gorm:"not null;unique"` // email is used as account number

	FirstName   string
	LastName    string
	PhoneNumber string `gorm:"not null;unique"`
	PassportNo  string
	Password    string `gorm:"not null"`

	// an extra column/property that tells us if the agent is a super agent
	SuperAgent string `gorm:"default:'0'"` // PS: bool values dont work well with gorm during updates

	gorm.Model
}

type Statement struct {
	ID           uuid.UUID
	Operation    value_objects.TxnOperation
	DebitAmount  float64
	CreditAmount float64
	UserID       uuid.UUID
	AccountID    uuid.UUID
	CreatedAt    time.Time
}

func (s *Statement) BeforeCreate(tx *gorm.DB) error {
	s.ID, _ = uuid.NewV4()
	return nil
}

func (Statement) TableName() string {
	return "statements"
}

// Merchant
type Merchant struct {
	ID    uuid.UUID
	Email string `gorm:"not null;unique"` // email is used as account number

	FirstName   string
	LastName    string
	PhoneNumber string `gorm:"not null;unique"`
	PassportNo  string
	Password    string `gorm:"not null"`

	// a merchant is usually assigned a till number they use to accept
	// payments from other customers
	TillNumber string `gorm:"column:till_number;unique"`

	gorm.Model
}

// BeforeCreate hook will be used to add uuid to entity before adding to db
func (u *Merchant) BeforeCreate(tx *gorm.DB) error {
	u.ID, _ = uuid.NewV4()
	return nil
}

// Subscriber
type Subscriber struct {
	ID    uuid.UUID
	Email string `gorm:"not null;unique"` // email is used as account number

	FirstName   string
	LastName    string
	PhoneNumber string `gorm:"not null;unique"`
	PassportNo  string
	Password    string `gorm:"not null"`

	gorm.Model
}

// BeforeCreate hook will be used to add uuid to entity before adding to db
func (u *Subscriber) BeforeCreate(tx *gorm.DB) error {
	u.ID, _ = uuid.NewV4()
	return nil
}

type Charge struct {
	ID uuid.UUID

	Transaction         string `gorm:"uniqueIndex:idx_unique_tx_identity"`
	SourceUserType      string `gorm:"uniqueIndex:idx_unique_tx_identity"`
	DestinationUserType string `gorm:"uniqueIndex:idx_unique_tx_identity"`
	Fee                 uint

	gorm.Model
}

func (t *Charge) BeforeCreate(tx *gorm.DB) error {
	t.ID, _ = uuid.NewV4()
	return nil
}
