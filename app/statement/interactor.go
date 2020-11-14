package statement

import (
	"time"

	"github.com/gofrs/uuid"
)

const miniStatementCount = uint(5)

type Interactor interface {
	GetStatement(userId uuid.UUID) ([]Statement, error)
}

type interactor struct {
	repository Repository
}

func NewInteractor(repository Repository) Interactor {
	return &interactor{repository}
}

func (i interactor) GetStatement(userID uuid.UUID) ([]Statement, error) {
	now := time.Now()
	transactions, err := i.repository.GetStatements(userID, now, miniStatementCount)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}
