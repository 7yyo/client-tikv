package tikv

import (
	"fmt"
	"github.com/tikv/client-go/v2/tikv"
	"tikv-client/util"
	"time"
)

type truncateExecutor struct {
	table
	client *tikv.RawKVClient
}

func (te truncateExecutor) Valid() error {
	if err := checkTblName(te); err != nil {
		return err
	}
	return nil
}

// Exec TODO deleteRange unsupported truncate tikv when key, value is nil
func (te truncateExecutor) Exec() error {
	t := time.Now()
	var startKey []byte
	var endKey []byte
	if err := te.client.DeleteRange(startKey, endKey); err != nil {
		return err
	}
	fmt.Println(util.Red("Only the syntax takes effect."))
	util.QueryOkNRows(0, t)
	return nil
}
