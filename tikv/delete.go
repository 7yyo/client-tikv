package tikv

import (
	"github.com/tikv/client-go/v2/tikv"
	"tikv-client/pd"
	"tikv-client/util"
	"time"
)

type deleteExecutor struct {
	sql             string
	placementDriver pd.PlacementDriverGroup
	table
	kvPairs []kvPair
	client  *tikv.RawKVClient
}

func (d deleteExecutor) Exec() error {
	t := time.Now()
	var rows int
	var err error
	keys := getKeys(d.kvPairs)
	if len(d.kvPairs) == 1 {
		r, err := d.client.Get(keys[0])
		if err != nil {
			return err
		}
		if r == nil {
			util.NRowsAffected(0, t)
			return nil
		}
		err = d.client.Delete(d.kvPairs[0].key)
		rows = 1
	} else {
		rs, err := d.client.BatchGet(keys)
		if err != nil {
			return err
		}
		var ks [][]byte
		for i, r := range rs {
			if r != nil {
				ks = append(ks, keys[i])
				rows++
			}
		}
		err = d.client.BatchDelete(ks)
	}
	if err != nil {
		return err
	}
	util.NRowsAffected(rows, t)
	return nil
}

func (d deleteExecutor) Valid() error {
	if d.table.name == "" {
		return errUnsupportedSQL(d.sql)
	}
	if err := checkTblName(d); err != nil {
		return err
	}
	if len(d.kvPairs) == 0 {
		return errUnsupportedTiKVDelete()
	}
	return nil
}
