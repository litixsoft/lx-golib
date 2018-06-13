package lxAudit

import "time"

type IAudit interface {
	SetupAudit() error
	Log(user, message, data interface{}) chan bool
}

type AuditModel struct {
	TimeStamp   time.Time   `json:"timestamp"`
	ServiceName string      `json:"service_name"`
	ServiceHost string      `json:"service_host"`
	User        interface{} `json:"user"`
	Message     interface{} `json:"msg"`
	Data        interface{} `json:"data"`
}
