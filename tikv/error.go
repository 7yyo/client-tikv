package tikv

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

func errTableName(s ...string) Error {
	msg := fmt.Sprintf("Illegal table name: `%s`", s[0])
	return Error{message: msg}
}

func errUnsupportedTiKVScan() Error {
	msg := fmt.Sprintf("Sorry for not currently supporting querying all data in tikv.")
	return Error{message: msg}
}

func errUnsupportedTiKVDelete() Error {
	msg := fmt.Sprintf("Sorry for not currently supporting delete all data in tikv.")
	return Error{message: msg}
}

func errUnsupportedSQL(s ...string) Error {
	msg := fmt.Sprintf("Unsupported SQL `%s`", s[0])
	return Error{message: msg}
}

func errQueryScan() Error {
	return Error{message: "The query range must be a closed interval."}
}

func errUnknownColumn(s ...string) Error {
	msg := fmt.Sprintf("Unknown column `%s` in 'field list'.", s[0])
	return Error{message: msg}
}

func errParseJson(s ...string) Error {
	msg := fmt.Sprintf("Failed to parse json, please confirm whether it is standard json, error: (%s), value: %s.", s[0], s[1])
	return Error{message: msg}
}

func errInsertValueCount() Error {
	msg := fmt.Sprintf("Number of illegal insertions, kv should only have two values.")
	return Error{message: msg}
}

func errColumnDontMatch() Error {
	msg := fmt.Sprintf("Column count doesn't match value count.")
	return Error{message: msg}
}

func errOrderbyCount() Error {
	msg := fmt.Sprintf("Only supports single-column sorting.")
	return Error{message: msg}
}

func errScanLimit() error {
	msg := fmt.Sprintf("Limit should be less than 10240.")
	return Error{message: msg}
}
