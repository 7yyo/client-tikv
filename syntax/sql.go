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

const (
	Get        = "select"
	Put        = "insert"
	Kv         = "kv"
	From       = "from"
	Tikv       = "tikv"
	Where      = "where"
	Key        = "k"
	Eq         = "="
	In         = "in"
	Apostrophe = "'';"
	BracketsIn = "('','');"
	Exit       = "exit"
	Quit       = "quit"
	Q          = "\\q"
	Into       = "into"
	Values     = "values"
)

var End = []int{5, 8}

func Parser(astNode *ast.StmtNode) *SQL {
	s := SQL{}
	(*astNode).Accept(&s)
	return &s
}

func (s *SQL) Enter(astNode ast.Node) (ast.Node, bool) {

	switch node := astNode.(type) {
	case *ast.SelectStmt:
		if node.From == nil {
			s.unsupportedSQL(astNode.Text())
			return astNode, true
		}
		s.Operate = "get"
		s.KvPairs = ParseKvPairs(node)
		s.Fields = node.Fields.Fields
	case *ast.InsertStmt:
		s.Operate = "put"
		s.KvPairs = ParseKvPairs(node)
		if s.KvPairs == nil {
			s.Error = "Currently only supports put single kv pair."
		}
	default:
		s.unsupportedSQL(astNode.Text())
	}
	s.Table.TableName(astNode)
	return astNode, true
}

func (s *SQL) Leave(node ast.Node) (ast.Node, bool) {
	return node, true
}

func (s *SQL) unsupportedSQL(sql string) {
	s.Error = fmt.Sprintf("Unsupported SQL: '%s'", sql)
}
