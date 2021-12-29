package tikv

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"tikv-client/syntax"
)

func (c *Completer) Executor(s string) {

	if strings.TrimSpace(strings.ToUpper(s)) == "EXIT" {
		fmt.Println("GoodBye!")
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

	// For get some information of tikv
	switch strings.ToUpper(sql.Table.Name) {
	case "REGIONS":
		result, err = c.GetRegionInfo(sql)
		fmt.Println(result)
		return
	}

	// For get, put, delete...
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

func ExecuteAndGetResult(s string) (string, error) {

	out := &bytes.Buffer{}
	cmd := exec.Command("/bin/sh", "-c", s)
	cmd.Stdin = os.Stdin
	cmd.Stdout = out
	if err := cmd.Run(); err != nil {
		return "", err
	}
	r := string(out.Bytes())
	return r, nil

}
