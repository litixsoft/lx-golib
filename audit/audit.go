package lxAudit

import "time"

// AuditModel, model for audit entry
type AuditModel struct {
	TimeStamp   time.Time   `json:"timestamp"`
	ServiceName string      `json:"service_name"`
	ServiceHost string      `json:"service_host"`
	User        interface{} `json:"user"`
	Message     interface{} `json:"msg"`
	Data        interface{} `json:"data"`
}
