package user

import (
	"simple-wallet/app/errors"
	"simple-wallet/app/models"
	"simple-wallet/app/storage"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

type Repository interface {
	Add(models.User) (models.User, error)
	Delete(models.User) error
	GetByID(uuid.UUID) (models.User, error)
	GetByEmail(string) (models.User, error)
	GetByPhoneNumber(string) (models.User, error)
	GetByEmailOrPhone(string) (models.User, error)
	Update(models.User) error
}

func NewRepository(database *storage.Database) Repository {
	return &repository{db: database}
}

type repository struct {
	db *storage.Database
}

func (r repository) searchBy(row models.User) (models.User, error) {
	var user models.User
	result := r.db.Where(row).First(&user)
	// check if no record found.
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return models.User{}, errors.Error{Code: errors.ENOTFOUND}
	}
	if err := result.Error; err != nil {
		return models.User{}, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	return user, nil
}

// Add a user if already not in db.
func (r repository) Add(user models.User) (models.User, error) {
	// add new user to users table, if query return violation of unique key column,
	// we know that the user with given record exists and return that user instead
	result := r.db.Model(models.User{}).Create(&user)
	if err := result.Error; err != nil {
		// we check if the error is a postgres unique constraint violation
		if pgerr, ok := err.(*pgconn.PgError); ok && pgerr.Code == "23505" {
			return user, errors.Error{Code: errors.ECONFLICT, Message: errors.ErrUserExists}
		}
		return models.User{}, errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}

	return user, nil
}

// Delete a user
func (r repository) Delete(user models.User) error {
	result := r.db.Delete(&user)
	if result.Error != nil {
		return errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}
	return nil
}

// GetByID searches user by primary id
func (r repository) GetByID(id uuid.UUID) (models.User, error) {
	user, err := r.searchBy(models.User{ID: id})
	return user, err
}

// GetByEmail searches user by email
func (r repository) GetByEmail(email string) (models.User, error) {
	user, err := r.searchBy(models.User{Email: email})
	return user, err
}

// GetByPhoneNumber searches user by phone number
func (r repository) GetByPhoneNumber(phoneNo string) (models.User, error) {
	user, err := r.searchBy(models.User{PhoneNumber: phoneNo})
	return user, err
}

// GetByEmailOrPhone
func (r repository) GetByEmailOrPhone(value string) (models.User, error) {
	var user models.User
	result := r.db.Model(models.User{}).Where(models.User{Email: value}).Or(models.User{PhoneNumber: value}).First(&user)
	// check if no record found.
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return models.User{}, errors.Error{Code: errors.ENOTFOUND}
	}
	if err := result.Error; err != nil {
		return models.User{}, errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}

	return user, nil
}

// Update
func (r repository) Update(user models.User) error {
	var u models.User
	result := r.db.Model(&u).Omit("id").Updates(user)
	if err := result.Error; err != nil {
		return errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}
	return nil
}
