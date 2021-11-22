package syntax

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/pingcap/parser/ast"
	"os"
)

type Table struct {
	Name string
}

func (t *Table) TableName(node ast.Node) {
	if r, ok := node.(*ast.TableName); ok {
		t.Name = r.Name.String()
	}
	if r, ok := node.(*ast.SelectStmt); ok {
		t.Name = r.From.TableRefs.Left.(*ast.TableSource).Source.(*ast.TableName).Name.String()
	}
}

func (t *Table) CheckTableName() string {
	if t.Name != "tikv" {
		return fmt.Sprintf("Illegal table name: %s, must be tikv.", t.Name)
	}
	return ""
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

func NewKvTableRow(tbl table.Writer, kv KvPair, value []byte) {
	var row []interface{}
	row = append(row, kv.Key)
	row = append(row, string(value))
	tbl.AppendRow(row)
}
