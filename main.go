package main

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	_ "github.com/pingcap/parser/test_driver"
	"os"
	"tikv-client/tikv"
)

func main() {

	//pdEndPoints := tikv.ParseArgs()
	c, err := tikv.NewCompleter([]string{"172.16.5.133:2379"})
	if err != nil {
		fmt.Printf("%s, exit", err)
		os.Exit(0)
	}

	welcome()
	defer bye()

	p := prompt.New(
		c.Executor,
		c.Complete,
		prompt.OptionPrefix(label()),
	)

	p.Run()

}

func welcome() {
	fmt.Println("Welcome to tikv client. Commands exit or quit to exit. \nServer version: 0.0.20 PingCAP Community Server - GPL")
}

func bye() string {
	return "Bye"
}

func label() string {
	return "tikv > "
}
