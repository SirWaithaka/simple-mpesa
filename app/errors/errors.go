package errors

import (
	"bytes"
	"errors"
	"fmt"
)

// ERCode defines a machine readable error code type
type ERCode string

// ErrorMessage defines a string type of an error description
type ErrorMessage string

const (
	ECONFLICT  = ERCode("conflict")        // action cannot be performed e.g. when inserting existing record to db
	EINTERNAL  = ERCode("internal")        // internal error
	EINVALID   = ERCode("invalid")         // validation failed
	ENOTFOUND  = ERCode("not_found")       // entity does not exist
)

// Error is our standard application error
type Error struct {
	// 	Machine readable error code
	Code ERCode

	// 	Human readable message
	Message ErrorMessage

	// nested error
	Err error
}

func (e Error) Error() string {
	var buf bytes.Buffer

	// 	If wrapping an error, print its Error() message.
	// otherwise print the error code & message.
	if e.Err != nil {
		buf.WriteString(e.Err.Error())
	} else {
		if e.Code != "" {
			fmt.Fprintf(&buf, "<%s> ", e.Code)
		}
		buf.WriteString(string(e.Message))
	}
	return buf.String()
}

// ErrorCode returns the code of the root error, if available. Otherwise returns EINTERNAL.
func ErrorCode(err error) ERCode {
	if err == nil {
		return ""
	} else if e, ok := err.(Error); ok && e.Code != "" {
		return e.Code
	} else if ok && e.Err != nil {
		return ErrorCode(err)
	}
	return EINTERNAL
}

type ErrorT struct {
	Message string
	Err     error
}

// Is reports whether any error in err's chain matches target.
func Is(err, target error) bool {
	return errors.Is(err, target)
}
