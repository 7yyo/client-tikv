package builder

import (
	"fmt"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/test_driver"
	"strings"
	error2 "tikv-client/error"
	pd2 "tikv-client/pd"
	"time"
)

const (
	TidbParseTso = "tidb_parse_tso"
)

const (
	Tso = "tiup ctl:%s pd -u http://%s tso %d"
)

type FuncExecutor struct {
	funcName        string
	pdEndPoints     string
	placementDriver pd2.Pd
	arg01           interface{}
	arg02           interface{}
	arg03           interface{}
	arg04           interface{}
	arg05           interface{}
	arg06           interface{}
	arg07           interface{}
	arg08           interface{}
	arg09           interface{}
	arg10           interface{}
}

func (c *Completer) buildFunc(node *ast.SelectStmt) (Executor, error) {
	var err error
	var e Executor
	switch node.Fields.Fields[0].Expr.(*ast.FuncCallExpr).FnName.String() {
	case TidbParseTso:
		e, err = c.buildTidbParseTso(node)
	default:
	}
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (c *Completer) buildTidbParseTso(node *ast.SelectStmt) (Executor, error) {
	var f FuncExecutor
	if len(node.Fields.Fields[0].Expr.(*ast.FuncCallExpr).Args) == 0 {
		return nil, error2.ErrColumnDontMatch()
	}
	f.arg01 = node.Fields.Fields[0].Expr.(*ast.FuncCallExpr).Args[0].(*test_driver.ValueExpr).Datum.GetValue()
	f.funcName = node.Fields.Fields[0].Expr.(*ast.FuncCallExpr).FnName.String()
	f.pdEndPoints = c.pdEndPoint[0]
	f.placementDriver = c.Pd
	return f, nil
}

func (f FuncExecutor) CheckValid() error {
	//TODO implement me
	panic("implement me")
}

func (f FuncExecutor) Execute() error {

	t := time.Now()

	var err error
	var r string
	var colNames []interface{}
	var row []interface{}
	switch f.funcName {
	case TidbParseTso:
		if r, err = ExecuteAndGetResult(fmt.Sprintf(Tso, f.placementDriver.Version, f.pdEndPoints, f.arg01)); err != nil {
			return err
		}
		colNames = append(colNames, "system")
		colNames = append(colNames, "logic")
		tbl := NewNormalDisplayTable(colNames)
		system := strings.TrimSpace(strings.ReplaceAll(strings.Split(strings.Split(r, "\n")[0], "system")[1], ":", ""))
		logic := strings.TrimSpace(strings.Split(strings.Split(r, "\n")[1], ":")[1])
		row := append(row, system)
		row = append(row, logic)

		tbl.AppendRow(row)

		tbl.Render()
		fmt.Println(NRowsInSet(1, t))

	}
	return nil
}
