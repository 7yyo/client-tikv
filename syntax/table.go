package syntax

import (
	"github.com/pingcap/parser/ast"
	e "tikv-client/error"
)

type Table struct {
	Name string
}

func (t *Table) setTableName(astNode ast.Node) {

	switch node := astNode.(type) {
	case *ast.TableName:
		t.Name = node.Name.String()
	case *ast.SelectStmt:
		t.Name = node.From.TableRefs.Left.(*ast.TableSource).Source.(*ast.TableName).Name.String()
	case *ast.InsertStmt:
		t.Name = node.Table.TableRefs.Left.(*ast.TableSource).Source.(*ast.TableName).Name.String()
	}

}

func (t *Table) CheckTableName() error {
	if t.Name != Tikv {
		return e.ErrTableName(t.Name)
	}
	return nil
}
