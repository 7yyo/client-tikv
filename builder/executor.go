package builder

import (
	"fmt"
	"github.com/pingcap/parser/ast"
	"os"
	"strings"
	"tikv-client/session"
)

func (c *Completer) Executor(sql string) {

	if strings.TrimSpace(strings.ToUpper(sql)) == "EXIT" {
		if err := c.Client.Close(); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("GoodBye!")
		os.Exit(0)
		return
	}

	if strings.TrimSpace(sql) == "" {
		return
	}

	var stmtNodes []ast.StmtNode
	var err error

	if stmtNodes, err = session.ParseSQL(sql); err != nil {
		fmt.Println(err)
		return
	}

	var e Executor

	if e, err = c.build(stmtNodes[0]); err != nil {
		fmt.Println(err)
		return
	}

	if err = e.Execute(); err != nil {
		fmt.Println(err)
		return
	}

}
