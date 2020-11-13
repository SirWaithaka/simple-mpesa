package merchant

import (
	"simple-mpesa/app/errors"
	"simple-mpesa/app/models"
	"simple-mpesa/app/storage"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

type Repository interface {
	Add(models.Merchant) (models.Merchant, error)
	Delete(models.Merchant) error
	FetchAll() ([]models.Merchant, error)
	FindByID(uuid.UUID) (models.Merchant, error)
	FindByEmail(string) (models.Merchant, error)
	Update(models.Merchant) error
}

func NewRepository(database *storage.Database) Repository {
	return &repository{db: database}
}

type repository struct {
	db *storage.Database
}

func (r repository) searchBy(row models.Merchant) (models.Merchant, error) {
	var merchant models.Merchant
	result := r.db.Where(row).First(&merchant)
	// check if no record found.
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return models.Merchant{}, errors.Error{Code: errors.ENOTFOUND}
	}
	if err := result.Error; err != nil {
		return models.Merchant{}, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	return merchant, nil
}

// Add a merchant if already not in db.
func (r repository) Add(merchant models.Merchant) (models.Merchant, error) {
	// add new merchant to merchants table, if query return violation of unique key column,
	// we know that the merchant with given record exists and return that merchant instead
	result := r.db.Model(models.Merchant{}).Create(&merchant)
	if err := result.Error; err != nil {
		// we check if the error is a postgres unique constraint violation
		if pgerr, ok := err.(*pgconn.PgError); ok && pgerr.Code == "23505" {
			return merchant, errors.Error{Code: errors.ECONFLICT, Message: errors.ErrUserExists}
		}
		return models.Merchant{}, errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}

	return merchant, nil
}

// Delete a merchant
func (r repository) Delete(merchant models.Merchant) error {
	result := r.db.Delete(&merchant)
	if result.Error != nil {
		return errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}
	return nil
}

// FetchAll gets all merchants in db
func (r repository) FetchAll() ([]models.Merchant, error) {
	var merchants []models.Merchant
	result := r.db.Find(&merchants)
	if err := result.Error; err != nil {
		return nil, errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}

	// we might not need to return this error
	if result.RowsAffected == 0 {
		return nil, errors.Error{Code: errors.ENOTFOUND}
	}

	return merchants, nil
}

// FindByID searches merchant by primary id
func (r repository) FindByID(id uuid.UUID) (models.Merchant, error) {
	merchant, err := r.searchBy(models.Merchant{ID: id})
	return merchant, err
}

// FindByEmail searches merchant by email
func (r repository) FindByEmail(email string) (models.Merchant, error) {
	merchant, err := r.searchBy(models.Merchant{Email: email})
	return merchant, err
}

// Update
func (r repository) Update(merchant models.Merchant) error {
	var merch models.Merchant
	result := r.db.Model(&merch).Omit("id").Updates(merchant)
	if err := result.Error; err != nil {
		return errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}
	return nil
}
