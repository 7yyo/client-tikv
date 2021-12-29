package tikv

import (
	"encoding/json"
	"flag"
	"github.com/pingcap/log"
	"github.com/tikv/client-go/v2/config"
	"github.com/tikv/client-go/v2/tikv"
	p "github.com/tikv/pd/client"
	"go.uber.org/zap/zapcore"
	"strings"
	e "tikv-client/error"
	"tikv-client/syntax"
	"tikv-client/util"
	"time"
)

func ParseArgs() []string {
	var pd string
	flag.StringVar(&pd, "pd", "", "pd endpoints")
	flag.Parse()
	return strings.Split(pd, ",")
}

func NewKvClient(pdEndPoint []string) (*tikv.RawKVClient, error) {

	// In order to don't print log from client-go, set log level = panic.
	setLogLevel(zapcore.PanicLevel)

	c, err := tikv.NewRawKVClient(pdEndPoint, config.DefaultConfig().Security, p.WithMaxErrorRetry(1))
	if err != nil {
		return nil, err
	}
	return c, nil
}

func setLogLevel(l zapcore.Level) {
	log.SetLevel(l)
}

func (c *Completer) Put(sql *syntax.SQL) (string, error) {
	t := time.Now()
	if sql.KvPairs == nil {
		return "", e.ErrInsertValueCount()
	}
	for _, kv := range sql.KvPairs {
		err := c.client.Put([]byte(kv.Key), []byte(kv.Value))
		if err != nil {
			return "", err
		}
	}
	return util.QueryOkNRows(len(sql.KvPairs), t), nil
}

// Get or batchGet
func (c *Completer) Get(sql *syntax.SQL) (string, error) {
	if syntax.IsSearchKvPairs(sql.Fields) {
		return c.getKv2Pairs(sql)
	} else {
		return c.getKv2Field(sql)
	}
}

// Get kv pairs
func (c *Completer) getKv2Pairs(sql *syntax.SQL) (string, error) {
	t := time.Now()

	tbl := util.NewKvDisplayTable()
	keys := syntax.GetKeys(sql.KvPairs)

	rs, err := c.client.BatchGet(keys)
	if err != nil {
		return "", err
	}

	for i, kv := range sql.KvPairs {
		kv.Value = string(rs[i])
		if kv.Value != "" {
			util.NewKvTableRow(tbl, kv)
		}
	}

	if tbl.Length() == 0 {
		return util.EmptyResult(), nil
	}
	tbl.Render()

	return util.QueryOkNRows(tbl.Length(), t), nil
}

// Get fields from kv by json
func (c *Completer) getKv2Field(sql *syntax.SQL) (string, error) {

	t := time.Now()

	var fieldNames []interface{}
	for _, field := range sql.Fields {
		fieldNames = append(fieldNames, field.Text())
	}

	tbl := util.NewNormalDisplayTable(fieldNames)
	m := make(map[string]string)
	for _, kv := range sql.KvPairs {
		value, err := c.client.Get([]byte(kv.Key))
		if err != nil {
			return "", err
		}
		if value == nil {
			continue
		}
		err = json.Unmarshal(value, &m)
		if err != nil {
			return "", e.ErrParseJson(err.Error(), string(value))
		}
		err = syntax.HasField(m, sql.Fields)
		if err != nil {
			return "", err
		}
		var row []interface{}
		for _, field := range sql.Fields {
			row = append(row, m[field.Text()])
		}
		tbl.AppendRow(row)
	}

	if tbl.Length() == 0 {
		return util.EmptyResult(), nil
	}
	tbl.Render()

	return util.NRowsInSet(tbl.Length(), t), nil
}
