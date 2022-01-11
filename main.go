package main

import (
	"flag"
	"fmt"
	"github.com/c-bata/go-prompt"
	_ "github.com/pingcap/parser/test_driver"
	"os"
	"tikv-client/client"
	"tikv-client/util"
)

func main() {

	pdEndPoint := ParseArgs()
	// For self debug
	if pdEndPoint == "" {
		pdEndPoint = "172.16.101.3:2479"
	}
	c, err := client.NewCompleter(pdEndPoint)
	if err != nil {
		fmt.Printf(util.Red("Can't connect to placement driver\nError: %s\n"), err)
		os.Exit(0)
	}

	client.Welcome(*c.PlacementDriverGroup, client.Ticlient_)

	defer func() {
		fmt.Println("GoodBye!")
	}()

	prpt := prompt.New(
		c.Executor,
		c.Complete,
		prompt.OptionPrefix(">>>> "),
		prompt.OptionSuggestionBGColor(prompt.Red),
		prompt.OptionSuggestionTextColor(prompt.White),
		prompt.OptionDescriptionBGColor(prompt.White),
		prompt.OptionDescriptionTextColor(prompt.Black),
		prompt.OptionTitle(client.Ticlient_),
	)

	prpt.Run()

}

func ParseArgs() string {
	var pdEndPoint string
	flag.StringVar(&pdEndPoint, "pd", "", "Placement driver endpoint")
	flag.Parse()
	return pdEndPoint
}
