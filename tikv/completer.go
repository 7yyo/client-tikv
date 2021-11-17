package tikv

import (
	"github.com/c-bata/go-prompt"
	"github.com/tikv/client-go/v2/tikv"
)

type Completer struct {
	client *tikv.RawKVClient
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
	if d.TextBeforeCursor() == "" {
		return []prompt.Suggest{}
	}
	s := []prompt.Suggest{
		{Text: "insert into tikv values", Description: "put kv to tikv."},
		{Text: "select * from tikv where k = ", Description: "get kv from tikv."},
		{Text: "delete from tikv where k = ", Description: "delete kv from tikv."},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}
