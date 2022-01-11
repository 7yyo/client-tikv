package cdc

import (
	"go.etcd.io/etcd/clientv3"
	"tikv-client/pd"
)

const (
	Cdc_ = "cdc"
)

type Cdc struct {
	*pd.PlacementDriverGroup
	EtcdClient *clientv3.Client
	Captures   []Capture
	Changefeed
	command
}

// Init when choose cdc, flush ticdc cluster info all
func Init(pdGroup pd.PlacementDriverGroup, client *clientv3.Client) (*Cdc, error) {
	captures, err := getCaptures(client)
	if err != nil {
		return nil, err
	}
	changefeeds, err := getChangefeeds(captures[0].Address)
	if err != nil {
		return nil, err
	}
	cdcCluster := Cdc{
		PlacementDriverGroup: &pdGroup,
		EtcdClient:           client,
		Captures:             captures,
		Changefeed:           changefeeds,
	}
	return &cdcCluster, nil
}

func (c *Cdc) Run(a []string) error {
	var e cmdExecutor
	var err error
	if e, err = c.build(a); err != nil {
		return err
	}
	if err := exec(a, e); err != nil {
		return err
	}
	return nil
}

func (c *Cdc) build(a []string) (cmdExecutor, error) {
	// cdc changefeed(00) task01(01) query(02)
	if len(a) < 3 {
		return nil, errInvalidCommand(a)
	}
	c.arg02 = a[2]
	var e cmdExecutor
	switch c.arg00 {
	case Changefeed_:
		e = changefeedBuild(c)
	case Capture_:
		e = captureBuild(c)
	default:
		return nil, errInvalidCommand(a)
	}
	return e, nil
}

func exec(a []string, e cmdExecutor) error {
	var err error
	switch a[2] {
	case Query_:
		err = e.query()
	case Query_detail_:
		err = e.queryDetail()
	case Resign_owner:
		err = e.resignOwner()
	case Rebalance_table:
		err = e.rebalanceTable()
	default:
		return errInvalidCommand(a)
	}
	if err != nil {
		return err
	}
	return nil
}
