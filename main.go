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
	//c, err := tikv.NewCompleter(pdEndPoints)
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
		prompt.OptionSuggestionBGColor(prompt.Black),
		prompt.OptionSuggestionTextColor(prompt.White),
		prompt.OptionDescriptionBGColor(prompt.Red),
		prompt.OptionDescriptionTextColor(prompt.White),
		prompt.OptionTitle("tikv-client"),
	)

	p.Run()

}

func welcome() {
	fmt.Println()
	fmt.Println("Welcome to tikv-client. Commands exit to exit. \nServer version: alpha PingCAP Community Server - GPL")
	fmt.Println()
}

func bye() string {
	return "Bye"
}

func label() string {
	return "tikv > "
}
