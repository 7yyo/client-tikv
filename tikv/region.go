package tikv

import (
	"encoding/json"
	"github.com/jedib0t/go-pretty/v6/table"
	"reflect"
	"strconv"
	http2 "tikv-client/http"
	"tikv-client/syntax"
	"tikv-client/util"
	"time"
)

type Regions struct {
	Count   int
	Regions []Region
}

type Region struct {
	Id               int
	Start_key        string
	End_key          string
	Peers            []Peer
	Leader           Leader
	Read_bytes       int
	Written_bytes    int
	Read_keys        int
	Approximate_size int
	Approximate_keys int
}

type Leader struct {
	Id        int
	Store_id  int
	Role_name string
}

type Peer struct {
	Id        int
	Store_id  int
	Role_name string
}

type RegionTable struct {
	REGION_ID        int
	START_KEY        string
	END_KEY          string
	LEADER_ID        int
	LEADER_STORE_ID  int
	PEERS            string
	WRITTEN_BYTES    int
	READ_BYTES       int
	APPROXIMATE_SIZE int
	APPROXIMATE_KEYS int
}

func (c *Completer) GetRegionInfo(sql *syntax.SQL) (string, error) {

	t := time.Now()

	body, err := http2.ReqGet("http://%s/pd/api/v1/regions", c.pdEndPoint[0])
	if err != nil {
		return "", err
	}

	var regions Regions
	err = json.Unmarshal(body, &regions)
	if err != nil {
		return "", err
	}

	var rts []RegionTable
	var rt RegionTable

	for _, region := range regions.Regions {
		rt.REGION_ID = region.Id
		rt.START_KEY = region.Start_key
		rt.END_KEY = region.End_key
		rt.LEADER_ID = region.Leader.Id
		rt.LEADER_STORE_ID = region.Leader.Store_id
		if len(region.Peers) != 0 {
			rt.PEERS = strconv.Itoa(region.Peers[0].Id)
		}
		i := 0
		for _, p := range region.Peers {
			if i == 0 {
				i++
				continue
			}
			rt.PEERS += "," + strconv.Itoa(p.Id)
		}
		rt.WRITTEN_BYTES = region.Written_bytes
		rt.READ_BYTES = region.Read_bytes
		rt.APPROXIMATE_SIZE = region.Approximate_size
		rt.APPROXIMATE_KEYS = region.Approximate_keys
		rts = append(rts, rt)
	}

	var tbl table.Writer

	if len(rts) > 0 {

		var fn []interface{}
		tt := reflect.TypeOf(rts[0])

		for i := 0; i < tt.NumField(); i++ {
			fn = append(fn, tt.Field(i).Name)
		}

		tbl = util.NewNormalDisplayTable(fn)

	}

	for _, r := range rts {
		var row []interface{}
		row = append(row, r.REGION_ID)
		row = append(row, r.START_KEY)
		row = append(row, r.END_KEY)
		row = append(row, r.LEADER_ID)
		row = append(row, r.LEADER_STORE_ID)
		row = append(row, r.PEERS)
		row = append(row, r.WRITTEN_BYTES)
		row = append(row, r.READ_BYTES)
		row = append(row, r.APPROXIMATE_SIZE)
		row = append(row, r.APPROXIMATE_KEYS)
		tbl.AppendRow(row)
	}

	tbl.SortBy([]table.SortBy{{Name: sql.OrderBy, Mode: table.AscNumeric}})

	tbl.Render()

	return util.NRowsInSet(regions.Count, t), nil

}
