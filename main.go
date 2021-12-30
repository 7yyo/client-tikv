package main

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	_ "github.com/pingcap/parser/test_driver"
	"os"
	"tikv-client/builder"
	p "tikv-client/pd"
)

func main() {

	pdEndPoints := builder.ParseArgs()
	//pdEndPoints := []string{"172.16.101.3:2479"}
	c, err := builder.NewCompleter(pdEndPoints)
	if err != nil {
		fmt.Printf("%s, exit", err)
		os.Exit(0)
	}

	welcome(c.Pd)

	defer bye()

	prpt := prompt.New(
		c.Executor,
		c.Complete,
		prompt.OptionPrefix(">>>> "),
		prompt.OptionSuggestionBGColor(prompt.Red),
		prompt.OptionSuggestionTextColor(prompt.White),
		prompt.OptionTitle("tikv-console"),
	)
	prpt.Run()

}

func welcome(p p.Pd) {

	fmt.Printf("\n"+
		"Welcome to tikv-console. Commands exit to exit. \n"+
		"Server version: alpha PingCAP Community Server - GPL\n"+
		"\n"+
		"PD version: %s\n"+
		"Build ts: %s\n"+
		"Git hash: %s\n"+
		"\n",
		p.Version,
		p.Build_ts,
		p.Git_hash)

}

func bye() string {
	return "GoodBye!"
}
