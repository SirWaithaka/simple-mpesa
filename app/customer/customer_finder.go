package customer

import (
	"simple-mpesa/app/agent"
	"simple-mpesa/app/merchant"
	"simple-mpesa/app/models"
	"simple-mpesa/app/subscriber"
)

type Finder interface {
	FindAgentByEmail(string) (models.Agent, error)
	FindMerchantByEmail(string) (models.Merchant, error)
	FindSubscriberByEmail(string) (models.Subscriber, error)
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
	return f.agentRepo.FindByEmail(email)
}

func (f finder) FindMerchantByEmail(email string) (models.Merchant, error) {
	return f.merchRepo.FindByEmail(email)
}

func (f finder) FindSubscriberByEmail(email string) (models.Subscriber, error) {
	return f.subRepo.FindByEmail(email)
}
