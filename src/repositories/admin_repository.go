package repositories

import (
	"simple-mpesa/src/errors"
	"simple-mpesa/src/models"
	"simple-mpesa/src/storage"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

// NewAdminRepository creates and returns a new instance of admin repository
func NewAdminRepository(database *storage.Database) *Admin {
	return &Admin{db: database}
}

type Admin struct {
	db *storage.Database
}

func (r Admin) searchBy(row models.Admin) (models.Admin, error) {
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
func (r Admin) Add(admin models.Admin) (models.Admin, error) {
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
func (r Admin) Delete(admin models.Admin) error {
	// TODO("If deleting permanently, delete record from users table too")
	result := r.db.Delete(&admin)
	if result.Error != nil {
		return errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}
	return nil
}

// GetByID searches admin by primary id
func (r Admin) GetByID(id uuid.UUID) (models.Admin, error) {
	admin, err := r.searchBy(models.Admin{ID: id})
	return admin, err
}

// GetByEmail searches admin by email
func (r Admin) GetByEmail(email string) (models.Admin, error) {
	admin, err := r.searchBy(models.Admin{Email: email})
	return admin, err
}

// Update details of an amin
func (r Admin) Update(admin models.Admin) error {
	var u models.Admin
	result := r.db.Model(&u).Omit("id").Updates(admin)
	if err := result.Error; err != nil {
		return errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}
	return nil
}
