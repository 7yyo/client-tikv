package tikv

import (
	"github.com/c-bata/go-prompt"
	"github.com/tikv/client-go/v2/tikv"
	"strings"
	"tikv-client/option"
	"tikv-client/syntax"
)

var s []prompt.Suggest

type Completer struct {
	client      *tikv.RawKVClient
	operateType string
}

func NewCompleter(pdEndPoint []string) (*Completer, error) {
	c, err := NewKvClient(pdEndPoint)
	if err != nil {
		return nil, err
	}
	return &Completer{
		client: c,
	}, nil
}

func (c *Completer) Complete(d prompt.Document) []prompt.Suggest {

	args := specificationArgs(d)

	if len(args) == 1 {
		s = []prompt.Suggest{
			{Text: syntax.Get, Description: "Get the key-value pairs of the corresponding key from tikv."},
			{Text: syntax.Put, Description: "Insert key-value pairs into tikv, each values() should only have two values of kv."},
			{Text: syntax.Exit, Description: "Exit tikv-client."},
		}
	}

	if len(args) >= 2 {
		switch args[0] {
		case syntax.Get:
			c.operateType = syntax.Get
			s = option.GetOption(args)
		case syntax.Put:
			c.operateType = syntax.Put
			s = option.PutOption(args)
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
