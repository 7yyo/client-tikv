package syntax

import (
	"fmt"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/mysql"
	"github.com/pingcap/parser/opcode"
	"github.com/pingcap/parser/test_driver"
	"strings"
	e "tikv-client/error"
)

type SQL struct {
	Table
	Operate string
	Fields  []*ast.SelectField
	Cols    []*ast.ColumnDef
	OrderBy string
	KvPairs []KvPair
	Error   string
}

func Parser(astNode *ast.StmtNode) *SQL {
	s := SQL{}
	(*astNode).Accept(&s)
	return &s
}

func (s *SQL) Enter(astNode ast.Node) (ast.Node, bool) {

	s.checkSyntax(astNode)
	// Check syntax, if error, return
	if s.Error != "" {
		return astNode, true
	}

	s.operateType(astNode)
	s.fields(astNode)
	s.orderBy(astNode)
	s.kvPairs(astNode)
	s.Table.tableName(astNode)

	return astNode, true
}

func (s *SQL) Leave(node ast.Node) (ast.Node, bool) {
	return node, true
}

func (s *SQL) operateType(astNode ast.Node) {
	switch astNode.(type) {
	case *ast.SelectStmt:
		s.Operate = "get"
	case *ast.InsertStmt:
		s.Operate = "put"
	default:
		s.unsupportedSQL(astNode.Text())
	}
}

func (s *SQL) fields(astNode ast.Node) {
	switch node := astNode.(type) {
	case *ast.SelectStmt:
		s.Fields = node.Fields.Fields
	}
}

func (s *SQL) orderBy(astNode ast.Node) {
	switch node := astNode.(type) {
	case *ast.SelectStmt:
		if node.OrderBy != nil {
			s.OrderBy = strings.ToUpper(node.OrderBy.Items[0].Expr.(*ast.ColumnNameExpr).Name.OrigColName())
		}
	}
}

func (s *SQL) kvPairs(astNode ast.Node) {
	var kvPairs []KvPair
	var kv KvPair
	switch node := astNode.(type) {
	case *ast.SelectStmt:
		if r, ok := node.Where.(*ast.BinaryOperationExpr); ok {
			if _, ok := r.R.(*ast.ColumnNameExpr); ok {
				s.unsupportedSQL(astNode.Text())
				return
			}
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
		for _, list := range node.Lists {
			kvPairs = append(kvPairs, KvPair{
				Key:   list[0].(*test_driver.ValueExpr).GetDatumString(),
				Value: list[1].(*test_driver.ValueExpr).GetDatumString(),
			})
		}
	default:
	}
	s.KvPairs = kvPairs
}

func (s *SQL) checkSyntax(astNode ast.Node) {
	switch node := astNode.(type) {
	case *ast.SelectStmt:
		tName := GetTableName(astNode)
		if node.OrderBy != nil {
			if len(node.OrderBy.Items) > 1 {
				s.Error = "Only support 1 order by."
			}
		}
		if tName != "regions" {
			// Normal select should be `=` or `in`
			switch whereNode := node.Where.(type) {
			case *ast.BinaryOperationExpr:
				if whereNode.Op != opcode.EQ {
					s.unsupportedSQL(astNode.Text())
				}
			case *ast.PatternInExpr:
				if whereNode.Not {
					s.unsupportedSQL(astNode.Text())
				}
			default:
				s.unsupportedSQL(astNode.Text())
			}
		} else {
			if node.Where != nil {
				s.unsupportedSQL(astNode.Text())
			}
		}
	case *ast.InsertStmt:
		for _, list := range node.Lists {
			if len(list) != 2 {
				s.Error = "Number of illegal insertions, kv should only have two values."
			}
			return
		}

	}
}

func HasField(m map[string]string, fields []*ast.SelectField) error {
	for _, f := range fields {
		if m[f.Text()] == "" {
			return e.ErrUnknownColumn(f.Text())
		}
	}
	return nil
}

func IsSearchKvPairs(fields []*ast.SelectField) bool {
	if len(fields) == 1 && strings.ToUpper(fields[0].Text()) == "" {
		return true
	}
	return false
}

func (s *SQL) unsupportedSQL(sql string) {
	s.Error = fmt.Sprintf("Unsupported SQL: '%s'", sql)
}
