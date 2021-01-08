package customer

import (
	"simple-mpesa/src/agent"
	"simple-mpesa/src/errors"
	"simple-mpesa/src/merchant"
	"simple-mpesa/src/value_objects"
	"simple-mpesa/src/subscriber"

	"github.com/gofrs/uuid"
)

type Finder interface {
	FindAgentByEmail(string) (agent.Agent, error)
	FindMerchantByEmail(string) (merchant.Merchant, error)
	FindSubscriberByEmail(string) (subscriber.Subscriber, error)
	FindIDByEmail(string, value_objects.UserType) (uuid.UUID, error)
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

func (f finder) FindAgentByEmail(email string) (agent.Agent, error) {
	agt, err := f.agentRepo.FindByEmail(email)
	if errors.ErrorCode(err) == errors.ENOTFOUND {
		return agent.Agent{}, errors.Error{Err: err, Message: errors.ErrUserNotFound}
	} else if err != nil {
		return agent.Agent{}, err
	}

	return agt, nil
}

func (f finder) FindMerchantByEmail(email string) (merchant.Merchant, error) {
	merch, err := f.merchRepo.FindByEmail(email)
	if errors.ErrorCode(err) == errors.ENOTFOUND {
		return merchant.Merchant{}, errors.Error{Err: err, Message: errors.ErrUserNotFound}
	} else if err != nil {
		return merchant.Merchant{}, err
	}

	return merch, nil
}

func (f finder) FindSubscriberByEmail(email string) (subscriber.Subscriber, error) {
	sub, err := f.subRepo.FindByEmail(email)
	if errors.ErrorCode(err) == errors.ENOTFOUND {
		return subscriber.Subscriber{}, errors.Error{Err: err, Message: errors.ErrUserNotFound}
	} else if err != nil {
		return subscriber.Subscriber{}, err
	}

	return sub, nil
}

func (f finder) FindIDByEmail(email string, userType value_objects.UserType) (uuid.UUID, error) {
	switch userType {
	case value_objects.UserTypAgent:
		agt, err := f.FindAgentByEmail(email)
		if err != nil {
			return uuid.Nil, err
		}
		return agt.ID, nil
	case value_objects.UserTypMerchant:
		merch, err := f.FindMerchantByEmail(email)
		if err != nil {
			return uuid.Nil, err
		}
		return merch.ID, nil
	case value_objects.UserTypSubscriber:
		sub, err := f.FindSubscriberByEmail(email)
		if err != nil {
			return uuid.Nil, err
		}
		return sub.ID, nil
	}
	return uuid.Nil, errors.Error{Code: errors.EINVALID, Message: errors.ErrUserNotFound}
}