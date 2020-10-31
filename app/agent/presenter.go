package agent

import (
	"simple-mpesa/app/data"
	"simple-mpesa/app/models"
)

func parseToNewAgent(agent models.Agent) data.CustomerContract {
	return data.CustomerContract{
		UserID: agent.ID,
	}
}
