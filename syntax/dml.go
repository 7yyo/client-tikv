package syntax

import (
	"fmt"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/mysql"
	"github.com/pingcap/parser/test_driver"
)

type Sql struct {
	Type   string
	Fields []*ast.SelectField
	Table
}

const selectNum = 1
const insertNum = 2

type Table struct {
	Name    string
	Columns []Column
}

type Column struct {
	ColValue string
	ColType  byte
}

func Parser(node *ast.StmtNode) *Sql {
	s := getParseInfo(node)
	if !s.checkTable() || !s.checkColumns() || !s.checkFields() {
		return nil
	}
	return s
}

func getParseInfo(node *ast.StmtNode) *Sql {
	sql := Sql{}
	(*node).Accept(&sql)
	return &sql
}

func (s *Sql) Enter(node ast.Node) (ast.Node, bool) {

	s.setSqlType(node)
	s.setTableName(node)

	switch r := node.(type) {
	case *ast.InsertStmt:
		for _, v := range r.Lists[0] {
			c := Column{}
			if r, ok := v.(*test_driver.ValueExpr); ok {
				c.ColType = r.Type.Tp
				switch r.Type.Tp {
				case mysql.TypeVarString:
					c.ColValue = r.GetString()
				default:
				}
				s.Table.Columns = append(s.Table.Columns, c)
			}
		}
	case *ast.SelectStmt:
		c := Column{
			ColValue: r.Where.(*ast.BinaryOperationExpr).R.(*test_driver.ValueExpr).Datum.GetString(),
			ColType:  r.Where.GetType().Tp,
		}
		s.Table.Columns = append(s.Table.Columns, c)
		s.Fields = r.Fields.Fields
	case *ast.DeleteStmt:
		c := Column{
			ColValue: r.Where.(*ast.BinaryOperationExpr).R.(*test_driver.ValueExpr).Datum.GetString(),
			ColType:  r.Where.GetType().Tp,
		}
		s.Table.Columns = append(s.Table.Columns, c)
	case *ast.ExplainStmt:
		fmt.Println("This is a kv storage, no need to execute explain.")
		return node, true
	}
	return node, false
}

func (s *Sql) Leave(node ast.Node) (ast.Node, bool) {
	return node, true
}

func (s *Sql) checkTable() bool {
	switch s.Type {
	case "explain":
		return false
	default:
		return s.Table.checkTableName()
	}
}

func (t *Table) checkTableName() bool {
	if t.Name != "tikv" {
		fmt.Printf("Unknown table: %s, should be tikv\n", t.Name)
		return false
	}
	return true
}

func (s *Sql) checkColumns() bool {
	return s.checkColumnLength() || s.checkColumnType()
}

func (s *Sql) checkColumnLength() bool {
	switch s.Type {
	case "insert":
		if len(s.Table.Columns) != insertNum {
			fmt.Printf("Number of illegal fields: %d, should be %d.\n", len(s.Table.Columns), insertNum)
			return false
		}
	case "select":
	case "delete":
		if len(s.Table.Columns) != selectNum {
			fmt.Printf("Number of illegal fields: %d, should be %d.\n", len(s.Table.Columns), selectNum)
			return false
		}
	default:
	}
	return true
}

func (s *Sql) checkColumnType() bool {
	for _, c := range s.Table.Columns {
		if c.ColType != mysql.TypeVarString {
			fmt.Printf("Illegal column type: %d, should be %s\n", c.ColType, "var_string")
			return false
		}
	}
	return true
}

func (s *Sql) checkFields() bool {
	switch s.Type {
	case "select":
		if len(s.Fields) != 1 || s.Fields[0].Text() != "" {
			fmt.Println("Illegal field name, should be *")
			return false
		}
	default:
	}
	return true
}

func (s *Sql) setSqlType(node ast.Node) {
	if _, ok := node.(*ast.InsertStmt); ok {
		s.Type = "insert"
	} else if _, ok := node.(*ast.SelectStmt); ok {
		s.Type = "select"
	} else if _, ok := node.(*ast.DeleteStmt); ok {
		s.Type = "delete"
	} else if _, ok := node.(*ast.ExplainStmt); ok {
		s.Type = "explain"
	}
}

func (s *Sql) setTableName(node ast.Node) {
	if r, ok := node.(*ast.TableName); ok {
		s.Table.Name = r.Name.String()
	}
	if r, ok := node.(*ast.SelectStmt); ok {
		s.Table.Name = r.From.TableRefs.Left.(*ast.TableSource).Source.(*ast.TableName).Name.String()
	}
}
