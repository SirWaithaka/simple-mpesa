package errors

const (
	// ErrUserExists returned when adding a user with
	// phone number or email number that are already in the db.
	ErrUserExists         = ERMessage("user already exists")
	ErrUserNotFound       = ERMessage("user not found")
	ErrAgentNotSuperAgent = ERMessage("given agent is not a super agent")
)

// PasswordHashError
type PasswordHashError struct {
	Err error
}

func (e PasswordHashError) Error() string {
	return "error hashing password"
}

func (e PasswordHashError) Debug() error {
	return e.Err
}
