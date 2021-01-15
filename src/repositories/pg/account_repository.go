package pg

import (
	"simple-mpesa/src/domain/account"
	"simple-mpesa/src/domain/value_objects"
	"simple-mpesa/src/errors"
	"simple-mpesa/src/repositories/schema"
	"simple-mpesa/src/storage"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

func NewAccountRepository(database *storage.Database) *AccountRepository {
	return &AccountRepository{db: database}
}

type AccountRepository struct {
	db *storage.Database
}

// GetAccountByUserID fetches an account tied to a user's id
func (r AccountRepository) GetAccountByUserID(userID uuid.UUID) (account.Account, error) {
	var row schema.Account
	result := r.db.Where(&schema.Account{UserID: userID}).First(&row)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return account.Account{}, errors.Error{Code: errors.ENOTFOUND}
	} else if result.Error != nil {
		return account.Account{}, errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}

	return account.Account{
		ID:               row.ID,
		AvailableBalance: value_objects.Cents(row.AvailableBalance),
		Status:           account.Status(row.Status),
		AccountType:      account.Type(row.AccountType),
		UserID:           row.UserID,
	}, nil
}

// UpdateBalance
func (r AccountRepository) UpdateBalance(amount value_objects.Cents, userID uuid.UUID) (account.Account, error) {
	var row schema.Account
	result := r.db.Model(&schema.Account{}).
		Where(schema.Account{UserID: userID}).
		Updates(schema.Account{AvailableBalance: uint(amount)}).
		Scan(&row)
	if err := result.Error; err != nil {
		return account.Account{}, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	return account.Account{
		ID:               row.ID,
		AvailableBalance: value_objects.Cents(row.AvailableBalance),
		Status:           account.Status(row.Status),
		AccountType:      account.Type(row.AccountType),
		UserID:           row.UserID,
	}, nil
}

// Create a now account for userId
func (r AccountRepository) Create(userId uuid.UUID) (account.Account, error) {
	// check if user has an account and return it, otherwise create an account for user
	row := zeroAccount(userId)
	result := r.db.Where(schema.Account{UserID: userId}).FirstOrCreate(&row)
	if err := result.Error; err != nil {
		// we check if the error is a postgres unique constraint violation
		if pgerr, ok := result.Error.(*pgconn.PgError); ok && pgerr.Code == "23505" {
			return account.Account{}, errors.Error{Code: errors.ECONFLICT, Message: errors.ErrUserHasAccount(userId, row.ID)}
		}
		return account.Account{}, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	return account.Account{
		ID:               row.ID,
		AvailableBalance: value_objects.Cents(row.AvailableBalance),
		Status:           account.Status(row.Status),
		AccountType:      account.Type(row.AccountType),
		UserID:           row.UserID,
	}, nil
}

func zeroAccount(userId uuid.UUID) schema.Account {
	id, _ := uuid.NewV4()

	return schema.Account{
		ID: id,
		// balance:     0, // no need to initialize with zero value, Go will do that for us
		Status:      string(account.StatusActive),
		AccountType: string(account.AccTypeCurrent),
		UserID:      userId,
	}
}
