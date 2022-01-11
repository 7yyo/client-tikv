package tikv

import (
	"encoding/json"
	tbl "github.com/jedib0t/go-pretty/v6/table"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/opcode"
	"github.com/tikv/client-go/v2/tikv"
	"strings"
	"tikv-client/pd"
	"tikv-client/util"
	"time"
)

type selectExecutor struct {
	sql             string
	placementDriver pd.PlacementDriverGroup
	table
	fields  []*ast.SelectField
	kvPairs []kvPair
	orderBy string
	where
	limit  int64
	client *tikv.RawKVClient
}

type where struct {
	op       opcode.Op
	moreThan string // <
	lessThan string // >
}

func (s selectExecutor) Valid() error {
	if s.table.name == "" {
		return errUnsupportedSQL(s.sql)
	}
	if err := checkTblName(s); err != nil {
		return err
	}
	if len(s.kvPairs) == 0 && s.limit == 0 && s.where.op == 0 {
		return errUnsupportedTiKVScan()
	}
	return nil
}

func (s selectExecutor) Exec() error {
	switch strings.ToLower(s.table.name) {
	case TikvTable_:
		if s.limit != 0 || s.op != 0 {
			return scan(s)
		}
		//if isOnlyLimit(s) {
		//	return limit(s)
		//}
		//if isScanLimit(s) {
		//	return scanLimit(s)
		//}
		if searchKvPairs(s.fields) {
			return keyValue(s)
		}
		return fields(s)
	}
	return nil
}

func keyValue(s selectExecutor) error {
	t := time.Now()
	tw := util.NewKvDisplayTable()
	var r []byte
	var rs [][]byte
	var err error
	if len(s.kvPairs) != 1 {
		rs, err = s.client.BatchGet(getKeys(s.kvPairs))
	} else {
		r, err = s.client.Get(s.kvPairs[0].key)
		rs = append(rs, r)
	}
	for i, kv := range s.kvPairs {
		kv.value = rs[i]
		if len(kv.value) != 0 {
			newKvTableRow(tw, kv)
		}
	}
	if err != nil {
		return err
	}
	if tw.Length() == 0 {
		util.EmptyResult(t)
		return nil
	}
	tw.SortBy([]tbl.SortBy{{Name: strings.ToUpper(s.orderBy), Mode: tbl.AscNumeric}})
	tw.Render()
	util.QueryOkNRows(tw.Length(), t)
	return nil
}

func fields(s selectExecutor) error {
	t := time.Now()
	var fieldNames []interface{}
	for _, field := range s.fields {
		fieldNames = append(fieldNames, field.Text())
	}
	tw := util.NewNormalDisplayTable(fieldNames)
	m := make(map[string]string)
	var r []byte
	var rs [][]byte
	var err error
	if len(s.kvPairs) == 1 {
		r, err = s.client.Get(s.kvPairs[0].key)
		rs = append(rs, r)
	} else {
		rs, err = s.client.BatchGet(getKeys(s.kvPairs))
	}
	if err != nil {
		return err
	}
	for i := range s.kvPairs {
		value := rs[i]
		if len(value) == 0 {
			continue
		}
		err = json.Unmarshal(rs[i], &m)
		if err != nil {
			return errParseJson(err.Error(), string(value))
		}
		err = hasField(m, s.fields)
		if err != nil {
			return err
		}
		var row []interface{}
		for _, field := range s.fields {
			row = append(row, m[field.Text()])
		}
		tw.AppendRow(row)
	}
	if tw.Length() == 0 {
		util.EmptyResult(t)
		return nil
	}
	tw.SortBy([]tbl.SortBy{{Name: s.orderBy, Mode: tbl.AscNumeric}})
	tw.Render()
	util.NRowsInSet(tw.Length(), t)
	return nil
}

func scan(s selectExecutor) error {
	t := time.Now()
	var keys [][]byte
	var values [][]byte
	var err error
	startKey := []byte(s.where.lessThan)
	endKey := []byte(s.where.moreThan)
	// TODO Should be `(startKey, endKey)`
	if s.limit == 0 {
		keys, values, err = s.client.Scan(startKey, endKey, tikv.MaxRawKVScanLimit)
	} else {
		keys, values, err = s.client.Scan(startKey, endKey, int(s.limit))
	}
	if err != nil {
		if err.Error() == "limit should be less than MaxRawKVScanLimit" {
			return errScanLimit()
		}
		return err
	}
	table := util.NewKvDisplayTable()
	var kv kvPair
	for i, _ := range keys {
		kv.key = keys[i]
		kv.value = values[i]
		if kv.value != nil {
			newKvTableRow(table, kv)
		}
	}
	if table.Length() == 0 {
		util.EmptyResult(t)
		return nil
	}
	table.Render()
	util.QueryOkNRows(table.Length(), t)
	return nil
}

//func limit(s selectExecutor) error {
//	t := time.Now()
//	var keys [][]byte
//	var values [][]byte
//	var err error
//	startKey := []byte(s.where.lessThan)
//	endKey := []byte(s.where.moreThan)
//	keys, values, err = s.client.Scan(startKey, endKey, int(s.limit))
//	if err != nil {
//		if err.Error() == "limit should be less than MaxRawKVScanLimit" {
//			return errScanLimit()
//		}
//		return err
//	}
//	table := util.NewKvDisplayTable()
//	var kv kvPair
//	for i, _ := range keys {
//		kv.key = keys[i]
//		kv.value = values[i]
//		if kv.value != nil {
//			newKvTableRow(table, kv)
//		}
//	}
//	if table.Length() == 0 {
//		util.EmptyResult(t)
//		return nil
//	}
//	table.Render()
//	util.QueryOkNRows(table.Length(), t)
//	return nil
//}
//
//func scanLimit(s selectExecutor) error {
//	t := time.Now()
//	var keys [][]byte
//	var values [][]byte
//	var err error
//
//	startKey := []byte(s.where.lessThan)
//	endKey := []byte(s.where.moreThan)
//	keys, values, err = s.client.Scan(startKey, endKey, int(s.limit))
//	if err != nil {
//		if err.Error() == "limit should be less than MaxRawKVScanLimit" {
//			return errScanLimit()
//		}
//		return err
//	}
//	table := util.NewKvDisplayTable()
//	var kv kvPair
//	for i, _ := range keys {
//		kv.key = keys[i]
//		kv.value = values[i]
//		if kv.value != nil {
//			newKvTableRow(table, kv)
//		}
//	}
//	if table.Length() == 0 {
//		util.EmptyResult(t)
//		return nil
//	}
//	table.Render()
//	util.QueryOkNRows(table.Length(), t)
//	return nil
//}

func hasField(m map[string]string, fields []*ast.SelectField) error {
	for _, field := range fields {
		if m[field.Text()] == "" {
			return errUnknownColumn(field.Text())
		}
	}
	return nil
}

// select * from t where k...
func isOnlyScan(s selectExecutor) bool {
	return s.limit == 0 && s.where.op != 0
}

// select * from t limit N
func isOnlyLimit(s selectExecutor) bool {
	return s.limit != 0 && s.where.op == 0
}

// select * from t where k... limit N
func isScanLimit(s selectExecutor) bool {
	return s.limit != 0 && s.where.op != 0
}
