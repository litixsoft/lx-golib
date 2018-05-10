package lxHelper

import "github.com/litixsoft/lx-golib/db"

type M map[string]interface{}

type ReqGetData struct {
	Options lxDb.Options `json:"opts"`
	Query   M  `json:"query"`
}
