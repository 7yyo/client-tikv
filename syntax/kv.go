package syntax

import (
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/mysql"
	"github.com/pingcap/parser/test_driver"
)

type KvPair struct {
	Key   string
	Value string
}

func ParseKvPairs(astNode ast.Node) []KvPair {
	var kvPairs []KvPair
	var kv KvPair
	switch node := astNode.(type) {
	case *ast.SelectStmt:
		if r, ok := node.Where.(*ast.BinaryOperationExpr); ok {
			valueExpr := r.R.(*test_driver.ValueExpr)
			if valueExpr.Type.Tp == mysql.TypeVarString {
				kv.Key = valueExpr.Datum.GetString()
			}
			kvPairs = append(kvPairs, kv)
		} else if r, ok := node.Where.(*ast.PatternInExpr); ok {
			for _, pair := range r.List {
				valueExpr := pair.(*test_driver.ValueExpr)
				if valueExpr.Type.Tp == mysql.TypeVarString {
					kv.Key = valueExpr.Datum.GetString()
				}
				kvPairs = append(kvPairs, kv)
			}
		}
	case *ast.InsertStmt:
		if len(node.Lists) != 1 || len(node.Lists[0]) == 0 {
			return nil
		}
		for _, list := range node.Lists {
			kvPairs = append(kvPairs, KvPair{
				Key:   list[0].(*test_driver.ValueExpr).GetDatumString(),
				Value: list[1].(*test_driver.ValueExpr).GetDatumString(),
			})
		}
	default:
	}
	return kvPairs
}
