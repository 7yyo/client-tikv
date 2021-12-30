package builder

import (
	"encoding/json"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/mysql"
	"github.com/pingcap/parser/test_driver"
	"github.com/tikv/client-go/v2/tikv"
	"reflect"
	"strconv"
	"strings"
	error2 "tikv-client/error"
	http2 "tikv-client/http"
	pd2 "tikv-client/pd"
	"time"
)

type SelectExecutor struct {
	pdEndPoints     string
	placementDriver pd2.Pd
	SQL             string
	TableName       string
	Fields          []*ast.SelectField
	KvPairs         []kvPair
	OrderBy         string
	client          *tikv.RawKVClient
}

var validTables = map[string]string{
	"tikv":    "y",
	"regions": "y",
}

func (c *Completer) buildSelect(node *ast.SelectStmt) (Executor, error) {

	var s SelectExecutor
	var kv kvPair

	s.SQL = node.Text()
	s.Fields = node.Fields.Fields
	s.placementDriver = c.Pd
	s.pdEndPoints = c.pdEndPoint[0]

	if node.OrderBy != nil {
		s.OrderBy = strings.ToUpper(node.OrderBy.Items[0].Expr.(*ast.ColumnNameExpr).Name.OrigColName())
	}

	if r, ok := node.Where.(*ast.BinaryOperationExpr); ok {
		if _, ok := r.R.(*ast.ColumnNameExpr); ok {
			return nil, error2.ErrUnsupportedSQL()
		}
		valueExpr := r.R.(*test_driver.ValueExpr)
		if valueExpr.Type.Tp == mysql.TypeVarString {
			kv.Key = valueExpr.Datum.GetString()
		}
		s.KvPairs = append(s.KvPairs, kv)
	} else if r, ok := node.Where.(*ast.PatternInExpr); ok {
		for _, pair := range r.List {
			valueExpr := pair.(*test_driver.ValueExpr)
			if valueExpr.Type.Tp == mysql.TypeVarString {
				kv.Key = valueExpr.Datum.GetString()
			}
			s.KvPairs = append(s.KvPairs, kv)
		}
	}

	s.TableName = node.From.TableRefs.Left.(*ast.TableSource).Source.(*ast.TableName).Name.String()

	if err := s.CheckValid(); err != nil {
		return nil, err
	}
	s.client = c.Client

	return &s, nil
}

func (s SelectExecutor) CheckValid() error {
	if s.TableName == "" {
		return error2.ErrUnsupportedSQL(s.SQL)
	}
	if validTables[s.TableName] != "y" {
		return error2.ErrTableName(s.TableName)
	}
	return nil
}

func (s SelectExecutor) Execute() error {
	switch strings.ToUpper(s.TableName) {
	case "TIKV":
		if searchKvPairs(s.Fields) {
			return kv(s)
		} else {
			return fields(s)
		}
	case "REGIONS":
		if err := s.GetRegionInfo(); err != nil {
			return err
		}
	}
	return nil
}

/**
>>>> select * from tikv where k = 'jim'
+-----+-----------------------------------------------+
| KEY | VALUE                                         |
+-----+-----------------------------------------------+
| jim | {"id":"1","name":"Jim","url":"www.baidu.com"} |
+-----+-----------------------------------------------+
*/
func kv(s SelectExecutor) error {

	t := time.Now()
	tbl := NewKvDisplayTable()

	rs, err := s.client.BatchGet(getKeys(s.KvPairs))
	if err != nil {
		return err
	}

	for i, kv := range s.KvPairs {
		kv.Value = string(rs[i])
		if kv.Value != "" {
			NewKvTableRow(tbl, kv)
		}
	}

	if tbl.Length() == 0 {
		fmt.Println(EmptyResult(t))
		return nil
	}

	tbl.Render()
	fmt.Println(QueryOkNRows(tbl.Length(), t))

	return nil
}

