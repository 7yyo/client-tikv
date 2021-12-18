package option

import (
	"github.com/c-bata/go-prompt"
	"tikv-client/syntax"
)

func PutOption(args []string) []prompt.Suggest {

	var s []prompt.Suggest

	switch len(args) {
	case 2:
		s = []prompt.Suggest{
			{Text: syntax.Into, Description: "No need to explain it."},
		}
	case 3:
		s = []prompt.Suggest{
			{Text: syntax.Tikv, Description: "We have defined that the table must be `tikv`."},
		}
	case 4:
		s = []prompt.Suggest{
			{Text: syntax.Values, Description: "No need to explain it."},
		}
	case 5:
		s = []prompt.Suggest{
			{Text: syntax.BracketsIn, Description: "No need to explain it."},
		}
	default:
		if len(args) != 1 {
			s = []prompt.Suggest{
				{Text: "," + syntax.BracketsIn, Description: "No need to explain it."},
			}
		}
	}
	return s
}
