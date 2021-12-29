package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	Status = "http://%s/pd/api/v1/status"
	Tso    = "tiup ctl:%s pd -u http://%s tso %d"
)

func ReqGet(cmd string, pdEndPoint string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf(cmd, pdEndPoint), nil)
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
