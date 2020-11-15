package agent

import (
	"log"

	"simple-mpesa/app"
	"simple-mpesa/app/data"
	"simple-mpesa/app/errors"
	"simple-mpesa/app/helpers"
	"simple-mpesa/app/models"
)

type Interactor interface {
	AuthenticateByEmail(email, password string) (models.Agent, error)
	Register(RegistrationParams) (models.Agent, error)
	UpdateSuperAgentStatus(email string) error
}

func NewInteractor(config app.Config, agentRepo Repository, custChan data.ChanNewCustomers) Interactor {
	return &interactor{
		config:           config,
		repository:       agentRepo,
		customersChannel: custChan,
	}
}

type interactor struct {
	customersChannel data.ChanNewCustomers
	config           app.Config
	repository       Repository
}

// AuthenticateByEmail verifies an agent by the provided unique email address
func (ui interactor) AuthenticateByEmail(email, password string) (models.Agent, error) {
	// search for agent by email.
	agent, err := ui.repository.FindByEmail(email)
	if errors.ErrorCode(err) == errors.ENOTFOUND {
		return models.Agent{}, errors.Error{Err: err, Message: errors.ErrUserNotFound}
	} else if err != nil {
		return models.Agent{}, err
	}

	// validate password
	if err := helpers.ComparePasswordToHash(agent.Password, password); err != nil {
		return models.Agent{}, errors.Unauthorized{Message: errors.ErrorMessage(err)}
	}

	return agent, nil
}

// Register takes in a agent object and adds the agent to db.
func (ui interactor) Register(params RegistrationParams) (models.Agent, error) {
	agent := models.Agent{
		FirstName:   params.FirstName,
		LastName:    params.LastName,
		Email:       params.Email,
		PhoneNumber: params.PhoneNumber,
		Password:    params.Password,
		PassportNo:  params.PassportNo,
	}

	// hash agent password before adding to db.
	passwordHash, err := helpers.HashPassword(agent.Password)
	if err != nil { // if we get an error, it means our hashing func dint work
		return models.Agent{}, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	// change password to hashed string
	agent.Password = passwordHash
	agt, err := ui.repository.Add(agent)
	if err != nil {
		return models.Agent{}, err
	}

	// tell channel listeners that a new agent has been created.
	ui.postNewAgentToChannel(&agt)
	return agt, nil
}

func (ui interactor) UpdateSuperAgentStatus(email string) error {
	agent, err := ui.repository.FindByEmail(email)
	if errors.ErrorCode(err) == errors.ENOTFOUND {
		return errors.Error{Err: err, Message: errors.ErrUserNotFound}
	} else if err != nil {
		return err
	}

	log.Printf("status before: %v", agent.SuperAgent)
	agent.SuperAgent = agent.SuperAgent.Not()
	log.Printf("status after: %v", agent.SuperAgent)
	err = ui.repository.Update(agent)
	if err != nil {
		return err
	}

	return nil
}

// take the newly created agent and post them to channel
// that listens for newly created customers and acts upon them
// like creating an account for them automatically.
func (ui interactor) postNewAgentToChannel(agent *models.Agent) {
	newAgent := parseToNewAgent(*agent)
	go func() { ui.customersChannel.Writer <- newAgent }()
}
