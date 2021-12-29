package tikv

import (
	"github.com/c-bata/go-prompt"
	"github.com/tikv/client-go/v2/tikv"
	"strings"
)

var s []prompt.Suggest

type Completer struct {
	client      *tikv.RawKVClient
	operateType string
	pdEndPoint  []string
}

func NewCompleter(pdEndPoint []string) (*Completer, error) {
	c, err := NewKvClient(pdEndPoint)
	if err != nil {
		return nil, err
	}
	return &Completer{
		client:     c,
		pdEndPoint: pdEndPoint,
	}, nil
}

func (c *Completer) Complete(d prompt.Document) []prompt.Suggest {

	if d.TextBeforeCursor() == "" {
		return []prompt.Suggest{}
	}

	args := specificationArgs(d)

	s = []prompt.Suggest{
		{Text: "select"},
		{Text: "insert into"},
		{Text: "from"},
		{Text: "where"},
		{Text: "values"},
		{Text: "order by"},
	}

	if len(args) > 1 {
		switch strings.ToUpper(args[len(args)-2]) {
		case "FROM":
			s = []prompt.Suggest{
				{Text: "tikv"},
				{Text: "regions"},
			}
		case "INTO":
			s = []prompt.Suggest{
				{Text: "tikv"},
			}
		case "SELECT":
			s = []prompt.Suggest{
				{Text: "tidb_parse_tso()"},
			}
		default:
		}
	}

	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)

}

func specificationArgs(d prompt.Document) []string {
	args := strings.Split(d.TextBeforeCursor(), " ")
	for i, v := range args {
		if v == " " {
			args = append(args[:i], args[i+1:]...)
		}
	}
	return args
}
