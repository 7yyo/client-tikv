package tikv

import (
	"github.com/tikv/client-go/v2/tikv"
	"tikv-client/util"
	"time"
)

type insertExecutor struct {
	table
	kvPairs []kvPair
	client  *tikv.RawKVClient
}

func (ie insertExecutor) Valid() error {
	if err := checkTblName(ie); err != nil {
		return err
	}
	return nil
}

func (ie insertExecutor) Exec() error {
	t := time.Now()
	if ie.kvPairs == nil {
		return errInsertValueCount()
	}
	for _, kv := range ie.kvPairs {
		if err := ie.client.Put(kv.key, kv.value); err != nil {
			return err
		}
	}
	util.QueryOkNRows(len(ie.kvPairs), t)
	return nil
}
