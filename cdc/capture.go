package cdc

import (
	"context"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"tikv-client/pd"
	"tikv-client/util"
)

const (
	capture_version  = "version"
	capture_id       = "id"
	capture_address  = "address"
	capture_is_owner = "is_owner"
)

type Capture struct {
	ID      string `json:"id"`
	Address string `json:"address"`
	Version string `json:"version"`
	owner
}

type captureExecutor struct {
	address  string
	action   string
	Captures []Capture
}

func (c *captureExecutor) rebalanceTable() error {
	//TODO implement me
	panic("implement me")
}

func (c *captureExecutor) resignOwner() error {
	if err := pd.HttpPost(fmt.Sprintf("http://%s/api/v1/owner/resign", c.address)); err != nil {
		return err
	}
	util.RunSuccess(c.action)
	return nil
}

func (c captureExecutor) query() error {
	for _, cp := range c.Captures {
		if cp.Address == c.address {
			r, err := pd.HttpGet(fmt.Sprintf("http://%s/api/v1/status", cp.Address))
			if err != nil {
				return err
			}
			var o owner
			err = json.Unmarshal(r, &o)
			if err != nil {
				return err
			}
			var title []interface{}
			title = append(title, "version")
			title = append(title, "id")
			title = append(title, "address")
			title = append(title, "is_owner")
			tbl := util.NewNormalDisplayTable(title)
			var row []interface{}
			row = append(row, cp.Version)
			row = append(row, cp.ID)
			row = append(row, cp.Address)
			row = append(row, o.IsOwner)
			tbl.AppendRow(row)
			tbl.Render()
		}
	}
	return nil
}

func (c captureExecutor) queryDetail() error {
	if err := c.query(); err != nil {
		return err
	}
	return nil
}

func captureBuild(c *Cdc) cmdExecutor {
	e := captureExecutor{
		address:  c.arg01,
		action:   c.arg02,
		Captures: c.Captures,
	}
	return &e
}

type owner struct {
	IsOwner bool `json:"is_owner"`
}

func getCaptures(e *clientv3.Client) ([]Capture, error) {
	r, err := e.Get(context.TODO(), "/tidb/cdc/capture/", clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	var c Capture
	var cs []Capture
	for _, kv := range r.Kvs {
		err = json.Unmarshal(kv.Value, &c)
		if err != nil {
			return nil, err
		}
		cs = append(cs, c)
	}
	return cs, nil
}
