package builder

import (
	"bytes"
	"github.com/pingcap/parser/ast"
	"os"
	"os/exec"
	e "tikv-client/error"
)

type Executor interface {
	Execute() error
	CheckValid() error
}

func (c *Completer) build(astNode ast.Node) (Executor, error) {

	switch node := astNode.(type) {
	case *ast.SelectStmt:
		if len(node.Fields.Fields) != 0 {
			if _, ok := node.Fields.Fields[0].Expr.(*ast.FuncCallExpr); ok {
				return c.buildFunc(node)
			}
		}
		return c.buildSelect(node)
	case *ast.InsertStmt:
		return c.buildInsert(node)
	default:
		return nil, e.ErrUnsupportedSQL(node.Text())
	}

}

func ExecuteAndGetResult(s string) (string, error) {

	out := &bytes.Buffer{}
	cmd := exec.Command("/bin/sh", "-c", s)
	cmd.Stdin = os.Stdin
	cmd.Stdout = out
	if err := cmd.Run(); err != nil {
		return "", err
	}
	r := string(out.Bytes())
	return r, nil

}

//func (s *SQL) checkSyntax(astNode ast.Node) {
//	switch node := astNode.(type) {
//	case *ast.SelectStmt:
//		tName := GetTableName(astNode)
//
//		if tName != "regions" {
//			// Normal select should be `=` or `in`
//			switch whereNode := node.Where.(type) {
//			case *ast.BinaryOperationExpr:
//				if whereNode.Op != opcode.EQ {
//					s.unsupportedSQL(astNode.Text())
//				}
//			case *ast.PatternInExpr:
//				if whereNode.Not {
//					s.unsupportedSQL(astNode.Text())
//				}
//			default:
//				s.unsupportedSQL(astNode.Text())
//			}
//		} else {
//			if node.Where != nil {
//				s.unsupportedSQL(astNode.Text())
//			}
//		}
//	case *ast.InsertStmt:
//		for _, list := range node.Lists {
//			if len(list) != 2 {
//				s.Error = "Number of illegal insertions, kv should only have two values."
//			}
//			return
//		}
//
//	}
//}
