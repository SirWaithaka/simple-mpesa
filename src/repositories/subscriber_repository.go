package repositories

import (
	"simple-mpesa/src/errors"
	"simple-mpesa/src/models"
	"simple-mpesa/src/storage"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

func NewSubscriberRepository(database *storage.Database) *Subscriber {
	return &Subscriber{db: database}
}

type Subscriber struct {
	db *storage.Database
}

func (r Subscriber) searchBy(row models.Subscriber) (models.Subscriber, error) {
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
func (r Subscriber) Add(subscriber models.Subscriber) (models.Subscriber, error) {
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
func (r Subscriber) Delete(subscriber models.Subscriber) error {
	result := r.db.Delete(&subscriber)
	if result.Error != nil {
		return errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}
	return nil
}

// FetchAll gets all subscribers in db
func (r Subscriber) FetchAll() ([]models.Subscriber, error) {
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
func (r Subscriber) FindByID(id uuid.UUID) (models.Subscriber, error) {
	subscriber, err := r.searchBy(models.Subscriber{ID: id})
	return subscriber, err
}

// FindByEmail searches subscriber by email
func (r Subscriber) FindByEmail(email string) (models.Subscriber, error) {
	subscriber, err := r.searchBy(models.Subscriber{Email: email})
	return subscriber, err
}

// Update
func (r Subscriber) Update(subscriber models.Subscriber) error {
	var sub models.Subscriber
	result := r.db.Model(&sub).Omit("id").Updates(subscriber)
	if err := result.Error; err != nil {
		return errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}
	return nil
}

