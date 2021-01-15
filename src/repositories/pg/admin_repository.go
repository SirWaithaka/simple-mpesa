package pg

import (
	"simple-mpesa/src/domain/administrator"
	"simple-mpesa/src/errors"
	"simple-mpesa/src/repositories/schema"
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

func (r AdminRepository) searchBy(row schema.Admin) (administrator.Admin, error) {
	var adm schema.Admin
	result := r.db.Where(row).First(&adm)
	// check if no record found.
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return administrator.Admin{}, errors.Error{Code: errors.ENOTFOUND}
	}
	if err := result.Error; err != nil {
		return administrator.Admin{}, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	return administrator.Admin{
		ID:        adm.ID,
		FirstName: adm.FirstName,
		LastName:  adm.LastName,
		Email:     adm.Email,
		Password:  adm.Password,
	}, nil
}

// Add an admin if already not in db.
func (r AdminRepository) Add(admin administrator.Admin) (administrator.Admin, error) {
	adm := schema.Admin{
		FirstName: admin.FirstName,
		LastName:  admin.LastName,
		Email:     admin.Email,
		Password:  admin.Password,
	}

	// add new admin to administrators table, if query return violation of unique key column,
	// we know that the admin with given record exists and return that admin instead
	result := r.db.Model(schema.Admin{}).Create(&adm)
	if err := result.Error; err != nil {
		// we check if the error is a postgres unique constraint violation
		if pgerr, ok := err.(*pgconn.PgError); ok && pgerr.Code == "23505" {
			return admin, errors.Error{Code: errors.ECONFLICT, Message: errors.ErrUserExists}
		}
		return administrator.Admin{}, errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}

	return administrator.Admin{
		ID:        adm.ID,
		FirstName: adm.FirstName,
		LastName:  adm.LastName,
		Email:     adm.Email,
		Password:  adm.Password,
	}, nil
}

// Delete a user
func (r AdminRepository) Delete(admin administrator.Admin) error {
	// TODO("If deleting permanently, delete record from users table too")
	result := r.db.Delete(&schema.Admin{ID: admin.ID})
	if result.Error != nil {
		return errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}
	return nil
}

// GetByID searches admin by primary id
func (r AdminRepository) GetByID(id uuid.UUID) (administrator.Admin, error) {
	adm, err := r.searchBy(schema.Admin{ID: id})
	return adm, err
}

// GetByEmail searches admin by email
func (r AdminRepository) GetByEmail(email string) (administrator.Admin, error) {
	adm, err := r.searchBy(schema.Admin{Email: email})
	return adm, err
}

// Update details of an amin
func (r AdminRepository) Update(admin administrator.Admin) error {
	result := r.db.Model(&schema.Admin{}).Where(&schema.Admin{ID: admin.ID}).Omit("id").Updates(admin)
	if err := result.Error; err != nil {
		return errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}
	return nil
}
