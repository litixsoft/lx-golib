package lxDb_test

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"log"
	"io/ioutil"
	"encoding/json"
	"testing"
	"github.com/smartystreets/goconvey/convey"
	"sort"
	"os"
	"github.com/litixsoft/lx-golib/db"
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
	Db   = "lx_golib_test"
	Coll = "users"
)

func getConn () *mgo.Session {

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
	conn.DB(Db).C(Coll).DropCollection()

	// Setup indexes
	indexes := []mgo.Index{
		{Key: []string{"name"}},
		{Key: []string{"email"}, Unique: true},
	}

	// Ensure indexes
	col := conn.DB(Db).C(Coll)
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
		if err := conn.DB(Db).C(Coll).Insert(users[i]); err != nil {
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

func TestMongoDb_Create(t *testing.T) {
	conn := getConn()
	defer conn.Close()

	// Tests
	db := &lxDb.MongoDb{
		Connection: conn,
		Name:   Db,
		Collection: Coll,
	}

	convey.Convey("Given a new user should be stored in the database", t, func() {
		convey.Convey("When create a correct new user", func() {
			tu := TestUser{Id: bson.NewObjectId(), Name: "Test User",Gender:"Male", Email: "t.user@gmail.com", IsActive:true}
			err := db.Create(tu)
			convey.So(err, convey.ShouldBeNil)

			convey.Convey("Then this user should be found in the database", func() {
				var chkResult TestUser
				err := db.Connection.DB(db.Name).C(db.Collection).Find(bson.M{"_id": tu.Id}).One(&chkResult)
				convey.So(err, convey.ShouldBeNil)
			})
		})
		convey.Convey("When we create a new incorrect user without id", func() {
			tu := TestUser{Name: "Test User",Gender:"Male", Email: "t.user@gmail.com", IsActive:true}

			convey.Convey("Then should be return a error", func() {
				err := db.Create(tu)
				convey.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMongoDb_GetAll(t *testing.T) {
	conn := getConn()
	defer conn.Close()

	// Create test users
	testUsers := setupData(conn)

	// Tests
	db := &lxDb.MongoDb{
		Connection: conn,
		Name:   Db,
		Collection: Coll,
	}

	convey.Convey("Given all users should be read from the database", t, func() {
		convey.Convey("When: get all users without query and options", func() {
			var result []TestUser
			var opts lxDb.Options
			n, err := db.GetAll(nil, &result, &opts)

			// Check err
			convey.So(err, convey.ShouldBeNil)

			convey.Convey("Then result should be contain all users", func() {
				// Expect all users
				var expect []TestUser
				for _, u := range testUsers {
					expect = append(expect, u)
				}

				// Check count and result
				convey.So(result, convey.ShouldResemble, expect)
				convey.So(n, convey.ShouldEqual, 0) // Count 0 when count option is false
			})
		})
		convey.Convey("When get all users with query is_active:true and option count", func() {
			var result []TestUser
			opts := lxDb.Options{Count: true}
			n, err := db.GetAll(bson.M{"is_active": true}, &result, &opts)

			// Check err
			convey.So(err, convey.ShouldBeNil)

			convey.Convey("Then result should be contain all active users", func() {
				// Expect all active users
				var expect []TestUser
				for _, u := range testUsers {
					if u.IsActive {
						expect = append(expect, u)
					}
				}

				// Check count and result
				convey.So(result, convey.ShouldResemble, expect)
				convey.So(n, convey.ShouldEqual, len(expect))
			})
		})
		convey.Convey("When get all users without query and options skip=5 limit=5 and count=true", func() {
			var result []TestUser
			opts := lxDb.Options{Skip: 5, Limit: 5, Count: true}
			n, err := db.GetAll(nil, &result, &opts)

			// Check err
			convey.So(err, convey.ShouldBeNil)

			convey.Convey("Then: result should be contain only users u2 and u3", func() {
				//sort.Slice(testUsers[:], func(i, j int) bool {
				//	return testUsers[i].Id < expect[j].Id
				//})
				var expect []TestUser
				i:=0
				for _, u := range testUsers {
					if i >4 && i <10 {
						expect = append(expect, u)
					}
					i++
				}

				// Check count and result
				convey.So(result, convey.ShouldResemble, expect)
				convey.So(n, convey.ShouldEqual, len(testUsers)) // Count complete skip is a filter
			})
		})
	})
}
//
//func TestMongoDb_GetCount(t *testing.T) {
//	// Setup
//	conn, err := mgo.Dial(dbHost)
//	if err != nil {
//		log.Fatal(err)
//	}
//	conn.SetMode(mgo.Monotonic, true)
//	defer conn.Close()
//
//	if err := baseSetup(conn); err != nil {
//		log.Fatal(err)
//	}
//
//	// Tests
//	rm := &lxBaseRepo.MongoDbBase{
//		Conn: conn,
//		Db:   Db,
//		Coll: Coll,
//	}
//
//	convey.Convey("Given: Test data with 4 users", t, func() {
//		convey.Convey("When: get count without query ", func() {
//			n, err := rm.GetCount(nil)
//
//			// Check err
//			convey.So(err, convey.ShouldBeNil)
//
//			convey.Convey("Then: response count should be len(testUsers)", func() {
//				convey.So(n, convey.ShouldEqual, len(baseTestUsers))
//			})
//		})
//		convey.Convey("When: get count with query .langbein@litixsoft.de", func() {
//			query := bson.M{"email": bson.M{"$regex": ".langbein@litixsoft.de"}}
//			n, err := rm.GetCount(query)
//
//			// Check err
//			convey.So(err, convey.ShouldBeNil)
//
//			convey.Convey("Then: response count should be 2", func() {
//				convey.So(n, convey.ShouldEqual, 2)
//			})
//		})
//	})
//}
//
//func TestMongoDb_GetOne(t *testing.T) {
//	// Setup
//	conn, err := mgo.Dial(dbHost)
//	if err != nil {
//		log.Fatal(err)
//	}
//	conn.SetMode(mgo.Monotonic, true)
//	defer conn.Close()
//
//	if err := baseSetup(conn); err != nil {
//		log.Fatal(err)
//	}
//
//	// Tests
//	rm := &lxBaseRepo.MongoDbBase{
//		Conn: conn,
//		Db:   Db,
//		Coll: Coll,
//	}
//
//	// Get user2 with id
//	var result TestUser
//	if assert.NoError(t, rm.GetOne(bson.M{"_id": baseTestUsers["u2"].Id}, &result)) {
//		assert.Equal(t, baseTestUsers["u2"], result)
//	}
//
//	// Get first .langbein@litixsoft.de user (Dennis)
//	result = TestUser{}
//	if assert.NoError(t, rm.GetOne(bson.M{"email": bson.M{"$regex": ".langbein@litixsoft.de"}}, &result)) {
//		assert.Equal(t, baseTestUsers["u4"], result)
//	}
//}
//
//func TestMongoDb_Update(t *testing.T) {
//	// Setup
//	conn, err := mgo.Dial(dbHost)
//	if err != nil {
//		log.Fatal(err)
//	}
//	conn.SetMode(mgo.Monotonic, true)
//	defer conn.Close()
//
//	if err := baseSetup(conn); err != nil {
//		log.Fatal(err)
//	}
//
//	// Tests
//	rm := &lxBaseRepo.MongoDbBase{
//		Conn: conn,
//		Db:   Db,
//		Coll: Coll,
//	}
//
//	// Update linus with id
//	assert.NoError(t, rm.Update(bson.M{"_id": baseTestUsers["u3"].Id}, bson.M{"name": "linus_updated"}))
//
//	// Check user linus should be updated
//	var chkResult TestUser
//	if assert.NoError(t, rm.Conn.DB(rm.Db).C(rm.Coll).Find(bson.M{"_id": baseTestUsers["u3"].Id}).One(&chkResult)) {
//		assert.Equal(t, "linus_updated", chkResult.Name)
//	}
//}
//
//func TestMongoDb_UpdateAll(t *testing.T) {
//	// Setup
//	conn, err := mgo.Dial(dbHost)
//	if err != nil {
//		log.Fatal(err)
//	}
//	conn.SetMode(mgo.Monotonic, true)
//	defer conn.Close()
//
//	if err := baseSetup(conn); err != nil {
//		log.Fatal(err)
//	}
//
//	// Tests
//	rm := &lxBaseRepo.MongoDbBase{
//		Conn: conn,
//		Db:   Db,
//		Coll: Coll,
//	}
//
//	// Update all users with email .langbein@litixsoft.de to NewName
//	info, err := rm.UpdateAll(bson.M{"email": bson.M{"$regex": ".langbein@litixsoft.de"}}, bson.M{"name": "NewName"})
//	if assert.NoError(t, err) {
//		expect := lxBaseRepo.ChangeInfo{
//			Updated: 2,
//			Removed: 0,
//			Matched: 2,
//		}
//		assert.Equal(t, expect, info)
//	}
//
//	// Check user linus and dennis should be updated
//	var chkResult []TestUser
//	if assert.NoError(t, rm.Conn.DB(rm.Db).C(rm.Coll).
//		Find(bson.M{"email": bson.M{"$regex": ".langbein@litixsoft.de"}}).All(&chkResult)) {
//		for _, u := range chkResult {
//			assert.Equal(t, "NewName", u.Name)
//		}
//	}
//}
//
//func TestMongoDb_Delete(t *testing.T) {
//	// Setup
//	conn, err := mgo.Dial(dbHost)
//	if err != nil {
//		log.Fatal(err)
//	}
//	conn.SetMode(mgo.Monotonic, true)
//	defer conn.Close()
//
//	if err := baseSetup(conn); err != nil {
//		log.Fatal(err)
//	}
//
//	// Tests
//	rm := &lxBaseRepo.MongoDbBase{
//		Conn: conn,
//		Db:   Db,
//		Coll: Coll,
//	}
//
//	// Delete xenia with loginName
//	assert.NoError(t, rm.Delete(bson.M{"login_name": "xenia"}))
//
//	// Check user xenia should be deleted
//	var chkResult TestUser
//	assert.Error(t, rm.Conn.DB(rm.Db).C(rm.Coll).Find(bson.M{"login_name": "xenia"}).One(&chkResult))
//}
//
//func TestMongoDb_DeleteAll(t *testing.T) {
//	// Setup
//	conn, err := mgo.Dial(dbHost)
//	if err != nil {
//		log.Fatal(err)
//	}
//	conn.SetMode(mgo.Monotonic, true)
//	defer conn.Close()
//
//	if err := baseSetup(conn); err != nil {
//		log.Fatal(err)
//	}
//
//	// Tests
//	rm := &lxBaseRepo.MongoDbBase{
//		Conn: conn,
//		Db:   Db,
//		Coll: Coll,
//	}
//
//	// Delete all users with email .liebetrau@litixsoft.de
//	info, err := rm.DeleteAll(bson.M{"email": bson.M{"$regex": ".liebetrau@litixsoft.de"}})
//	if assert.NoError(t, err) {
//		expect := lxBaseRepo.ChangeInfo{
//			Updated: 0,
//			Removed: 2,
//			Matched: 2,
//		}
//		assert.Equal(t, expect, info)
//	}
//
//	// Check user timo and xenia should be deleted
//	var chkResult []TestUser
//	assert.Empty(t, rm.Conn.DB(rm.Db).C(rm.Coll).
//		Find(bson.M{"email": bson.M{"$regex": ".liebetrau@litixsoft.de"}}).All(&chkResult))
//}
