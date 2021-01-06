package repositories

import (
	"simple-mpesa/src/errors"
	"simple-mpesa/src/models"
	"simple-mpesa/src/storage"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

func NewMerchantRepository(database *storage.Database) *Merchant {
	return &Merchant{db: database}
}

type Merchant struct {
	db *storage.Database
}

func (r Merchant) searchBy(row models.Merchant) (models.Merchant, error) {
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
func (r Merchant) Add(merchant models.Merchant) (models.Merchant, error) {
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
func (r Merchant) Delete(merchant models.Merchant) error {
	result := r.db.Delete(&merchant)
	if result.Error != nil {
		return errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}
	return nil
}

// FetchAll gets all merchants in db
func (r Merchant) FetchAll() ([]models.Merchant, error) {
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
func (r Merchant) FindByID(id uuid.UUID) (models.Merchant, error) {
	merchant, err := r.searchBy(models.Merchant{ID: id})
	return merchant, err
}

// FindByEmail searches merchant by email
func (r Merchant) FindByEmail(email string) (models.Merchant, error) {
	merchant, err := r.searchBy(models.Merchant{Email: email})
	return merchant, err
}

// Update
func (r Merchant) Update(merchant models.Merchant) error {
	var merch models.Merchant
	result := r.db.Model(&merch).Omit("id").Updates(merchant)
	if err := result.Error; err != nil {
		return errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}
	return nil
}
