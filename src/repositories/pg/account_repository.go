package pg

import (
	"simple-mpesa/src/domain/account"
	"simple-mpesa/src/domain/value_objects"
	"simple-mpesa/src/errors"
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
	var acc account.Account
	result := r.db.Where(account.Account{UserID: userID}).First(&acc)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return account.Account{}, errors.Error{Code: errors.ENOTFOUND}
	} else if result.Error != nil {
		return account.Account{}, errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}

	return acc, nil
}

// UpdateBalance
func (r AccountRepository) UpdateBalance(amount value_objects.Cents, userID uuid.UUID) (account.Account, error) {
	var acc account.Account
	result := r.db.Model(account.Account{}).Where(account.Account{UserID: userID}).Updates(account.Account{AvailableBalance: amount}).Scan(&acc)
	if err := result.Error; err != nil {
		return account.Account{}, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	return acc, nil
}

// Create a now account for userId
func (r AccountRepository) Create(userId uuid.UUID) (account.Account, error) {
	// check if user has an account and return it, otherwise create an account for user
	acc := zeroAccount(userId)
	result := r.db.Where(account.Account{UserID: userId}).FirstOrCreate(&acc)
	if err := result.Error; err != nil {
		// we check if the error is a postgres unique constraint violation
		if pgerr, ok := result.Error.(*pgconn.PgError); ok && pgerr.Code == "23505" {
			return account.Account{}, errors.Error{Code: errors.ECONFLICT, Message: errors.ErrUserHasAccount(userId, acc.ID)}
		}
		return account.Account{}, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	return acc, nil
}

func zeroAccount(userId uuid.UUID) account.Account {
	id, _ := uuid.NewV4()

	return account.Account{
		ID: id,
		// balance:     0, // no need to initialize with zero value, Go will do that for us
		Status:      account.StatusActive,
		AccountType: account.AccTypeCurrent,
		UserID:      userId,
	}
}
