package fixtures

import (
	"os"
	"log"
	"github.com/globalsign/mgo"
)

const (
	TestDbName  = "lx_golib_test"
)

var DbHost string

// getConn, get a new db connection
func GetMongoConn() *mgo.Session {

	// Check DbHost environment
	dbHost := os.Getenv("DBHOST")

	// When not defined set default host
	if dbHost == "" {
		dbHost = "localhost:27017"
	}

	// Set dbHost
	DbHost = dbHost

	log.Println("DBHOST:", dbHost)

	// Create new connection
	conn, err := mgo.Dial(dbHost)
	if err != nil {
		log.Fatal(err)
	}
	conn.SetMode(mgo.Monotonic, true)

	return conn
}