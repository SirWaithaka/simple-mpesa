package admin

import (
	"simple-mpesa/app/errors"
	"simple-mpesa/app/models"
	"simple-mpesa/app/storage"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

type Repository interface {
	Add(models.Admin) (models.Admin, error)
	Delete(models.Admin) error
	GetByID(uuid.UUID) (models.Admin, error)
	GetByEmail(string) (models.Admin, error)
	Update(models.Admin) error
}

// NewRepository creates and returns a new instance of admin repository
func NewRepository(database *storage.Database) Repository {
	return &repository{db: database}
}

type repository struct {
	db *storage.Database
}

func (r repository) searchBy(row models.Admin) (models.Admin, error) {
	var admin models.Admin
	result := r.db.Where(row).First(&admin)
	// check if no record found.
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return models.Admin{}, errors.Error{Code: errors.ENOTFOUND}
	}
	if err := result.Error; err != nil {
		return models.Admin{}, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	return admin, nil
}

// Add an admin if already not in db.
func (r repository) Add(admin models.Admin) (models.Admin, error) {
	// add new admin to administrators table, if query return violation of unique key column,
	// we know that the admin with given record exists and return that admin instead
	result := r.db.Model(models.Admin{}).Create(&admin)
	if err := result.Error; err != nil {
		// we check if the error is a postgres unique constraint violation
		if pgerr, ok := err.(*pgconn.PgError); ok && pgerr.Code == "23505" {
			return admin, errors.Error{Code: errors.ECONFLICT, Message: errors.ErrUserExists}
		}
		return models.Admin{}, errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}

	return admin, nil
}

// Delete a user
func (r repository) Delete(admin models.Admin) error {
	// TODO("If deleting permanently, delete record from users table too")
	result := r.db.Delete(&admin)
	if result.Error != nil {
		return errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}
	return nil
}

// GetByID searches admin by primary id
func (r repository) GetByID(id uuid.UUID) (models.Admin, error) {
	admin, err := r.searchBy(models.Admin{ID: id})
	return admin, err
}

// GetByEmail searches admin by email
func (r repository) GetByEmail(email string) (models.Admin, error) {
	admin, err := r.searchBy(models.Admin{Email: email})
	return admin, err
}

// Update details of an amin
func (r repository) Update(admin models.Admin) error {
	var u models.Admin
	result := r.db.Model(&u).Omit("id").Updates(admin)
	if err := result.Error; err != nil {
		return errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}
	return nil
}
