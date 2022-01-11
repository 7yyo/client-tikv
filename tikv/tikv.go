package tikv

import (
	tbl "github.com/jedib0t/go-pretty/v6/table"
	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	"github.com/tikv/client-go/v2/tikv"
	"strings"
	p "tikv-client/pd"
)

const (
	Tikv_ = "tikv"
)

type TiKV struct {
	GoClient             *tikv.RawKVClient
	PlacementDriverGroup *p.PlacementDriverGroup
}

const (
	Select_     = "select"
	InsertInto_ = "insert into"
	Delete_     = "delete"
	From_       = "from"
	Where_      = "where"
	Values_     = "values"
	OrderBy_    = "order by"
	Into_       = "into"
	Limit_      = "limit"
	Truncate_   = "truncate"
)

type kvPair struct {
	key   []byte
	value []byte
}

func (tikv TiKV) Run(cmd string) error {
	stmtNode, err := ParseSQL(cmd)
	if err != nil {
		return err
	}
	b := ExecutorBuilder{
		GoClient:             tikv.GoClient,
		PlacementDriverGroup: tikv.PlacementDriverGroup,
		StmtNode:             stmtNode,
	}
	e, err := b.Build()
	if err != nil {
		return err
	}
	if err := e.Valid(); err != nil {
		return err
	}
	if err = e.Exec(); err != nil {
		return err
	}
	return nil
}

func ParseSQL(s string) ([]ast.StmtNode, error) {
	psr := parser.New()
	var err error
	stmtNode, _, err := psr.Parse(s, "utf-8", "")
	if err != nil {
		return nil, err
	}
	return stmtNode, nil
}

// GetKeys Get keys from kvpairs.
func getKeys(kvPairs []kvPair) [][]byte {
	var keys [][]byte
	for _, kv := range kvPairs {
		keys = append(keys, kv.key)
	}
	return keys
}

func searchKvPairs(fields []*ast.SelectField) bool {
	if len(fields) == 1 && strings.ToUpper(fields[0].Text()) == "" {
		return true
	}
	return false
}

// newKvTableRow Appand kv table rows.
func newKvTableRow(t tbl.Writer, kv kvPair) {
	var row []interface{}
	row = append(row, kv.kToString())
	row = append(row, kv.vToString())
	t.AppendRow(row)
}

func (kv *kvPair) kToString() string {
	return string(kv.key)
}

func (kv *kvPair) vToString() string {
	return string(kv.value)
}
