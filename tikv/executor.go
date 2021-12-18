package tikv

import (
	"fmt"
	"os"
	"strings"
	"tikv-client/syntax"
)

func (c *Completer) Executor(s string) {

	// Execute cmd, clean up operateType
	c.operateType = ""

	if strings.TrimSpace(strings.ToUpper(s)) == syntax.Exit {
		fmt.Println("Bye")
		os.Exit(0)
		return
	}

	if strings.TrimSpace(s) == "" {
		return
	}

	sql, err := syntax.ParseSQL(s)
	if err != nil {
		fmt.Println(err)
		return
	}
	if sql.Error != "" {
		fmt.Println(sql.Error)
		return
	}

	if err := sql.Table.CheckTableName(); err != nil {
		fmt.Println(err.Error())
		return
	}

	var result string
	switch sql.Operate {
	case "get":
		result, err = c.Get(sql)
	case "put":
		result, err = c.Put(sql)
	}
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(result)
	}

}
