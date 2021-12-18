package util

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
	"tikv-client/syntax"
	"time"
)

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

func NewKvTableRow(tbl table.Writer, kv syntax.KvPair) {
	var row []interface{}
	row = append(row, kv.Key)
	row = append(row, kv.Value)
	tbl.AppendRow(row)
}

func EmptyResult() string {
	return "Empty set (0.00 sec)"
}

func QueryOk0Rows(t time.Time) string {
	return fmt.Sprintf("Query OK, 0 rows affected (%f sec)", Duration(t))
}

func QueryOkNRows(n int, t time.Time) string {
	return fmt.Sprintf("Query OK, %d rows affected (%f sec)", n, Duration(t))
}
