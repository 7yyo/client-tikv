package builder

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/pingcap/parser/ast"
	"os"
	"strings"
	"tikv-client/util"
	"time"
)

type kvPair struct {
	Key   string
	Value string
}

func getKeys(kvPairs []kvPair) [][]byte {
	var keys [][]byte
	for _, kv := range kvPairs {
		keys = append(keys, []byte(kv.Key))
	}
	return keys
}

func searchKvPairs(fields []*ast.SelectField) bool {
	if len(fields) == 1 && strings.ToUpper(fields[0].Text()) == "" {
		return true
	}
	return false
}

func NewKvDisplayTable() table.Writer {
	tbl := table.NewWriter()
	tbl.SetOutputMirror(os.Stdout)
	tbl.AppendHeader(table.Row{"Key", "Value"})
	return tbl
}

func NewNormalDisplayTable(titles []interface{}) table.Writer {
	tbl := table.NewWriter()
	tbl.SetOutputMirror(os.Stdout)
	t := table.Row{}
	for _, v := range titles {
		t = append(t, v)
	}
	tbl.AppendHeader(t)
	return tbl
}

func NewKvTableRow(tbl table.Writer, kv kvPair) {
	var row []interface{}
	row = append(row, kv.Key)
	row = append(row, kv.Value)
	tbl.AppendRow(row)
}

func EmptyResult(t time.Time) string {
	return fmt.Sprintf("Empty set (%f sec)", util.Duration(t))
}

func QueryOk0Rows(t time.Time) string {
	return fmt.Sprintf("Query OK, 0 rows affected (%f sec)", util.Duration(t))
}

func QueryOkNRows(n int, t time.Time) string {
	return fmt.Sprintf("Query OK, %d rows affected (%f sec)", n, util.Duration(t))
}

func NRowsInSet(n int, t time.Time) string {
	return fmt.Sprintf("%d rows in set (%f sec)", n, util.Duration(t))
}
