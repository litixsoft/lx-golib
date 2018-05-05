package lxBaseRepo_test

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/litixsoft/lx-golib/repos"
	"github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"sort"
	"testing"
	"time"
)

type TestUser struct {
	Id        bson.ObjectId `json:"id" bson:"_id"`
	Name      string        `json:"name" bson:"name"`
	LoginName string        `json:"login_name" bson:"login_name"`
	Email     string        `json:"email" bson:"email"`
}

type TestUser2 struct {
	Id bson.ObjectId `json:"id" bson:"_id"`
	Name string `json:"name" bson:"name"`
	Gender string `json:"gender" bson:"gender"`
	Email string `json:"email" bson:"email"`
	IsActive bool `json:"is_active" bson:"is_active"`
	LastActivity time.Time `json:"last_activity" bson:"last_activity"`

}

const (
	Db   = "bdu_repo_test"
	Coll = "users"
)

var (
	dbHost        string
	baseTestUsers = map[string]TestUser{
		"u1": {Id: bson.NewObjectId(), Name: "Timo Liebetrau", LoginName: "timo", Email: "t.liebetrau@litixsoft.de"},
		"u2": {Id: bson.NewObjectId(), Name: "Xenia Liebetrau", LoginName: "xenia", Email: "x.liebetrau@litixsoft.de"},
		"u3": {Id: bson.NewObjectId(), Name: "Linus Langbein", LoginName: "linus", Email: "l.langbein@litixsoft.de"},
		"u4": {Id: bson.NewObjectId(), Name: "Dennis Langbein", LoginName: "dennis", Email: "d.langbein@litixsoft.de"},
	}
)

func init() {
	// Check DbHost environment
	dbHost = os.Getenv("DBHOST")

	// When not defined set default host
	if dbHost == "" {
		dbHost = "mongodb://localhost:27017"
	}
}

func baseSetup(conn *mgo.Session) error {
	conn.DB(Db).C(Coll).DropCollection()

	indexes := []mgo.Index{
		{Key: []string{"email"}, Unique: true},
		{Key: []string{"login_name"}, Unique: true},
	}

	// Ensure indexes
	col := conn.DB(Db).C(Coll)
	for _, i := range indexes {
		if err := col.EnsureIndex(i); err != nil {
			return err
		}
	}

	// Save test data
	return conn.DB(Db).C(Coll).Insert(baseTestUsers["u1"], baseTestUsers["u2"], baseTestUsers["u3"], baseTestUsers["u4"])
}
//func TestMongoDb_Create2(t *testing.T) {
//	raw, err := ioutil.ReadFile("../tests/fixtures/MOCK_DATA.json")
//	if err != nil {
//		fmt.Println(err.Error())
//		os.Exit(1)
//	}
//
//	var c []TestUser2
//	if err := json.Unmarshal(raw, &c); err != nil {
//		log.Fatal(err)
//	}
//
//	for _, u := range c {
//		log.Println(u)
//	}
//
//}

func TestMongoDb_Create(t *testing.T) {
	// Setup
	conn, err := mgo.Dial(dbHost)
	if err != nil {
		log.Fatal(err)
	}
	conn.SetMode(mgo.Monotonic, true)
	defer conn.Close()

	if err := baseSetup(conn); err != nil {
		log.Fatal(err)
	}

	// Tests
	rm := &lxBaseRepo.MongoDb{
		Conn: conn,
		Db:   Db,
		Coll: Coll,
	}

	// Test create correct user
	tu := TestUser{Id: bson.NewObjectId(), Name: "Test User", LoginName: "tu", Email: "t.user@gmail.com"}
	assert.NoError(t, rm.Create(tu))

	// Check exists created user
	var chkResult TestUser
	if assert.NoError(t, rm.Conn.DB(rm.Db).C(rm.Coll).Find(bson.M{"_id": tu.Id}).One(&chkResult)) {
		assert.Equal(t, tu, chkResult)
	}

	// Test create incorrect user without id
	tu = TestUser{Name: "Test User2"}
	assert.Error(t, rm.Create(tu))
}

