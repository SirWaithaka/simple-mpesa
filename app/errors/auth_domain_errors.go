package errors

/*
	Auth domain errors
*/

// ErrTokenParsing
type ErrTokenParsing ErrorT

// Unauthorized
type Unauthorized struct {
	Message string
}

func (e Unauthorized) Error() string {
	return e.Message
}