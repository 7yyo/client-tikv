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

func ErrTableName(s ...string) Error {
	msg := fmt.Sprintf("Illegal table name: `%s`, must be `TiKV`.", s[0])
	return Error{message: msg}
}

func ErrInsertValueCount() Error {
	msg := fmt.Sprintf("Number of illegal insertions, kv should only have two values.")
	return Error{message: msg}
}
