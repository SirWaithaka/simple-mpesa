package repositories

import (
	"simple-mpesa/src/admin"
	"simple-mpesa/src/errors"
	"simple-mpesa/src/storage"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

// NewAdminRepository creates and returns a new instance of admin repository
func NewAdminRepository(database *storage.Database) *AdminRepository {
	return &AdminRepository{db: database}
}

type AdminRepository struct {
	db *storage.Database
}

func (r AdminRepository) searchBy(row admin.Administrator) (admin.Administrator, error) {
	var adm admin.Administrator
	result := r.db.Where(row).First(&adm)
	// check if no record found.
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return admin.Administrator{}, errors.Error{Code: errors.ENOTFOUND}
	}
	if err := result.Error; err != nil {
		return admin.Administrator{}, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	return adm, nil
}

// Add an admin if already not in db.
func (r AdminRepository) Add(adm admin.Administrator) (admin.Administrator, error) {
	// add new admin to administrators table, if query return violation of unique key column,
	// we know that the admin with given record exists and return that admin instead
	result := r.db.Model(admin.Administrator{}).Create(&adm)
	if err := result.Error; err != nil {
		// we check if the error is a postgres unique constraint violation
		if pgerr, ok := err.(*pgconn.PgError); ok && pgerr.Code == "23505" {
			return adm, errors.Error{Code: errors.ECONFLICT, Message: errors.ErrUserExists}
		}
		return admin.Administrator{}, errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}

	return adm, nil
}

// Delete a user
func (r AdminRepository) Delete(adm admin.Administrator) error {
	// TODO("If deleting permanently, delete record from users table too")
	result := r.db.Delete(&adm)
	if result.Error != nil {
		return errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}
	return nil
}

// GetByID searches admin by primary id
func (r AdminRepository) GetByID(id uuid.UUID) (admin.Administrator, error) {
	adm, err := r.searchBy(admin.Administrator{ID: id})
	return adm, err
}

// GetByEmail searches admin by email
func (r AdminRepository) GetByEmail(email string) (admin.Administrator, error) {
	adm, err := r.searchBy(admin.Administrator{Email: email})
	return adm, err
}

// Update details of an amin
func (r AdminRepository) Update(adm admin.Administrator) error {
	result := r.db.Model(&admin.Administrator{}).Omit("id").Updates(adm)
	if err := result.Error; err != nil {
		return errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}
	return nil
}
