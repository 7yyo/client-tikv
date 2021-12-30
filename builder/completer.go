package builder

import (
	"flag"
	"github.com/c-bata/go-prompt"
	"github.com/pingcap/log"
	"github.com/tikv/client-go/v2/config"
	"github.com/tikv/client-go/v2/tikv"
	p "github.com/tikv/pd/client"
	"go.uber.org/zap/zapcore"
	"strings"
	pd2 "tikv-client/pd"
)

var s []prompt.Suggest

type Completer struct {
	Client      *tikv.RawKVClient
	operateType string
	pdEndPoint  []string
	Pd          pd2.Pd
}

func NewCompleter(pdEndPoint []string) (*Completer, error) {
	c, err := newKvClient(pdEndPoint)
	p := pd2.PdInfo(pdEndPoint[0])
	if err != nil {
		return nil, err
	}
	return &Completer{
		Client:     c,
		pdEndPoint: pdEndPoint,
		Pd:         *p,
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

func ParseArgs() []string {
	var pd string
	flag.StringVar(&pd, "pd", "", "pd endpoints")
	flag.Parse()
	return strings.Split(pd, ",")
}

func newKvClient(pdEndPoint []string) (*tikv.RawKVClient, error) {

	// In order to don't print log from console-go, set log level = panic.
	setLogLevel(zapcore.PanicLevel)

	var c *tikv.RawKVClient
	var err error
	if c, err = tikv.NewRawKVClient(pdEndPoint, config.DefaultConfig().Security, p.WithMaxErrorRetry(1)); err != nil {
		return nil, err
	}

	return c, nil
}

func setLogLevel(level zapcore.Level) {
	log.SetLevel(level)
}
