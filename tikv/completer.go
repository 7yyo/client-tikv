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

	if d.GetWordBeforeCursor() == "" {
		return []prompt.Suggest{}
	}

	args := specificationArgs(d)

	if len(args) == 1 {
		s = []prompt.Suggest{
			{Text: syntax.Get, Description: "Get kv pairs from tikv."},
			{Text: syntax.Put, Description: "Put kv pairs to tikv."},
			{Text: syntax.Exit, Description: "Exit tikv-client."},
			{Text: syntax.Quit, Description: "Exit tikv-client."},
			{Text: syntax.Q, Description: "Exit tikv-client."},
		}
	}

	if len(args) >= 2 {
		switch args[0] {
		case syntax.Get:
			c.operateType = syntax.Get
		case syntax.Put:
			c.operateType = syntax.Put
		}
	}

	switch c.operateType {
	case syntax.Get:
		s = option.GetOption(args)
		c.clearOperate(8, args)
	case syntax.Put:
		s = option.PutOption(args)
		c.clearOperate(5, args)
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

func (c *Completer) clearOperate(n int, args []string) {
	if n == len(args) {
		c.operateType = ""
	}
}
