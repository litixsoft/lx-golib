package lxHelper

import (
	"github.com/litixsoft/lx-golib/db"
	"encoding/json"
)

type M map[string]interface{}

type ReqByQuery struct {
	Options lxDb.Options `json:"opts,omitempty"`
	Query   M            `json:"query"`
}

func NewReqByQuery(opts string) (*ReqByQuery, error) {
	var data ReqByQuery

	if len(opts) > 0 {
		err := json.Unmarshal([]byte(opts), &data)
		if err != nil {
			return nil, err
		}
	}

	return &data, nil
}