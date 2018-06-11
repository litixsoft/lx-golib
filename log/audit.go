package lxLog

import "time"

type AuditModel struct {
	TimeStamp time.Time `json:"@timestamp"`
	User string `json:"user"`
	Host string `json:"host"`
	Type string `json:"type"`
	Data interface{} `json:"data"`
}

type IAudit interface {
	Log()
}


func Log() {

}