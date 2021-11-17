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

	sql := syntax.ParseSQL(*astNode)
	if sql == nil {
		return
	}

	var msg string
	kv := NewKVPair(sql)
	switch sql.Type {
	case "insert":
		msg, err = c.Put(kv)
	case "select":
		msg, err = c.Get(kv)
	case "delete":
		msg, err = c.Delete(kv)
	default:
	}

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(msg)

}
