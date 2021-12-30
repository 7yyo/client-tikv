package pd

import (
	"encoding/json"
	"fmt"
	"os"
	http2 "tikv-client/http"
)

type Pd struct {
	Build_ts string
	Version  string
	Git_hash string
}

func PdInfo(pdEndPoint string) *Pd {
	body, err := http2.ReqGet(fmt.Sprintf(http2.Status, pdEndPoint))
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	var pd Pd
	err = json.Unmarshal(body, &pd)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	return &pd
}
