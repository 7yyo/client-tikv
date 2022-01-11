package pd

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

func ErrGetPlacementInfoFailed(s ...string) Error {
	msg := fmt.Sprintf("Get placement driver failed, error: %s", s[0])
	return Error{message: msg}
}
