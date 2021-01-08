package agent

import (
	"simple-mpesa/src/data"
)

func parseToNewAgent(agent Agent) data.CustomerContract {
	return data.CustomerContract{
		UserID: agent.ID,
	}
}
