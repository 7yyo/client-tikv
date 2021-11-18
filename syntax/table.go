package syntax

import (
	"fmt"
	"github.com/pingcap/parser/ast"
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