func TestMongoDb_GetAll(t *testing.T) {
	// Setup
	conn, err := mgo.Dial(dbHost)
	if err != nil {
		log.Fatal(err)
	}
	conn.SetMode(mgo.Monotonic, true)
	defer conn.Close()

	if err := baseSetup(conn); err != nil {
		log.Fatal(err)
	}

	// Tests
	rm := &lxBaseRepo.MongoDb{
		Conn: conn,
		Db:   Db,
		Coll: Coll,
	}

	convey.Convey("Given: Test data with 4 users", t, func() {
		convey.Convey("When: get all users without query and options", func() {
			var result []TestUser
			var opts lxBaseRepo.Options
			n, err := rm.GetAll(nil, &result, &opts)

			// Check err
			convey.So(err, convey.ShouldBeNil)

			// Sort result for compare
			sort.Slice(result[:], func(i, j int) bool {
				return result[i].Id < result[j].Id
			})

			convey.Convey("Then: result should be contain all users", func() {
				expect := []TestUser{baseTestUsers["u1"], baseTestUsers["u2"], baseTestUsers["u3"], baseTestUsers["u4"]}

				// Sort expect for compare
				sort.Slice(expect[:], func(i, j int) bool {
					return expect[i].Id < expect[j].Id
				})

				// Check count and result
				convey.So(result, convey.ShouldResemble, expect)
				convey.So(n, convey.ShouldEqual, 0) // Count 0 when count option is false
			})
		})
		convey.Convey("When: get all users with query .langbein@litixsoft.de and option count", func() {
			var result []TestUser
			opts := lxBaseRepo.Options{Count: true}
			n, err := rm.GetAll(bson.M{"email": bson.M{"$regex": ".langbein@litixsoft.de"}}, &result, &opts)

			// Check err
			convey.So(err, convey.ShouldBeNil)

			// Sort result for compare
			sort.Slice(result[:], func(i, j int) bool {
				return result[i].Id < result[j].Id
			})

			convey.Convey("Then: result should be contain only users u3 and u4", func() {
				expect := []TestUser{baseTestUsers["u3"], baseTestUsers["u4"]}

				// Sort expect for compare
				sort.Slice(expect[:], func(i, j int) bool {
					return expect[i].Id < expect[j].Id
				})

				// Check count and result
				convey.So(result, convey.ShouldResemble, expect)
				convey.So(n, convey.ShouldEqual, 2)
			})
		})
		convey.Convey("When: get all users without query and options skip=1 limit=2 and count=true", func() {
			var result []TestUser
			opts := lxBaseRepo.Options{Skip: 1, Limit: 2, Count: true}
			n, err := rm.GetAll(nil, &result, &opts)

			// Check err
			convey.So(err, convey.ShouldBeNil)

			// Sort result for compare
			sort.Slice(result[:], func(i, j int) bool {
				return result[i].Id < result[j].Id
			})

			convey.Convey("Then: result should be contain only users u2 and u3", func() {
				expect := []TestUser{baseTestUsers["u2"], baseTestUsers["u3"]}

				// Sort expect for compare
				sort.Slice(expect[:], func(i, j int) bool {
					return expect[i].Id < expect[j].Id
				})

				// Check count and result
				convey.So(result, convey.ShouldResemble, expect)
				convey.So(n, convey.ShouldEqual, 4)
			})
		})
	})
}

func TestMongoDb_GetCount(t *testing.T) {
	// Setup
	conn, err := mgo.Dial(dbHost)
	if err != nil {
		log.Fatal(err)
	}
	conn.SetMode(mgo.Monotonic, true)
	defer conn.Close()

	if err := baseSetup(conn); err != nil {
		log.Fatal(err)
	}

	// Tests
	rm := &lxBaseRepo.MongoDb{
		Conn: conn,
		Db:   Db,
		Coll: Coll,
	}

	convey.Convey("Given: Test data with 4 users", t, func() {
		convey.Convey("When: get count without query ", func() {
			n, err := rm.GetCount(nil)

			// Check err
			convey.So(err, convey.ShouldBeNil)

			convey.Convey("Then: response count should be len(testUsers)", func() {
				convey.So(n, convey.ShouldEqual, len(baseTestUsers))
			})
		})
		convey.Convey("When: get count with query .langbein@litixsoft.de", func() {
			query := bson.M{"email": bson.M{"$regex": ".langbein@litixsoft.de"}}
			n, err := rm.GetCount(query)

			// Check err
			convey.So(err, convey.ShouldBeNil)

			convey.Convey("Then: response count should be 2", func() {
				convey.So(n, convey.ShouldEqual, 2)
			})
		})
	})
}

func TestMongoDb_GetOne(t *testing.T) {
	// Setup
	conn, err := mgo.Dial(dbHost)
	if err != nil {
		log.Fatal(err)
	}
	conn.SetMode(mgo.Monotonic, true)
	defer conn.Close()

	if err := baseSetup(conn); err != nil {
		log.Fatal(err)
	}

	// Tests
	rm := &lxBaseRepo.MongoDb{
		Conn: conn,
		Db:   Db,
		Coll: Coll,
	}

	// Get user2 with id
	var result TestUser
	if assert.NoError(t, rm.GetOne(bson.M{"_id": baseTestUsers["u2"].Id}, &result)) {
		assert.Equal(t, baseTestUsers["u2"], result)
	}

	// Get first .langbein@litixsoft.de user (Dennis)
	result = TestUser{}
	if assert.NoError(t, rm.GetOne(bson.M{"email": bson.M{"$regex": ".langbein@litixsoft.de"}}, &result)) {
		assert.Equal(t, baseTestUsers["u4"], result)
	}
}

