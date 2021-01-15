package administrator

import (
	"simple-mpesa/src"
	"simple-mpesa/src/domain/account"
	"simple-mpesa/src/domain/customer"
	"simple-mpesa/src/domain/value_objects"
	"simple-mpesa/src/errors"
	"simple-mpesa/src/helpers"
)

type Interactor interface {
	AuthenticateByEmail(email, password string) (Admin, error)
	Register(RegistrationParams) (Admin, error)
	AssignFloat(AssignFloatParams) (float64, error)
}

func NewInteractor(config src.Config, adminsRepo Repository, accountant account.Accountant, finder customer.Finder) Interactor {
	return &interactor{
		config:         config,
		repository:     adminsRepo,
		accountant:     accountant,
		customerFinder: finder,
	}
}

type interactor struct {
	accountant     account.Accountant
	customerFinder customer.Finder
	config         src.Config
	repository     Repository
}

// AuthenticateByEmail verifies a admin by the provided unique email address
func (i interactor) AuthenticateByEmail(email, password string) (Admin, error) {
	// search for admin by email.
	admin, err := i.repository.GetByEmail(email)
	if errors.ErrorCode(err) == errors.ENOTFOUND {
		return Admin{}, errors.Error{Err: err, Message: errors.ErrUserNotFound}
	} else if err != nil {
		return Admin{}, err
	}

	// validate password
	if err := helpers.ComparePasswordToHash(admin.Password, password); err != nil {
		return Admin{}, errors.Unauthorized{Message: errors.ErrorMessage(err)}
	}

	return admin, nil
}

// Register takes in a admin object and adds the admin to db.
func (i interactor) Register(params RegistrationParams) (Admin, error) {
	admin := Admin{
		FirstName: params.FirstName,
		LastName:  params.LastName,
		Email:     params.Email,
		Password:  params.Password,
	}

	// hash admin password before adding to db.
	passwordHash, err := helpers.HashPassword(admin.Password)
	if err != nil { // if we get an error, it means our hashing func dint work
		return Admin{}, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	// change password to hashed string
	admin.Password = passwordHash
	adm, err := i.repository.Add(admin)
	if err != nil {
		return Admin{}, err
	}

	return adm, nil
}

// AssignFloat is an admin only operation that gives a super agent the initial amount of
// money. It can also be used in subsequent operations to increase the amount of money in
// the system.
func (i interactor) AssignFloat(params AssignFloatParams) (float64, error) {
	agent, err := i.customerFinder.FindAgentByEmail(params.AgentAccountNumber)
	if err != nil {
		return 0, err
	}

	// float is only assignable to a super agent
	if !agent.IsSuperAgent() {
		return 0, errors.Error{Code: errors.EINVALID, Message: errors.ErrAgentNotSuperAgent}
	}

	balance, err := i.accountant.CreditAccount(agent.ID, params.Amount.ToCents(), value_objects.TxnFloatAssignment)
	if err != nil {
		return 0, err
	}

	return balance, nil
}
