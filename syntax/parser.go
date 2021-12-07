package syntax

import (
	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
)

func Parse(sql string) (*ast.StmtNode, error) {
	p := parser.New()
	stmtNodes, _, err := p.Parse(sql, "utf-8", "")
	if err != nil {
		return nil, err
	}
	return &stmtNodes[0], nil
}

func ParseSQL(s string) (*SQL, error) {
	astNode, err := Parse(s)
	if err != nil {
		return nil, err
	}
	return Parser(astNode), nil
}
