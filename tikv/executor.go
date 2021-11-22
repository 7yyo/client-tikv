package tikv

import (
	"fmt"
	"os"
	"strings"
	"tikv-client/syntax"
)

func (c *Completer) Executor(s string) {

	if s == "exit" || s == "quit" {
		fmt.Println("Bye")
		os.Exit(0)
		return
	}

	if strings.TrimSpace(s) == "" {
		return
	}

	astNode, err := syntax.Parse(s)
	if err != nil {
		fmt.Println(err)
		return
	}

	sql := syntax.ParseSQL(astNode)
	if sql.Error != "" {
		fmt.Println(sql.Error)
		return
	}

	var r string
	switch sql.Operate {
	case "get":
		r, err = c.Get(sql)
	}
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(r)
	}

}
