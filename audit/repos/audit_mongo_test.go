package lxAuditRepos_test

import (
	"github.com/globalsign/mgo/bson"
	"github.com/litixsoft/lx-golib/audit"
	"github.com/litixsoft/lx-golib/audit/repos"
	"github.com/litixsoft/lx-golib/db"
	"github.com/litixsoft/lx-golib/helper"
	"github.com/litixsoft/lx-golib/tests/fixtures"
	"github.com/smartystreets/goconvey/convey"
	"log"
	"reflect"
	"testing"
	"time"
)

const (
	ServiceName     = "LxGoLib_Test_Service"
	ServiceHost     = "localhost:3101"
	AuditCollection = "audit"
)

func TestNewAuditMongo(t *testing.T) {
	// Db connect
	conn := fixtures.GetMongoConn()
	defer conn.Close()

	// Db base repo
	db := lxDb.NewMongoDb(conn, fixtures.TestDbName, AuditCollection)

	convey.Convey("Given db base repo", t, func() {
		convey.Convey("When create audit mongo repo", func() {
			repo := lxAuditRepos.NewAuditMongo(db, ServiceName, ServiceHost)

			convey.Convey("Then type should be *lxAuditRepos.auditMongo", func() {
				chkT := reflect.TypeOf(repo)
				convey.So(chkT.String(), convey.ShouldEqual, "*lxAuditRepos.auditMongo")
			})
		})
	})
}

func TestAuditMongo_SetupAudit(t *testing.T) {
	// Db connect
	conn := fixtures.GetMongoConn()
	defer conn.Close()

	// Delete collection
	conn.DB(fixtures.TestDbName).C(AuditCollection).DropCollection()

	// Test audit entry for create db and check indexes
	testEntry := lxAudit.AuditModel{
		TimeStamp:   time.Now(),
		ServiceName: ServiceName,
		ServiceHost: ServiceHost,
		User:        "TestUser",
		Message:     "TestMessage",
		Data: struct {
			Name  string
			Age   int
			Login time.Time
		}{
			Name:  "Test User",
			Age:   44,
			Login: time.Now(),
		},
	}
	if err := conn.DB(fixtures.TestDbName).C(AuditCollection).Insert(testEntry); err != nil {
		log.Fatal(err)
	}

	// Db base, repo
	db := lxDb.NewMongoDb(conn, fixtures.TestDbName, AuditCollection)
	repo := lxAuditRepos.NewAuditMongo(db, "TestService", "localhost:3101")

	convey.Convey("Given repo with deleted collection and one test entry", t, func() {
		convey.Convey("When check indexes before setup", func() {
			idx, err := conn.DB(fixtures.TestDbName).C(AuditCollection).Indexes()
			if err != nil {
				log.Fatal(err)
			}

			convey.Convey("Then indexes should be only contain _id", func() {
				convey.So(len(idx), convey.ShouldEqual, 1)
				convey.So(idx[0].Name, convey.ShouldEqual, "_id_")
			})
		})
		convey.Convey("When setup repo", func() {
			convey.So(repo.SetupAudit(), convey.ShouldBeNil)

			convey.Convey("Then indexes should be equal to expected values", func() {
				idx, err := conn.DB(fixtures.TestDbName).C(AuditCollection).Indexes()
				if err != nil {
					log.Fatal(err)
				}

				convey.So(len(idx), convey.ShouldEqual, 2)
				convey.So(idx[1].Name, convey.ShouldEqual, "timestamp_1")
			})
		})
	})
}

func TestAuditMongo_Log(t *testing.T) {
	// Db connect
	conn := fixtures.GetMongoConn()
	defer conn.Close()

	// Delete collection
	conn.DB(fixtures.TestDbName).C(AuditCollection).DropCollection()

	// Repo and setup
	// Db base, repo
	db := lxDb.NewMongoDb(conn, fixtures.TestDbName, AuditCollection)
	repo := lxAuditRepos.NewAuditMongo(db, "TestService", "localhost:3101")

	if err := repo.SetupAudit(); err != nil {
		log.Fatal(err)
	}

	convey.Convey("Given repo with deleted collection and indexes", t, func() {
		convey.Convey("When log a new entry", func() {
			done := repo.Log("test_user", "a audit message", lxHelper.M{"name": "test_name"})

			// wait for go routine is done
			<-done

			convey.Convey("Then audit entry should be found in db", func() {
				var result lxAudit.AuditModel
				if err := conn.DB(db.Name).C(db.Collection).Find(lxHelper.M{"user": "test_user"}).One(&result); err != nil {
					log.Fatal(err)
				}

				convey.So(result.ServiceName, convey.ShouldEqual, "TestService")
				convey.So(result.ServiceHost, convey.ShouldEqual, "localhost:3101")
				convey.So(result.User, convey.ShouldEqual, "test_user")
				convey.So(result.Message, convey.ShouldEqual, "a audit message")
				convey.So(result.Data, convey.ShouldResemble, bson.M{"name": "test_name"})

			})
		})
	})
}
