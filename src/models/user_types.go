package models

type UserType string

// IsAgent returns true if the user type matches
// UserTypAgent or UserTypSuperAgent
func (typ UserType) IsAgent() bool {
	return typ == UserTypAgent || typ == UserTypSuperAgent
}

const (
	UserTypAdmin      = UserType("administrator")
	UserTypAgent      = UserType("agent")
	UserTypMerchant   = UserType("merchant")
	UserTypSubscriber = UserType("subscriber")

	UserTypSuperAgent = UserType("super_agent")
)
