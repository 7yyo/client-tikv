package syntax

import (
	"fmt"
	"github.com/pingcap/parser/ast"
)

type SQL struct {
	Table
	Operate string
	Fields  []*ast.SelectField
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
		s.KvPairs = *NewKv(node)
		s.Fields = node.Fields.Fields
		s.Operate = "get"
	default:
		s.Error = fmt.Sprintf("Unsupport SQL: %s", astNode.Text())
		return astNode, true
	}
	return astNode, true
}

func (s *SQL) Leave(node ast.Node) (ast.Node, bool) {
	return node, true
}

func (s *SQL) CanOperate() string {
	if errMsg := s.Table.CheckTableName(); errMsg != "" {
		return errMsg
	}
	if errMsg := s.checkFields(); errMsg != "" {
		return errMsg
	}
	return ""
}

func (s *SQL) checkFields() string {
	if s.Fields[0].Text() != "kv" || len(s.Fields) != 1 {
		return "Illegal field name, must be kv"
	}
	return ""
}
