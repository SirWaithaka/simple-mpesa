package agent

import (
	"simple-mpesa/src/data"
	"simple-mpesa/src/models"
)

func parseToNewAgent(agent models.Agent) data.CustomerContract {
	return data.CustomerContract{
		UserID: agent.ID,
	}
}
