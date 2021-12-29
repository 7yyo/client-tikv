package syntax

import (
	"github.com/pingcap/parser/ast"
	"strings"
	e "tikv-client/error"
)

type Table struct {
	Name string
}

func GetTableName(astNode ast.Node) string {
	switch node := astNode.(type) {
	case *ast.TableName:
		return node.Name.String()
	case *ast.SelectStmt:
		if node.From != nil {
			return node.From.TableRefs.Left.(*ast.TableSource).Source.(*ast.TableName).Name.String()
		}
		return ""
	case *ast.InsertStmt:
		return node.Table.TableRefs.Left.(*ast.TableSource).Source.(*ast.TableName).Name.String()
	default:
		return ""
	}
}

func (t *Table) tableName(astNode ast.Node) {
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
	tName := strings.ToUpper(t.Name)
	if tName != "TIKV" && tName != "REGIONS" {
		return e.ErrTableName(t.Name)
	}
	return nil
}
