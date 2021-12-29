package main

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	_ "github.com/pingcap/parser/test_driver"
	"os"
	pd2 "tikv-client/pd"
	"tikv-client/tikv"
)

func main() {

	pdEndPoints := tikv.ParseArgs()
	//pdEndPoints := []string{"172.16.101.3:2379"}
	welcome(pdEndPoints[0])

	//c, err := tikv.NewCompleter(pdEndPoints)

	c, err := tikv.NewCompleter(pdEndPoints)
	if err != nil {
		fmt.Printf("%s, exit", err)
		os.Exit(0)
	}

	defer bye()

	p := prompt.New(
		c.Executor,
		c.Complete,
		prompt.OptionPrefix(">>>> "),
		prompt.OptionDescriptionBGColor(prompt.Red),
		prompt.OptionDescriptionTextColor(prompt.White),
		prompt.OptionTitle("tikv-client"),
	)
	p.Run()

}

func welcome(pdEndPoint string) {

	pd := pd2.PdInfo(pdEndPoint)

	fmt.Println()
	fmt.Printf("Welcome to tikv-client. Commands exit to exit. \n"+
		"Server version: alpha PingCAP Community Server - GPL\n\n"+
		"---------------------------------------------------\n"+
		"PD version: %s\nBuild ts: %s\nGit hash: %s\n"+
		"---------------------------------------------------\n",
		pd.Version, pd.Build_ts, pd.Git_hash)
	fmt.Println()

}

func bye() string {
	return "GoodBye!"
}
