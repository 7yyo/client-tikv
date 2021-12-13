package tikv

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/pingcap/parser/ast"
	"github.com/tikv/client-go/v2/config"
	"github.com/tikv/client-go/v2/tikv"
	p "github.com/tikv/pd/client"
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
	c, err := tikv.NewRawKVClient(pdEndPoint, config.DefaultConfig().Security, p.WithMaxErrorRetry(1))
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Completer) Put(sql *syntax.SQL) (string, error) {
	t := time.Now()
	for _, kv := range sql.KvPairs {
		err := c.client.Put([]byte(kv.Key), []byte(kv.Value))
		if err != nil {
			return "", err
		}
	}
	return queryOkNRows(1, t), nil
}

func (c *Completer) Get(sql *syntax.SQL) (string, error) {
	if isSearchKv(sql.Fields) {
		return c.getKv2Pairs(sql)
	} else {
		return c.getKv2Field(sql)
	}
}

func (c *Completer) getKv2Pairs(sql *syntax.SQL) (string, error) {
	t := time.Now()
	tbl := syntax.NewKvDisplayTable()
	for _, kv := range sql.KvPairs {
		value, err := c.client.Get([]byte(kv.Key))
		if err != nil {
			return "", err
		}
		if value != nil {
			syntax.NewKvTableRow(tbl, kv, value)
		}
	}
	if tbl.Length() == 0 {
		return emptyResult(), nil
	}
	tbl.Render()
	return queryOkNRows(tbl.Length(), t), nil
}

func (c *Completer) getKv2Field(sql *syntax.SQL) (string, error) {

	t := time.Now()

	var fieldNames []interface{}
	for _, field := range sql.Fields {
		fieldNames = append(fieldNames, field.Text())
	}

	tbl := syntax.NewNormalDisplayTable(fieldNames)
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
		err = hasField(m, sql.Fields)
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
		return emptyResult(), nil
	}
	tbl.Render()
	return fmt.Sprintf("%d rows in set (%f sec)", tbl.Length(), util.Duration(t)), nil
}

func hasField(m map[string]string, fields []*ast.SelectField) error {
	for _, f := range fields {
		if m[f.Text()] == "" {
			return e.ErrUnknownColumn(f.Text())
		}
	}
	return nil
}

func isSearchKv(fields []*ast.SelectField) bool {
	if len(fields) == 1 && fields[0].Text() == "kv" {
		return true
	}
	return false
}

func emptyResult() string {
	return "Empty set (0.00 sec)"
}

func queryOk0Rows(t time.Time) string {
	return fmt.Sprintf("Query OK, 0 rows affected (%f sec)", util.Duration(t))
}

func queryOkNRows(n int, t time.Time) string {
	return fmt.Sprintf("Query OK, %d rows affected (%f sec)", n, util.Duration(t))
}
