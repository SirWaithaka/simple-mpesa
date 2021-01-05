package customer

import (
	"simple-mpesa/src/agent"
	"simple-mpesa/src/merchant"
	"simple-mpesa/src/models"
	"simple-mpesa/src/subscriber"
)

type CustomersFetcher interface {
	GetAllAgents() ([]models.Agent, error)
	GetAllMerchants() ([]models.Merchant, error)
	GetAllSubscribers() ([]models.Subscriber, error)
}

type customerFetcher struct {
	agentRepo agent.Repository
	merchRepo merchant.Repository
	subRepo   subscriber.Repository
}

func (fetcher customerFetcher) GetAllAgents() ([]models.Agent, error) {
	return fetcher.agentRepo.FetchAll()
}

func (fetcher customerFetcher) GetAllMerchants() ([]models.Merchant, error) {
	return fetcher.merchRepo.FetchAll()
}

func (fetcher customerFetcher) GetAllSubscribers() ([]models.Subscriber, error) {
	return fetcher.subRepo.FetchAll()
}
