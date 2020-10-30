package agent

import (
	"simple-wallet/app/data"
	"simple-wallet/app/models"
)

func parseToNewAgent(agent models.Agent) data.CustomerContract {
	return data.CustomerContract{
		UserID: agent.ID,
	}
}
