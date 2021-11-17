package main

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	_ "github.com/pingcap/parser/test_driver"
	"os"
	"tikv-client/tikv"
)

func main() {

	pdEndPoints := tikv.ParseArgs()
	c, err := tikv.NewCompleter(pdEndPoints)
	if err != nil {
		fmt.Printf("Error: %s, exit", err)
		os.Exit(0)
	}

	fmt.Println("Welcome to tikv client. Commands exit or quit to exit.\nServer version: 0.0.20 PingCAP Community Server - GPL")
	defer fmt.Println("Bye")

	p := prompt.New(
		c.Executor,
		c.Complete,
		prompt.OptionPrefix("tikv> "),
	)

	p.Run()

}
