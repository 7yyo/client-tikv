package cdc

import (
	"fmt"
	"tikv-client/util"
)

type Error struct {
	message string
}

func (e Error) Error() string {
	return util.Red(fmt.Sprintf("%s", e.message))
}

func errInvalidCommand(s []string) error {
	msg := fmt.Sprintf("Invalid command: %s", s)
	return Error{message: msg}
}

func errCanNotFind(s ...string) error {
	msg := fmt.Sprintf("Can not find `%s` for %s", s[0], s[1])
	return Error{message: msg}
}
