package customer

import (
	"simple-mpesa/app/agent"
	"simple-mpesa/app/errors"
	"simple-mpesa/app/merchant"
	"simple-mpesa/app/models"
	"simple-mpesa/app/subscriber"

	"github.com/gofrs/uuid"
)

type Finder interface {
	FindAgentByEmail(string) (models.Agent, error)
	FindMerchantByEmail(string) (models.Merchant, error)
	FindSubscriberByEmail(string) (models.Subscriber, error)
	FindIDByEmail(string, models.UserType) (uuid.UUID, error)
}

func NewFinder(agentRepo agent.Repository, merchRepo merchant.Repository, subRepo subscriber.Repository) Finder {
	return &finder{
		agentRepo: agentRepo,
		merchRepo: merchRepo,
		subRepo:   subRepo,
	}
}

type finder struct {
	agentRepo agent.Repository
	merchRepo merchant.Repository
	subRepo   subscriber.Repository
}

func (f finder) FindAgentByEmail(email string) (models.Agent, error) {
	agt, err := f.agentRepo.FindByEmail(email)
	if errors.ErrorCode(err) == errors.ENOTFOUND {
		return models.Agent{}, errors.Error{Err: err, Message: errors.ErrUserNotFound}
	} else if err != nil {
		return models.Agent{}, err
	}

	return agt, nil
}

func (f finder) FindMerchantByEmail(email string) (models.Merchant, error) {
	merch, err := f.merchRepo.FindByEmail(email)
	if errors.ErrorCode(err) == errors.ENOTFOUND {
		return models.Merchant{}, errors.Error{Err: err, Message: errors.ErrUserNotFound}
	} else if err != nil {
		return models.Merchant{}, err
	}

	return merch, nil
}

func (f finder) FindSubscriberByEmail(email string) (models.Subscriber, error) {
	sub, err := f.subRepo.FindByEmail(email)
	if errors.ErrorCode(err) == errors.ENOTFOUND {
		return models.Subscriber{}, errors.Error{Err: err, Message: errors.ErrUserNotFound}
	} else if err != nil {
		return models.Subscriber{}, err
	}

	return sub, nil
}

func (f finder) FindIDByEmail(email string, userType models.UserType) (uuid.UUID, error) {
	switch userType {
	case models.UserTypAgent:
		agt, err := f.FindAgentByEmail(email)
		if err != nil {
			return uuid.Nil, err
		}
		return agt.ID, nil
	case models.UserTypMerchant:
		merch, err := f.FindMerchantByEmail(email)
		if err != nil {
			return uuid.Nil, err
		}
		return merch.ID, nil
	case models.UserTypSubscriber:
		sub, err := f.FindSubscriberByEmail(email)
		if err != nil {
			return uuid.Nil, err
		}
		return sub.ID, nil
	}
	return uuid.Nil, errors.Error{Code: errors.EINVALID, Message: errors.ErrUserNotFound}
}