func TestMongoDb_Update(t *testing.T) {
	// Setup
	conn, err := mgo.Dial(dbHost)
	if err != nil {
		log.Fatal(err)
	}
	conn.SetMode(mgo.Monotonic, true)
	defer conn.Close()

	if err := baseSetup(conn); err != nil {
		log.Fatal(err)
	}

	// Tests
	rm := &lxBaseRepo.MongoDb{
		Conn: conn,
		Db:   Db,
		Coll: Coll,
	}

	// Update linus with id
	assert.NoError(t, rm.Update(bson.M{"_id": baseTestUsers["u3"].Id}, bson.M{"name": "linus_updated"}))

	// Check user linus should be updated
	var chkResult TestUser
	if assert.NoError(t, rm.Conn.DB(rm.Db).C(rm.Coll).Find(bson.M{"_id": baseTestUsers["u3"].Id}).One(&chkResult)) {
		assert.Equal(t, "linus_updated", chkResult.Name)
	}
}

func TestMongoDb_UpdateAll(t *testing.T) {
	// Setup
	conn, err := mgo.Dial(dbHost)
	if err != nil {
		log.Fatal(err)
	}
	conn.SetMode(mgo.Monotonic, true)
	defer conn.Close()

	if err := baseSetup(conn); err != nil {
		log.Fatal(err)
	}

	// Tests
	rm := &lxBaseRepo.MongoDb{
		Conn: conn,
		Db:   Db,
		Coll: Coll,
	}

	// Update all users with email .langbein@litixsoft.de to NewName
	info, err := rm.UpdateAll(bson.M{"email": bson.M{"$regex": ".langbein@litixsoft.de"}}, bson.M{"name": "NewName"})
	if assert.NoError(t, err) {
		expect := lxBaseRepo.ChangeInfo{
			Updated: 2,
			Removed: 0,
			Matched: 2,
		}
		assert.Equal(t, expect, info)
	}

	// Check user linus and dennis should be updated
	var chkResult []TestUser
	if assert.NoError(t, rm.Conn.DB(rm.Db).C(rm.Coll).
		Find(bson.M{"email": bson.M{"$regex": ".langbein@litixsoft.de"}}).All(&chkResult)) {
		for _, u := range chkResult {
			assert.Equal(t, "NewName", u.Name)
		}
	}
}

func TestMongoDb_Delete(t *testing.T) {
	// Setup
	conn, err := mgo.Dial(dbHost)
	if err != nil {
		log.Fatal(err)
	}
	conn.SetMode(mgo.Monotonic, true)
	defer conn.Close()

	if err := baseSetup(conn); err != nil {
		log.Fatal(err)
	}

	// Tests
	rm := &lxBaseRepo.MongoDb{
		Conn: conn,
		Db:   Db,
		Coll: Coll,
	}

	// Delete xenia with loginName
	assert.NoError(t, rm.Delete(bson.M{"login_name": "xenia"}))

	// Check user xenia should be deleted
	var chkResult TestUser
	assert.Error(t, rm.Conn.DB(rm.Db).C(rm.Coll).Find(bson.M{"login_name": "xenia"}).One(&chkResult))
}

func TestMongoDb_DeleteAll(t *testing.T) {
	// Setup
	conn, err := mgo.Dial(dbHost)
	if err != nil {
		log.Fatal(err)
	}
	conn.SetMode(mgo.Monotonic, true)
	defer conn.Close()

	if err := baseSetup(conn); err != nil {
		log.Fatal(err)
	}

	// Tests
	rm := &lxBaseRepo.MongoDb{
		Conn: conn,
		Db:   Db,
		Coll: Coll,
	}

	// Delete all users with email .liebetrau@litixsoft.de
	info, err := rm.DeleteAll(bson.M{"email": bson.M{"$regex": ".liebetrau@litixsoft.de"}})
	if assert.NoError(t, err) {
		expect := lxBaseRepo.ChangeInfo{
			Updated: 0,
			Removed: 2,
			Matched: 2,
		}
		assert.Equal(t, expect, info)
	}

	// Check user timo and xenia should be deleted
	var chkResult []TestUser
	assert.Empty(t, rm.Conn.DB(rm.Db).C(rm.Coll).
		Find(bson.M{"email": bson.M{"$regex": ".liebetrau@litixsoft.de"}}).All(&chkResult))
}
