package lxAuditRepos

// IAudit, interface for audit repositories
type IAudit interface {
	SetupAudit() error
	Log(user, message, data interface{}) chan bool
}
