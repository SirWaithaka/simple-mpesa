package pg

import (
	"simple-mpesa/src/domain/merchant"
	"simple-mpesa/src/errors"
	"simple-mpesa/src/storage"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

func NewMerchantRepository(database *storage.Database) *MerchantRepository {
	return &MerchantRepository{db: database}
}

type MerchantRepository struct {
	db *storage.Database
}

func (r MerchantRepository) searchBy(row merchant.Merchant) (merchant.Merchant, error) {
	var m merchant.Merchant
	result := r.db.Where(row).First(&m)
	// check if no record found.
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return merchant.Merchant{}, errors.Error{Code: errors.ENOTFOUND}
	}
	if err := result.Error; err != nil {
		return merchant.Merchant{}, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	return m, nil
}

// Add a merchant if already not in db.
func (r MerchantRepository) Add(m merchant.Merchant) (merchant.Merchant, error) {
	// add new merchant to merchants table, if query return violation of unique key column,
	// we know that the merchant with given record exists and return that merchant instead
	result := r.db.Model(merchant.Merchant{}).Create(&m)
	if err := result.Error; err != nil {
		// we check if the error is a postgres unique constraint violation
		if pgerr, ok := err.(*pgconn.PgError); ok && pgerr.Code == "23505" {
			return m, errors.Error{Code: errors.ECONFLICT, Message: errors.ErrUserExists}
		}
		return merchant.Merchant{}, errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}

	return m, nil
}

// Delete a merchant
func (r MerchantRepository) Delete(m merchant.Merchant) error {
	result := r.db.Delete(&m)
	if result.Error != nil {
		return errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}
	return nil
}

// FetchAll gets all merchants in db
func (r MerchantRepository) FetchAll() ([]merchant.Merchant, error) {
	var merchants []merchant.Merchant
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
func (r MerchantRepository) FindByID(id uuid.UUID) (merchant.Merchant, error) {
	m, err := r.searchBy(merchant.Merchant{ID: id})
	return m, err
}

// FindByEmail searches merchant by email
func (r MerchantRepository) FindByEmail(email string) (merchant.Merchant, error) {
	m, err := r.searchBy(merchant.Merchant{Email: email})
	return m, err
}

// Update
func (r MerchantRepository) Update(m merchant.Merchant) error {
	result := r.db.Model(&merchant.Merchant{}).Omit("id").Updates(m)
	if err := result.Error; err != nil {
		return errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}
	return nil
}
