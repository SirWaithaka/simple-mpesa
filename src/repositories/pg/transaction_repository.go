package pg

// type TransactionRepository struct {
// 	database *storage.Database
// }
//
// func NewTransactionRepository(db *storage.Database) *TransactionRepository {
// 	return &TransactionRepository{db}
// }
//
// func (r TransactionRepository) Add(tx transaction.Statement) (transaction.Statement, error) {
// 	result := r.database.Create(&tx)
// 	if err := result.Error; err != nil {
// 		return transaction.Statement{}, errors.Error{Err: err, Code: errors.EINTERNAL}
// 	}
//
// 	return tx, nil
// }
