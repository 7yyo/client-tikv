package error

import (
	"fmt"
)

type Error struct {
	message string
}

func (e Error) Error() string {
	return fmt.Sprintf("%s", e.message)
}

func ErrUnknownColumn(s ...string) Error {
	msg := fmt.Sprintf("Unknown column %s in 'field list'", s[0])
	return Error{message: msg}
}

func ErrParseJson(s ...string) Error {
	msg := fmt.Sprintf("Failed to parse json, please confirm whether it is standard json, error: %s, value: %s", s[0], s[1])
	return Error{message: msg}
}
