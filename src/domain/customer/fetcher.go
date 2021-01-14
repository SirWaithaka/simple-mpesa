package customer

import (
	"simple-mpesa/src/domain/agent"
	"simple-mpesa/src/domain/merchant"
	"simple-mpesa/src/domain/subscriber"
)

type CustomersFetcher interface {
	GetAllAgents() ([]agent.Agent, error)
	GetAllMerchants() ([]merchant.Merchant, error)
	GetAllSubscribers() ([]subscriber.Subscriber, error)
}

type customerFetcher struct {
	agentRepo agent.Repository
	merchRepo merchant.Repository
	subRepo   subscriber.Repository
}

func (fetcher customerFetcher) GetAllAgents() ([]agent.Agent, error) {
	return fetcher.agentRepo.FetchAll()
}

func (fetcher customerFetcher) GetAllMerchants() ([]merchant.Merchant, error) {
	return fetcher.merchRepo.FetchAll()
}

func (fetcher customerFetcher) GetAllSubscribers() ([]subscriber.Subscriber, error) {
	return fetcher.subRepo.FetchAll()
}
