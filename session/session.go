package session

import (
	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
)

func ParseSQL(s string) ([]ast.StmtNode, error) {
	p := parser.New()
	stmtNodes, _, err := p.Parse(s, "utf-8", "")
	if err != nil {
		return nil, err
	}
	return stmtNodes, nil
}
