package syntax

import (
	"fmt"
	"github.com/pingcap/parser/ast"
)

type SQL struct {
	Table
	Operate string
	Fields  []*ast.SelectField
	Cols    []*ast.ColumnDef
	KvPairs []KvPair
	Error   string
}

func Parser(astNode *ast.StmtNode) *SQL {
	s := SQL{}
	(*astNode).Accept(&s)
	return &s
}

func (s *SQL) Enter(astNode ast.Node) (ast.Node, bool) {

	s.TableName(astNode)

	switch node := astNode.(type) {
	case *ast.SelectStmt:
		s.Operate = "get"
		s.KvPairs = ParseKvPairs(node)
		s.Fields = node.Fields.Fields
	default:
		s.Error = fmt.Sprintf("Unsupported SQL: '%s'", astNode.Text())
		return astNode, true
	}
	return astNode, true
}

func (s *SQL) Leave(node ast.Node) (ast.Node, bool) {
	return node, true
}
