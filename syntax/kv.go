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

func NewKv(astNode ast.Node) *[]KvPair {
	var kvPairs []KvPair
	switch node := astNode.(type) {
	case *ast.SelectStmt:
		var kv KvPair
		if r, ok := node.Where.(*ast.BinaryOperationExpr); ok {
			valueExpr := r.R.(*test_driver.ValueExpr)
			if valueExpr.Type.Tp == mysql.TypeVarString {
				kv.Key = valueExpr.Datum.GetString()
			}
			kvPairs = append(kvPairs, kv)
		}
		if r, ok := node.Where.(*ast.PatternInExpr); ok {
			var kv KvPair
			for _, pair := range r.List {
				valueExpr := pair.(*test_driver.ValueExpr)
				if valueExpr.Type.Tp == mysql.TypeVarString {
					kv.Key = valueExpr.Datum.GetString()
				}
				kvPairs = append(kvPairs, kv)
			}
		}
	default:
	}
	return &kvPairs
}