/**
>>>> select id, name, url from tikv where k = 'jim'
+----+------+---------------+
| ID | NAME | URL           |
+----+------+---------------+
| 1  | Jim  | www.baidu.com |
+----+------+---------------+
*/
func fields(s SelectExecutor) error {

	t := time.Now()

	var fieldNames []interface{}
	for _, field := range s.Fields {
		fieldNames = append(fieldNames, field.Text())
	}

	tbl := NewNormalDisplayTable(fieldNames)

	m := make(map[string]string)

	rs, err := s.client.BatchGet(getKeys(s.KvPairs))
	if err != nil {
		return err
	}

	for i, _ := range s.KvPairs {

		value := string(rs[i])
		if value == "" {
			continue
		}

		err = json.Unmarshal(rs[i], &m)
		if err != nil {
			return error2.ErrParseJson(err.Error(), value)
		}

		err = hasField(m, s.Fields)
		if err != nil {
			return err
		}

		var row []interface{}
		for _, field := range s.Fields {
			row = append(row, m[field.Text()])
		}
		tbl.AppendRow(row)

	}

	if tbl.Length() == 0 {
		fmt.Println(EmptyResult(t))
		return nil
	}

	tbl.Render()
	fmt.Println(NRowsInSet(tbl.Length(), t))

	return nil

}

func hasField(m map[string]string, fields []*ast.SelectField) error {
	for _, field := range fields {
		if m[field.Text()] == "" {
			return error2.ErrUnknownColumn(field.Text())
		}
	}
	return nil
}

func (s SelectExecutor) GetRegionInfo() error {

	t := time.Now()

	body, err := http2.ReqGet(fmt.Sprintf("http://%s/pd/api/v1/regions", s.pdEndPoints))
	if err != nil {
		return err
	}

	var regions Regions
	err = json.Unmarshal(body, &regions)
	if err != nil {
		return err
	}

	var rts []RegionTable
	var rt RegionTable

	for _, region := range regions.Regions {
		rt.REGION_ID = region.Id
		rt.START_KEY = region.Start_key
		rt.END_KEY = region.End_key
		rt.LEADER_ID = region.Leader.Id
		rt.LEADER_STORE_ID = region.Leader.Store_id
		if len(region.Peers) != 0 {
			rt.PEERS = strconv.Itoa(region.Peers[0].Id)
		}
		i := 0
		for _, p := range region.Peers {
			if i == 0 {
				i++
				continue
			}
			rt.PEERS += "," + strconv.Itoa(p.Id)
		}
		rt.WRITTEN_BYTES = region.Written_bytes
		rt.READ_BYTES = region.Read_bytes
		rt.APPROXIMATE_SIZE = region.Approximate_size
		rt.APPROXIMATE_KEYS = region.Approximate_keys
		rts = append(rts, rt)
	}

	var tbl table.Writer

	if len(rts) > 0 {

		var fn []interface{}
		tt := reflect.TypeOf(rts[0])

		for i := 0; i < tt.NumField(); i++ {
			fn = append(fn, tt.Field(i).Name)
		}

		tbl = NewNormalDisplayTable(fn)

	}

	for _, r := range rts {
		var row []interface{}
		row = append(row, r.REGION_ID)
		row = append(row, r.START_KEY)
		row = append(row, r.END_KEY)
		row = append(row, r.LEADER_ID)
		row = append(row, r.LEADER_STORE_ID)
		row = append(row, r.PEERS)
		row = append(row, r.WRITTEN_BYTES)
		row = append(row, r.READ_BYTES)
		row = append(row, r.APPROXIMATE_SIZE)
		row = append(row, r.APPROXIMATE_KEYS)
		tbl.AppendRow(row)
	}

	tbl.SortBy([]table.SortBy{{Name: s.OrderBy, Mode: table.AscNumeric}})

	tbl.Render()

	fmt.Println(NRowsInSet(regions.Count, t))
	return nil

}
