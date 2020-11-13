package ports

import (
	"simple-mpesa/app/agent"
	"simple-mpesa/app/merchant"
	"simple-mpesa/app/models"
	"simple-mpesa/app/subscriber"
)

type CustomerFinder interface {
	FindAgentByEmail(string) (models.Agent, error)
	FindMerchantByEmail(string) (models.Merchant, error)
	FindSubscriberByEmail(string) (models.Subscriber, error)
}

func NewCustomerFinder(agentRepo agent.Repository, merchRepo merchant.Repository, subRepo subscriber.Repository) CustomerFinder {
	return &customerFinder{
		agentRepo: agentRepo,
		merchRepo: merchRepo,
		subRepo:   subRepo,
	}
}

type customerFinder struct {
	agentRepo agent.Repository
	merchRepo merchant.Repository
	subRepo   subscriber.Repository
}

func (finder customerFinder) FindAgentByEmail(email string) (models.Agent, error) {
	return finder.agentRepo.FindByEmail(email)
}

func (finder customerFinder) FindMerchantByEmail(email string) (models.Merchant, error) {
	return finder.merchRepo.FindByEmail(email)
}

func (finder customerFinder) FindSubscriberByEmail(email string) (models.Subscriber, error) {
	return finder.subRepo.FindByEmail(email)
}
