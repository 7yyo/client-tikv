package tikv

import (
	"flag"
	"fmt"
	"github.com/modood/table"
	"github.com/tikv/client-go/v2/config"
	"github.com/tikv/client-go/v2/tikv"
	p "github.com/tikv/pd/client"
	"strings"
	"tikv-client/syntax"
	"time"
)

func ParseArgs() []string {
	var pd string
	flag.StringVar(&pd, "pd", "", "pd endpoints")
	flag.Parse()
	return strings.Split(pd, ",")
}

func NewKvClient(pdEndPoint []string) (*tikv.RawKVClient, error) {
	c, err := tikv.NewRawKVClient(pdEndPoint, config.DefaultConfig().Security, p.WithMaxErrorRetry(1))
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Completer) Get(KvPairs *[]syntax.KvPair) (string, error) {
	var result []syntax.KvPair
	t := time.Now()
	for _, kv := range *KvPairs {
		value, err := c.client.Get([]byte(kv.Key))
		kv.Value = string(value)
		if err != nil {
			return "", err
		}
		if value == nil {
			return fmt.Sprintf("Empty set (0.00 sec)"), nil
		}
		result = append(result, kv)
	}
	fmt.Println(table.Table(result))
	return fmt.Sprintf("%d rows in set (%f sec)", len(*KvPairs), time.Now().Sub(t).Seconds()), nil
}
