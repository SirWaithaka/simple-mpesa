package repositories

import (
	"simple-mpesa/src/errors"
	"simple-mpesa/src/storage"
	"simple-mpesa/src/subscriber"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

func NewSubscriberRepository(database *storage.Database) *SubscriberRepository {
	return &SubscriberRepository{db: database}
}

type SubscriberRepository struct {
	db *storage.Database
}

func (r SubscriberRepository) searchBy(row subscriber.Subscriber) (subscriber.Subscriber, error) {
	var s subscriber.Subscriber
	result := r.db.Where(row).First(&s)
	// check if no record found.
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return subscriber.Subscriber{}, errors.Error{Code: errors.ENOTFOUND}
	}
	if err := result.Error; err != nil {
		return subscriber.Subscriber{}, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	return s, nil
}

// Add a subscriber if already not in db.
func (r SubscriberRepository) Add(s subscriber.Subscriber) (subscriber.Subscriber, error) {
	// add new subscriber to subscribers table, if query return violation of unique key column,
	// we know that the subscriber with given record exists and return that subscriber instead
	result := r.db.Model(subscriber.Subscriber{}).Create(&s)
	if err := result.Error; err != nil {
		// we check if the error is a postgres unique constraint violation
		if pgerr, ok := err.(*pgconn.PgError); ok && pgerr.Code == "23505" {
			return s, errors.Error{Code: errors.ECONFLICT, Message: errors.ErrUserExists}
		}
		return subscriber.Subscriber{}, errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}

	return s, nil
}

// Delete a subscriber
func (r SubscriberRepository) Delete(s subscriber.Subscriber) error {
	result := r.db.Delete(&s)
	if result.Error != nil {
		return errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}
	return nil
}

// FetchAll gets all subscribers in db
func (r SubscriberRepository) FetchAll() ([]subscriber.Subscriber, error) {
	var subscribers []subscriber.Subscriber
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
func (r SubscriberRepository) FindByID(id uuid.UUID) (subscriber.Subscriber, error) {
	s, err := r.searchBy(subscriber.Subscriber{ID: id})
	return s, err
}

// FindByEmail searches subscriber by email
func (r SubscriberRepository) FindByEmail(email string) (subscriber.Subscriber, error) {
	s, err := r.searchBy(subscriber.Subscriber{Email: email})
	return s, err
}

// Update
func (r SubscriberRepository) Update(s subscriber.Subscriber) error {
	result := r.db.Model(&subscriber.Subscriber{}).Omit("id").Updates(s)
	if err := result.Error; err != nil {
		return errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}
	return nil
}
