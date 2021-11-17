package tikv

import "tikv-client/syntax"

type KvPair struct {
	Key   string
	Value string
}

func NewKVPair(sql *syntax.Sql) *KvPair {
	var kv KvPair
	switch sql.Type {
	case "insert":
		kv.Key = sql.Table.Columns[0].ColValue
		kv.Value = sql.Table.Columns[1].ColValue
	case "select":
		fallthrough
	case "delete":
		kv.Key = sql.Table.Columns[0].ColValue
	default:
	}
	return &kv
}
