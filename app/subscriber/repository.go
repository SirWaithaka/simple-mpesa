package subscriber

import (
	"simple-mpesa/app/errors"
	"simple-mpesa/app/models"
	"simple-mpesa/app/storage"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

type Repository interface {
	Add(models.Subscriber) (models.Subscriber, error)
	Delete(models.Subscriber) error
	FetchAll() ([]models.Subscriber, error)
	FindByID(uuid.UUID) (models.Subscriber, error)
	FindByEmail(string) (models.Subscriber, error)
	Update(models.Subscriber) error
}

func NewRepository(database *storage.Database) Repository {
	return &repository{db: database}
}

type repository struct {
	db *storage.Database
}

func (r repository) searchBy(row models.Subscriber) (models.Subscriber, error) {
	var subscriber models.Subscriber
	result := r.db.Where(row).First(&subscriber)
	// check if no record found.
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return models.Subscriber{}, errors.Error{Code: errors.ENOTFOUND}
	}
	if err := result.Error; err != nil {
		return models.Subscriber{}, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	return subscriber, nil
}

// Add a subscriber if already not in db.
func (r repository) Add(subscriber models.Subscriber) (models.Subscriber, error) {
	// add new subscriber to subscribers table, if query return violation of unique key column,
	// we know that the subscriber with given record exists and return that subscriber instead
	result := r.db.Model(models.Subscriber{}).Create(&subscriber)
	if err := result.Error; err != nil {
		// we check if the error is a postgres unique constraint violation
		if pgerr, ok := err.(*pgconn.PgError); ok && pgerr.Code == "23505" {
			return subscriber, errors.Error{Code: errors.ECONFLICT, Message: errors.ErrUserExists}
		}
		return models.Subscriber{}, errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}

	return subscriber, nil
}

// Delete a subscriber
func (r repository) Delete(subscriber models.Subscriber) error {
	result := r.db.Delete(&subscriber)
	if result.Error != nil {
		return errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}
	return nil
}

// FetchAll gets all subscribers in db
func (r repository) FetchAll() ([]models.Subscriber, error) {
	var subscribers []models.Subscriber
	result := r.db.Find(&subscribers)
	if err := result.Error; err != nil {
		return nil, errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}

	// we might not need to return this error
	if result.RowsAffected == 0 {
		return nil, errors.Error{Code: errors.ENOTFOUND}
	}

	return subscribers, nil
}

// FindByID searches subscriber by primary id
func (r repository) FindByID(id uuid.UUID) (models.Subscriber, error) {
	subscriber, err := r.searchBy(models.Subscriber{ID: id})
	return subscriber, err
}

// FindByEmail searches subscriber by email
func (r repository) FindByEmail(email string) (models.Subscriber, error) {
	subscriber, err := r.searchBy(models.Subscriber{Email: email})
	return subscriber, err
}

// Update
func (r repository) Update(subscriber models.Subscriber) error {
	var sub models.Subscriber
	result := r.db.Model(&sub).Omit("id").Updates(subscriber)
	if err := result.Error; err != nil {
		return errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}
	return nil
}
