package lxDb_test

import (
	"encoding/json"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/litixsoft/lx-golib/db"
	"github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"sort"
	"testing"
)

// TestUser, struct for test users
type TestUser struct {
	Id       bson.ObjectId `json:"id" bson:"_id"`
	Name     string        `json:"name" bson:"name"`
	Gender   string        `json:"gender" bson:"gender"`
	Email    string        `json:"email" bson:"email"`
	IsActive bool          `json:"is_active" bson:"is_active"`
}

const (
	TestDbName     = "lx_golib_test"
	TestCollection = "users"
)

func getConn() *mgo.Session {

	// Check DbHost environment
	dbHost := os.Getenv("DBHOST")

	// When not defined set default host
	if dbHost == "" {
		dbHost = "localhost:27017"
	}

	log.Println("DBHOST:", dbHost)

	// Create new connection
	conn, err := mgo.Dial(dbHost)
	if err != nil {
		log.Fatal(err)
	}
	conn.SetMode(mgo.Monotonic, true)

	return conn
}

// setupData, create the test data and prepare the database
func setupData(conn *mgo.Session) []TestUser {
	// Delete collection if exists
	conn.DB(TestDbName).C(TestCollection).DropCollection()

	// Setup indexes
	indexes := []mgo.Index{
		{Key: []string{"name"}},
		{Key: []string{"email"}, Unique: true},
	}

	// Ensure indexes
	col := conn.DB(TestDbName).C(TestCollection)
	for _, i := range indexes {
		if err := col.EnsureIndex(i); err != nil {
			log.Fatal(err)
		}
	}

	// Load test data from json file
	raw, err := ioutil.ReadFile("../tests/fixtures/MOCK_DATA.json")
	if err != nil {
		log.Fatal(err)
	}

	// Convert
	var users []TestUser
	if err := json.Unmarshal(raw, &users); err != nil {
		log.Fatal(err)
	}

	// Make Test users map and insert test data in db
	for i := 0; i < len(users); i++ {
		users[i].Id = bson.NewObjectId()

		// Insert user
		if err := conn.DB(TestDbName).C(TestCollection).Insert(users[i]); err != nil {
			log.Fatal(err)
		}
	}

	// Sort test users with id for compare
	sort.Slice(users[:], func(i, j int) bool {
		return users[i].Id < users[j].Id
	})

	// Return mongo connection
	return users
}

func TestNewMongoDb(t *testing.T) {
	conn := getConn()
	defer conn.Close()

	convey.Convey("Given mongoDb connection", t, func() {
		convey.Convey("When create new mongoDb instance", func() {
			db := lxDb.NewMongoDb(conn, TestDbName, TestCollection)

			convey.Convey("Then type should be *lxDb.mongoDb", func() {
				chkT := reflect.TypeOf(db)
				convey.So(chkT.String(), convey.ShouldEqual, "*lxDb.mongoDb")
			})
			convey.Convey("And then test query should equal expected", func() {
				expected := setupData(conn)

				var result []TestUser
				db.Conn.DB(db.Name).C(db.Collection).Find(nil).All(&result)

				// Sort result for compare
				sort.Slice(result[:], func(i, j int) bool {
					return result[i].Id < result[j].Id
				})

				// Check result
				convey.So(result, convey.ShouldResemble, expected)
			})
		})
	})
}

func TestMongoDb_Setup(t *testing.T) {
	conn := getConn()
	defer conn.Close()

	// Delete collection if exists
	conn.DB(TestDbName).C(TestCollection).DropCollection()

	convey.Convey("Given mongoDb connection with drop collection", t, func() {
		convey.Convey("When db indexes setup", func() {
			db := lxDb.NewMongoDb(conn, TestDbName, TestCollection)
			db.Setup([]mgo.Index{
				{Key: []string{"name"}},
				{Key: []string{"email"}, Unique: true},
			})

			convey.Convey("Then index should be correct set", func() {
				idx, err := db.Conn.DB(db.Name).C(db.Collection).Indexes()
				convey.So(err, convey.ShouldBeNil)

				// Check indexes
				convey.So(idx[1].Name, convey.ShouldEqual, "email_1")
				convey.So(idx[1].Unique, convey.ShouldBeTrue)
				convey.So(idx[2].Name, convey.ShouldEqual, "name_1")
			})
		})
	})
}