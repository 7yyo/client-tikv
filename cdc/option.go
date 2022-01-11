package cdc

import (
	"fmt"
	"github.com/c-bata/go-prompt"
)

const (
	Capture_        = "capture"
	Changefeed_     = "changefeed"
	Query_          = "query"
	Query_detail_   = "query-detail"
	Resign_owner    = "resign-owner"
	Rebalance_table = "rebalance-table"
	back_           = "back"
)

func (c *Cdc) CdcCompleter(args []string, d prompt.Document) []prompt.Suggest {
	var suggest []prompt.Suggest
	switch len(args) {
	case 1:
		suggest = []prompt.Suggest{
			{Text: Capture_, Description: "Get all the capture for ticdc."},
			{Text: Changefeed_, Description: "Get all the changefeed tasks for ticdc."},
			{Text: back_, Description: "Back to main menu."},
		}
	case 2:
		c.command.arg00 = args[0]
		switch args[0] {
		case Changefeed_:
			for _, cf := range c.Changefeed {
				suggest = append(suggest, prompt.Suggest{Text: cf.ID, Description: cf.State})
			}
		case Capture_:
			for _, c := range c.Captures {
				suggest = append(suggest, prompt.Suggest{Text: c.Address, Description: c.ID})
			}
		}
	case 3:
		c.command.arg01 = args[1]
		switch args[0] {
		case Capture_:
			suggest = []prompt.Suggest{
				{Text: Query_, Description: fmt.Sprintf("Get %s information", args[0])},
				{Text: Resign_owner, Description: "Expel the current owner node of TiCDC, and trigger a new round of election to generate a new owner node."},
			}
		case Changefeed_:
			suggest = []prompt.Suggest{
				{Text: Query_, Description: fmt.Sprintf("Get %s information", args[0])},
				{Text: Query_detail_, Description: fmt.Sprintf("Get %s detail information", args[0])},
				{Text: Rebalance_table, Description: fmt.Sprintf("Rebalance table for %s", args[0])},
			}
		}
	default:
		return suggest
	}
	return prompt.FilterHasPrefix(suggest, d.GetWordBeforeCursor(), true)
}
