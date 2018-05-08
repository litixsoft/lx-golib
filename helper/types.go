package lxHelper

import (
	"github.com/litixsoft/lx-golib/db"
	"encoding/json"
)

type M map[string]interface{}

type queryStrConfig struct {
	Options lxDb.Options `json:"opts"`
	Query   interface{}  `json:"query"`
}

func NewQueryStrConfig(queryStr string) (*queryStrConfig, error) {
	var config queryStrConfig

	if len(queryStr) > 0 {
		err := json.Unmarshal([]byte(queryStr), &config)
		if err != nil {
			return nil, err
		}
	}

	return &config, nil
}
