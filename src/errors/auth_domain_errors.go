package errors


const (
	InvalidCredentials = ERMessage("provided credentials are invalid")
)

// Unauthorized
type Unauthorized struct {
	Message string
}

func (e Unauthorized) Error() string {
	return e.Message
}
