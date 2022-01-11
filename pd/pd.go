package pd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type PlacementDriverGroup struct {
	Header struct {
		ClusterID int64 `json:"cluster_id"`
	} `json:"header"`

	Members []struct {
		Name          string   `json:"name"`
		MemberID      int64    `json:"member_id"`
		PeerUrls      []string `json:"peer_urls"`
		ClientUrls    []string `json:"client_urls"`
		DeployPath    string   `json:"deploy_path"`
		BinaryVersion string   `json:"binary_version"`
		GitHash       string   `json:"git_hash"`
	} `json:"members"`

	Leader struct {
		Name          string   `json:"name"`
		MemberID      int64    `json:"member_id"`
		PeerUrls      []string `json:"peer_urls"`
		ClientUrls    []string `json:"client_urls"`
		DeployPath    string   `json:"deploy_path"`
		BinaryVersion string   `json:"binary_version"`
		GitHash       string   `json:"git_hash"`
	} `json:"leader"`

	EtcdLeader struct {
		Name          string   `json:"name"`
		MemberID      int64    `json:"member_id"`
		PeerUrls      []string `json:"peer_urls"`
		ClientUrls    []string `json:"client_urls"`
		DeployPath    string   `json:"deploy_path"`
		BinaryVersion string   `json:"binary_version"`
		GitHash       string   `json:"git_hash"`
	} `json:"etcd_leader"`
}

func PlacementDriverInfo(pdEndPoint string) (*PlacementDriverGroup, error) {
	members, err := HttpGet(fmt.Sprintf("http://%s/pd/api/v1/members", pdEndPoint))
	if err != nil {
		return nil, err
	}
	var pdGroup PlacementDriverGroup
	err = json.Unmarshal(members, &pdGroup)
	return &pdGroup, nil
}

func HttpGet(cmd string) ([]byte, error) {
	c := &http.Client{}
	req, err := http.NewRequest("GET", cmd, nil)
	if err != nil {
		return nil, err
	}
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func HttpPost(cmd string) error {
	c := &http.Client{}
	req, err := http.NewRequest("POST", cmd, nil)
	if err != nil {
		return err
	}
	if _, err = c.Do(req); err != nil {
		return err
	}
	return nil
}
