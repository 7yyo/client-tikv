package builder

import (
	"fmt"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/test_driver"
	"github.com/tikv/client-go/v2/tikv"
	e "tikv-client/error"
	"time"
)

type InsertExecutor struct {
	TableName string
	kvPairs   []kvPair
	client    *tikv.RawKVClient
}

func (c *Completer) buildInsert(node *ast.InsertStmt) (Executor, error) {

	var i InsertExecutor

	for _, list := range node.Lists {
		if len(list) != 2 {
			return nil, e.ErrColumnDontMatch()
		}
		i.kvPairs = append(i.kvPairs, kvPair{
			Key:   list[0].(*test_driver.ValueExpr).GetDatumString(),
			Value: list[1].(*test_driver.ValueExpr).GetDatumString(),
		})
	}

	i.TableName = node.Table.TableRefs.Left.(*ast.TableSource).Source.(*ast.TableName).Name.String()
	i.client = c.Client

	return i, nil
}

func (i InsertExecutor) Execute() error {
	t := time.Now()
	if i.kvPairs == nil {
		return e.ErrInsertValueCount()
	}
	for _, kv := range i.kvPairs {
		if err := i.client.Put([]byte(kv.Key), []byte(kv.Value)); err != nil {
			return err
		}
	}
	fmt.Println(QueryOkNRows(len(i.kvPairs), t))
	return nil
}

func (i InsertExecutor) CheckValid() error {
	return nil
}
