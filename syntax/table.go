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

func (t *Table) TableName(astNode ast.Node) {

	switch node := astNode.(type) {
	case *ast.TableName:
		t.Name = node.Name.String()
	case *ast.SelectStmt:
		t.Name = node.From.TableRefs.Left.(*ast.TableSource).Source.(*ast.TableName).Name.String()
	case *ast.InsertStmt:
		t.Name = node.Table.TableRefs.Left.(*ast.TableSource).Source.(*ast.TableName).Name.String()
	}

}

func (t *Table) CheckTableName() string {
	if t.Name != "tikv" {
		s := fmt.Sprintf("Illegal table name: `%s`, must be tikv.", t.Name)
		fmt.Println(s)
		return s
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
