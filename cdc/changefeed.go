package cdc

import (
	"encoding/json"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/list"
	"tikv-client/pd"
	"tikv-client/util"
)

type Changefeed []struct {
	ID             string      `json:"id"`
	State          string      `json:"state"`
	CheckpointTso  int64       `json:"checkpoint_tso"`
	CheckpointTime string      `json:"checkpoint_time"`
	Error          interface{} `json:"error"`
}

type changefeedDetail struct {
	ID             string      `json:"id"`
	SinkURI        string      `json:"sink_uri"`
	CreateTime     string      `json:"create_time"`
	StartTs        int64       `json:"start_ts"`
	TargetTs       int         `json:"target_ts"`
	CheckpointTso  int64       `json:"checkpoint_tso"`
	CheckpointTime string      `json:"checkpoint_time"`
	SortEngine     string      `json:"sort_engine"`
	State          string      `json:"state"`
	Error          interface{} `json:"error"`
	ErrorHistory   interface{} `json:"error_history"`
	CreatorVersion string      `json:"creator_version"`
	TaskStatus     []struct {
		CaptureID       string `json:"capture_id"`
		TableIds        []int  `json:"table_ids"`
		TableOperations struct {
		} `json:"table_operations"`
	} `json:"task_status"`
}

type changefeedExecutor struct {
	id     string
	action string
	Changefeed
	host string
}

func (c changefeedExecutor) rebalanceTable() error {
	if err := pd.HttpPost(fmt.Sprintf("http://%s/api/v1/changefeeds/%s/tables/rebalance_table", c.host, c.id)); err != nil {
		return err
	}
	util.RunSuccess(c.action)
	return nil
}

func (c changefeedExecutor) resignOwner() error {
	//TODO implement me
	panic("implement me")
}

func changefeedBuild(c *Cdc) cmdExecutor {
	e := changefeedExecutor{
		id:         c.arg01,
		action:     c.arg02,
		host:       c.Captures[0].Address,
		Changefeed: c.Changefeed,
	}
	return e
}

func (c changefeedExecutor) query() error {
	for _, cf := range c.Changefeed {
		if cf.ID == c.id {
			var title []interface{}
			title = append(title, "id")
			title = append(title, "state")
			title = append(title, "checkpoint_tso")
			title = append(title, "checkpoint_time")
			title = append(title, "error")
			tbl := util.NewNormalDisplayTable(title)
			var row []interface{}
			row = append(row, cf.ID)
			row = append(row, cf.State)
			row = append(row, cf.CheckpointTso)
			row = append(row, cf.CheckpointTime)
			row = append(row, util.IsNil(cf.Error))
			tbl.AppendRow(row)
			tbl.Render()
			return nil
		}
	}
	return errCanNotFind(c.id, Changefeed_)
}

func (c changefeedExecutor) queryDetail() error {
	r, err := pd.HttpGet(fmt.Sprintf("http://%s/api/v1/changefeeds/%s", c.host, c.id))
	if err != nil {
		return err
	}
	var cd changefeedDetail
	err = json.Unmarshal(r, &cd)
	if err != nil {
		return err
	}
	l := list.NewWriter()
	l.SetStyle(list.StyleConnectedRounded)
	l.AppendItem(util.Yellow(c.id))
	l.Indent()
	l.AppendItems([]interface{}{
		fmt.Sprintf("sink_url: %s", cd.SinkURI),
		fmt.Sprintf("create_time: %s", cd.CreateTime),
		fmt.Sprintf("start_ts: %d", cd.StartTs),
		fmt.Sprintf("target_ts: %d", cd.TargetTs),
		fmt.Sprintf("checkpoint_tso: %d", cd.CheckpointTso),
		fmt.Sprintf("checkpoint_time: %s", cd.CheckpointTime),
		fmt.Sprintf("sort_engine: %s", cd.SortEngine),
		fmt.Sprintf("state: %v", util.Status(cd.State)),
		fmt.Sprintf("error: %v", util.IsNil(cd.Error)),
		fmt.Sprintf("error_history: %v", util.IsNil(cd.ErrorHistory)),
		"task_status"})
	l.Indent()
	for _, r := range cd.TaskStatus {
		l.AppendItem("")
		l.AppendItems([]interface{}{
			fmt.Sprintf("capture_id: %s", r.CaptureID),
			fmt.Sprintf("table_id: %d", util.IsNil(r.TableIds)),
			fmt.Sprintf("table_operations: %s", util.IsNil(r.TableOperations))})
		l.AppendItem("")
	}
	fmt.Println()
	fmt.Println(l.Render())
	fmt.Println()
	return nil
}

func getChangefeeds(captureHost string) (Changefeed, error) {
	r, err := pd.HttpGet(fmt.Sprintf("http://%s/api/v1/changefeeds", captureHost))
	if err != nil {
		return nil, err
	}
	var cf Changefeed
	if err = json.Unmarshal(r, &cf); err != nil {
		return nil, err
	}
	return cf, nil
}
