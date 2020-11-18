package transaction_test

import (
	"testing"

	"simple-mpesa/app/models"
	"simple-mpesa/app/tariff"
	"simple-mpesa/app/transaction"

	"github.com/gofrs/uuid"
)

type testAccountant struct {
}

func (ta testAccountant) DebitAccount(userID uuid.UUID, amount models.Cents, reason models.TxnOperation) (float64, error) {
	return 0, nil
}

func (ta testAccountant) CreditAccount(userID uuid.UUID, amount models.Cents, reason models.TxnOperation) (float64, error) {
	return 0, nil
}

type testTariffManager struct {
}

func (tm testTariffManager) GetCharge(operation models.TxnOperation, src models.UserType, dest models.UserType) (models.Cents, error) {
	return 0, nil
}

func (tm testTariffManager) GetTariff() ([]tariff.Charge, error) {
	return nil, nil
}

func (tm testTariffManager) UpdateCharge(chargeID uuid.UUID, fee models.Cents) error {
	return nil
}

func TestTransactor_Transact(t *testing.T) {
	// 	we have some business rules we need to test that are enforced

	type testcase struct {
		name     string
		input    transaction.Transaction
	}

	agent1 := models.TxnCustomer{
		UserID:   uuid.FromStringOrNil("777b2533-f571-4e70-a366-7e69e9413ebf"),
		UserType: models.UserTypAgent,
	}

	subscriber1 := models.TxnCustomer{
		UserID:   uuid.FromStringOrNil("a6b8847a-0f18-45f9-98af-5b9f039076bd"),
		UserType: models.UserTypSubscriber,
	}

	testcases := []testcase{
		// we begin with the first rule
		{
			name: "cannot transact with own account: transaction with identical source and destination user id",
			input: transaction.Transaction{
				Source:       subscriber1,
				Destination:  subscriber1,
				TxnOperation: models.TxnOpTransfer,
				Amount:       models.Shillings(300),
			},
		},
		// we test now rules regarding deposits
		{
			name: "can only deposit at an agent",
			input: transaction.Transaction{
				Source:       subscriber1,
				Destination:  agent1,
				TxnOperation: models.TxnOpDeposit,
				Amount:       models.Shillings(300),
			},
		},
	}

	transactor := transaction.NewTransactor(testAccountant{}, testTariffManager{})

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := transactor.Transact(tc.input)
			if err == nil {
				t.Fatal("expected an error")
			}
		})
	}

}
