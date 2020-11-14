package models

type UserType string

const (
	UserTypAdmin      = UserType("administrator")
	UserTypAgent      = UserType("agent")
	UserTypMerchant   = UserType("merchant")
	UserTypSubscriber = UserType("subscriber")
)
