package lxHelper

import "github.com/litixsoft/lx-golib/db"

type M map[string]interface{}

type ReqByQuery struct {
	Options lxDb.Options `json:"opts,omitempty"`
	Query   M         `json:"query"`
}