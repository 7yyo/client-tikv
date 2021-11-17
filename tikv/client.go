package tikv

import (
	"flag"
	"fmt"
	"github.com/modood/table"
	"github.com/tikv/client-go/v2/config"
	"github.com/tikv/client-go/v2/tikv"
	p "github.com/tikv/pd/client"
	"strings"
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

func (c *Completer) Put(kv *KvPair) (string, error) {
	t := time.Now()
	err := c.client.Put([]byte(kv.Key), []byte(kv.Value))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Query OK, 1 row affected (%f sec)", time.Now().Sub(t).Seconds()), nil
}

func (c *Completer) Get(kv *KvPair) (string, error) {
	t := time.Now()
	value, err := c.client.Get([]byte(kv.Key))
	if err != nil {
		return "", err
	}
	if value == nil {
		return fmt.Sprintf("Empty set (0.00 sec)"), nil
	}
	result := []KvPair{
		{kv.Key, string(value)},
	}
	fmt.Println(table.Table(result))
	return fmt.Sprintf("1 rows in set (%f sec)", time.Now().Sub(t).Seconds()), nil
}

func (c *Completer) Delete(kv *KvPair) (string, error) {
	t := time.Now()
	v, err := c.client.Get([]byte(kv.Key))
	if v == nil {
		return fmt.Sprintf("Query OK, 0 row affected (%f sec)", time.Now().Sub(t).Seconds()), nil
	}
	err = c.client.Delete([]byte(kv.Key))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Query OK, 1 row affected (%f sec)", time.Now().Sub(t).Seconds()), nil
}
