package lxHelper

import (
	"github.com/litixsoft/lx-golib/db"
)

type M map[string]interface{}

type ReqGetData struct {
	Options lxDb.Options `json:"opts, omitempty"`
	Query   M            `json:"query"`
}

type ReqData struct {
	Query M `json:"query"`
	Data  M `json:"query"`
}
