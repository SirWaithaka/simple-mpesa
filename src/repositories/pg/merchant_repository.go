package pg

import (
	"simple-mpesa/src/domain/merchant"
	"simple-mpesa/src/errors"
	"simple-mpesa/src/repositories/schema"
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

func (r MerchantRepository) searchBy(search schema.Merchant) (merchant.Merchant, error) {
	var row schema.Merchant
	result := r.db.Where(search).First(&row)
	// check if no record found.
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return merchant.Merchant{}, errors.Error{Code: errors.ENOTFOUND}
	}
	if err := result.Error; err != nil {
		return merchant.Merchant{}, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	return merchant.Merchant{
		ID:          row.ID,
		Email:       row.Email,
		FirstName:   row.FirstName,
		LastName:    row.LastName,
		PhoneNumber: row.PhoneNumber,
		PassportNo:  row.PassportNo,
		Password:    row.Password,
		TillNumber:  row.TillNumber,
	}, nil
}

// Add a merchant if already not in db.
func (r MerchantRepository) Add(m merchant.Merchant) (merchant.Merchant, error) {
	row := schema.Merchant{
		Email:       m.Email,
		FirstName:   m.FirstName,
		LastName:    m.LastName,
		PhoneNumber: m.PhoneNumber,
		PassportNo:  m.PassportNo,
		Password:    m.Password,
		TillNumber:  m.TillNumber,
	}

	// add new merchant to merchants table, if query return violation of unique key column,
	// we know that the merchant with given record exists and return that merchant instead
	result := r.db.Model(schema.Merchant{}).Create(&row)
	if err := result.Error; err != nil {
		// we check if the error is a postgres unique constraint violation
		if pgerr, ok := err.(*pgconn.PgError); ok && pgerr.Code == "23505" {
			return m, errors.Error{Code: errors.ECONFLICT, Message: errors.ErrUserExists}
		}
		return merchant.Merchant{}, errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}

	return merchant.Merchant{
		ID:          row.ID,
		Email:       row.Email,
		FirstName:   row.FirstName,
		LastName:    row.LastName,
		PhoneNumber: row.PhoneNumber,
		PassportNo:  row.PassportNo,
		Password:    row.Password,
		TillNumber:  row.TillNumber,
	}, nil
}

// Delete a merchant
func (r MerchantRepository) Delete(m merchant.Merchant) error {
	result := r.db.Delete(&schema.Merchant{ID: m.ID})
	if result.Error != nil {
		return errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}
	return nil
}

// FetchAll gets all merchants in db
func (r MerchantRepository) FetchAll() ([]merchant.Merchant, error) {
	var rows []schema.Merchant
	result := r.db.Find(&rows)
	if err := result.Error; err != nil {
		return nil, errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}

	// we might not need to return this error
	if result.RowsAffected == 0 {
		return nil, errors.Error{Code: errors.ENOTFOUND}
	}

	var merchants []merchant.Merchant
	for _, row := range rows {
		merchants = append(merchants, merchant.Merchant{
			ID:          row.ID,
			Email:       row.Email,
			FirstName:   row.FirstName,
			LastName:    row.LastName,
			PhoneNumber: row.PhoneNumber,
			PassportNo:  row.PassportNo,
			Password:    row.Password,
			TillNumber:  row.TillNumber,
		})
	}

	return merchants, nil
}

// FindByID searches merchant by primary id
func (r MerchantRepository) FindByID(id uuid.UUID) (merchant.Merchant, error) {
	m, err := r.searchBy(schema.Merchant{ID: id})
	return m, err
}

// FindByEmail searches merchant by email
func (r MerchantRepository) FindByEmail(email string) (merchant.Merchant, error) {
	m, err := r.searchBy(schema.Merchant{Email: email})
	return m, err
}

// Update
func (r MerchantRepository) Update(m merchant.Merchant) error {
	row := schema.Merchant{
		Email:       m.Email,
		FirstName:   m.FirstName,
		LastName:    m.LastName,
		PhoneNumber: m.PhoneNumber,
		PassportNo:  m.PassportNo,
		Password:    m.Password,
		TillNumber:  m.TillNumber,
	}

	result := r.db.Model(&schema.Merchant{}).Where(schema.Merchant{ID: m.ID}).Omit("id").Updates(row)
	if err := result.Error; err != nil {
		return errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}
	return nil
}
