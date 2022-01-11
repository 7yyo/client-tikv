package tikv

import (
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/mysql"
	"github.com/pingcap/parser/opcode"
	"github.com/pingcap/parser/test_driver"
	"github.com/tikv/client-go/v2/tikv"
	"strings"
	"tikv-client/pd"
)

type ExecutorBuilder struct {
	GoClient             *tikv.RawKVClient
	PlacementDriverGroup *pd.PlacementDriverGroup
	StmtNode             []ast.StmtNode
}

type Executor interface {
	Exec() error
	Valid() error
}

func (b *ExecutorBuilder) Build() (Executor, error) {
	if err := isEqOrIn(b.StmtNode[0]); err != nil {
		return nil, err
	}
	switch node := b.StmtNode[0].(type) {
	case *ast.InsertStmt:
		return b.buildInsert(node)
	case *ast.SelectStmt:
		return b.buildSelect(node)
	case *ast.DeleteStmt:
		return b.buildDelete(node)
	case *ast.TruncateTableStmt:
		return b.buildTruncate(node)
	default:
		return nil, errUnsupportedSQL(node.Text())
	}
}

func (b *ExecutorBuilder) buildInsert(n *ast.InsertStmt) (Executor, error) {
	var kvPairs []kvPair
	if n.IsReplace || n.IgnoreErr || len(n.OnDuplicate) != 0 {
		return nil, errUnsupportedSQL(n.Text())
	}
	for _, list := range n.Lists {
		if len(list) != 2 {
			return nil, errColumnDontMatch()
		}
		kvPairs = append(kvPairs, kvPair{
			key:   []byte(list[0].(*test_driver.ValueExpr).GetDatumString()),
			value: []byte(list[1].(*test_driver.ValueExpr).GetDatumString()),
		})
	}
	t := table{
		name: n.Table.TableRefs.Left.(*ast.TableSource).Source.(*ast.TableName).Name.String(),
	}
	return insertExecutor{
		table:   t,
		client:  b.GoClient,
		kvPairs: kvPairs,
	}, nil
}

func (b *ExecutorBuilder) buildSelect(n *ast.SelectStmt) (Executor, error) {
	var err error
	var w where
	if r, ok := n.Where.(*ast.BinaryOperationExpr); ok {
		switch r.Op {
		case opcode.LogicAnd:
			w.op = opcode.LogicAnd
			switch r.L.(*ast.BinaryOperationExpr).Op {
			case opcode.LT:
				w.moreThan = r.L.(*ast.BinaryOperationExpr).R.(*test_driver.ValueExpr).Datum.GetString()
			case opcode.GT:
				w.lessThan = r.L.(*ast.BinaryOperationExpr).R.(*test_driver.ValueExpr).Datum.GetString()
			}
			switch r.R.(*ast.BinaryOperationExpr).Op {
			case opcode.LT:
				w.moreThan = r.R.(*ast.BinaryOperationExpr).R.(*test_driver.ValueExpr).Datum.GetString()
			case opcode.GT:
				w.lessThan = r.R.(*ast.BinaryOperationExpr).R.(*test_driver.ValueExpr).Datum.GetString()
			}
			// > and > , < and <
			if (w.lessThan != "" && w.moreThan == "") || (w.lessThan == "" && w.moreThan != "") {
				return nil, errQueryScan()
			}
		case opcode.GT:
			w.op = opcode.GT
			w.lessThan = r.R.(*test_driver.ValueExpr).Datum.GetString()
		case opcode.LT:
			w.op = opcode.LT
			w.moreThan = r.R.(*test_driver.ValueExpr).Datum.GetString()
		case opcode.EQ:
			break
		default:
			return nil, errUnsupportedSQL(n.Text())
		}
	}
	var kvPairs []kvPair
	if w.op == 0 {
		// Collect kv pairs in where
		kvPairs, err = collectWhereKvPairs(n.Where)
		if err != nil {
			return nil, err
		}
	}
	s := selectExecutor{
		client:          b.GoClient,
		placementDriver: *b.PlacementDriverGroup,
		sql:             n.Text(),
		fields:          n.Fields.Fields,
		kvPairs:         kvPairs,
		where:           w,
	}
	if n.From != nil {
		s.table = table{
			name: n.From.TableRefs.Left.(*ast.TableSource).Source.(*ast.TableName).Name.String(),
		}
	}
	if n.OrderBy != nil {
		if len(n.OrderBy.Items) > 1 {
			return nil, errOrderbyCount()
		}
		s.orderBy = strings.ToLower(n.OrderBy.Items[0].Expr.(*ast.ColumnNameExpr).Name.OrigColName())
	}
	if n.Limit != nil {
		s.limit = n.Limit.Count.(*test_driver.ValueExpr).Datum.GetInt64()
	}
	return s, nil
}

func (b *ExecutorBuilder) buildDelete(n *ast.DeleteStmt) (Executor, error) {
	tbl := table{
		name: n.TableRefs.TableRefs.Left.(*ast.TableSource).Source.(*ast.TableName).Name.String(),
	}
	kvPairs, err := collectWhereKvPairs(n.Where)
	if err != nil {
		return nil, err
	}
	d := deleteExecutor{
		client:          b.GoClient,
		placementDriver: *b.PlacementDriverGroup,
		sql:             n.Text(),
		table:           tbl,
		kvPairs:         kvPairs,
	}
	return d, nil
}

func (b *ExecutorBuilder) buildTruncate(n *ast.TruncateTableStmt) (Executor, error) {
	t := table{
		name: n.Table.Name.String(),
	}
	return truncateExecutor{
		table:  t,
		client: b.GoClient,
	}, nil
}

func collectWhereKvPairs(e ast.ExprNode) ([]kvPair, error) {
	var kvPairs []kvPair
	var kv kvPair
	if r, ok := e.(*ast.BinaryOperationExpr); ok {
		if _, ok := r.R.(*ast.ColumnNameExpr); ok {
			return nil, errUnsupportedSQL()
		}
		valueExpr := r.R.(*test_driver.ValueExpr)
		if valueExpr.Type.Tp == mysql.TypeVarString {
			kv.key = []byte(valueExpr.Datum.GetString())
		}
		kvPairs = append(kvPairs, kv)
	} else if r, ok := e.(*ast.PatternInExpr); ok {
		for _, pair := range r.List {
			valueExpr := pair.(*test_driver.ValueExpr)
			if valueExpr.Type.Tp == mysql.TypeVarString {
				kv.key = []byte(valueExpr.Datum.GetString())
			}
			kvPairs = append(kvPairs, kv)
		}
	}
	return kvPairs, nil
}

func isEqOrIn(astNode ast.StmtNode) error {
	result := true
	switch node := astNode.(type) {
	case *ast.SelectStmt:
		if node.Limit != nil {
			return nil
		}
		if r, ok := node.Where.(*ast.PatternInExpr); ok {
			result = isIn(r)
		}
	case *ast.DeleteStmt:
		if r, ok := node.Where.(*ast.BinaryOperationExpr); ok {
			result = isEq(r)
		}
		if r, ok := node.Where.(*ast.PatternInExpr); ok {
			result = isIn(r)
		}
	case *ast.InsertStmt, *ast.TruncateTableStmt:
		return nil
	default:
	}
	if !result {
		return errUnsupportedSQL(astNode.Text())
	}
	return nil
}

func isEq(node *ast.BinaryOperationExpr) bool {
	return !(node.Op != opcode.EQ)
}

func isIn(node *ast.PatternInExpr) bool {
	return !node.Not
}
