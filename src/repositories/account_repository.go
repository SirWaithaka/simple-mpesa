package repositories

import (
	"simple-mpesa/src/errors"
	"simple-mpesa/src/models"
	"simple-mpesa/src/storage"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

func NewAccountRepository(database *storage.Database) *Account {
	return &Account{db: database}
}

type Account struct {
	db *storage.Database
}

// GetAccountByUserID fetches an account tied to a user's id
func (r Account) GetAccountByUserID(userID uuid.UUID) (models.Account, error) {
	var acc models.Account
	result := r.db.Where(models.Account{UserID: userID}).First(&acc)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return models.Account{}, errors.Error{Code: errors.ENOTFOUND}
	} else if result.Error != nil {
		return models.Account{}, errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}

	return acc, nil
}

// UpdateBalance
func (r Account) UpdateBalance(amount models.Cents, userID uuid.UUID) (models.Account, error) {
	var acc models.Account
	result := r.db.Model(models.Account{}).Where(models.Account{UserID: userID}).Updates(models.Account{AvailableBalance: amount}).Scan(&acc)
	if err := result.Error; err != nil {
		return models.Account{}, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	return acc, nil
}

// Create a now account for userId
func (r Account) Create(userId uuid.UUID) (models.Account, error) {
	// check if user has an account and return it, otherwise create an account for user
	acc := zeroAccount(userId)
	result := r.db.Where(models.Account{UserID: userId}).FirstOrCreate(&acc)
	if err := result.Error; err != nil {
		// we check if the error is a postgres unique constraint violation
		if pgerr, ok := result.Error.(*pgconn.PgError); ok && pgerr.Code == "23505" {
			return models.Account{}, errors.Error{Code: errors.ECONFLICT, Message: errors.ErrUserHasAccount(userId, acc.ID)}
		}
		return models.Account{}, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	return acc, nil
}

func zeroAccount(userId uuid.UUID) models.Account {
	id, _ := uuid.NewV4()

	return models.Account{
		ID: id,
		// balance:     0, // no need to initialize with zero value, Go will do that for us
		Status:      models.StatusActive,
		AccountType: models.AccTypeCurrent,
		UserID:      userId,
	}
}
