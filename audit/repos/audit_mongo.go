package lxAuditRepos

import (
	"github.com/litixsoft/lx-golib/audit"
	"github.com/litixsoft/lx-golib/db"
	"github.com/globalsign/mgo"
	"time"
	"log"
)

type auditMongo struct {
	serviceName string
	serviceHost string
	db *lxDb.MongoDb
}

func NewAuditMongo(db *lxDb.MongoDb, serviceName, serviceHost string) lxAudit.IAudit {
	return &auditMongo{db:db, serviceName:serviceName, serviceHost:serviceHost}
}

func (repo *auditMongo) SetupAudit() error {
	// Copy mongo session (thread safe) and close after function
	conn := repo.db.Conn.Copy()
	defer conn.Close()

	// Setup indexes
	return repo.db.Setup([]mgo.Index{
		{Key: []string{"timestamp"}},
	})
}

func (repo *auditMongo) Log(user, message, data interface{}) chan bool {
	// channel for done
	done := make(chan bool, 1)

	go func() {
		// Copy mongo session (thread safe) and close after function
		conn := repo.db.Conn.Copy()
		defer conn.Close()

		entry := &lxAudit.AuditModel{
			TimeStamp:time.Now(),
			ServiceName:repo.serviceName,
			ServiceHost:repo.serviceHost,
			User: user,
			Message: message,
			Data:data,
		}

		// Insert entry
		if err := conn.DB(repo.db.Name).C(repo.db.Collection).Insert(entry); err != nil {
			log.Printf("mongoDb can't insert audit entry, error: %v\n", err)
		}

		done <- true
	}()

	return done
}