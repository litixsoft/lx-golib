package lxDb

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

func ConnectMongoDB(dbHost string) (*mgo.Session, error) {
	conn, err := mgo.Dial(dbHost)
	if err != nil {
		return nil, err
	}
	conn.SetMode(mgo.Monotonic, true)

	return conn, nil
}

// Db struct for mongodb
type MongoDb struct {
	Connection *mgo.Session
	Name       string
	Collection string
}

// Create, create new entity in collection
func (db MongoDb) Create(data interface{}) error {
	// Copy mongo session (thread safe) and close after function
	conn := db.Connection.Copy()
	defer conn.Close()

	// Insert data
	return conn.DB(db.Name).C(db.Collection).Insert(data)
}

// GetAll, get all entities by query in collection
func (db *MongoDb) GetAll(query interface{}, result interface{}, opts *Options) (int, error) {
	// Copy mongo session (thread safe) and close after function
	conn := db.Connection.Copy()
	defer conn.Close()

	// Set default count
	n := 0

	// Check if activate counter in options
	if opts.Count {
		var err error
		n, err = conn.DB(db.Name).C(db.Collection).Find(query).Count()
		if err != nil {
			return n, err
		}
	}

	// Find all with query in collection
	return n, conn.DB(db.Name).C(db.Collection).Find(query).Skip(opts.Skip).Limit(opts.Limit).All(result)
}

// GetCount, get count of entities by query in collection
func (db *MongoDb) GetCount(query interface{}) (int, error) {
	// Copy mongo session (thread safe) and close after function
	conn := db.Connection.Copy()
	defer conn.Close()

	// Find all with query in collection
	return conn.DB(db.Name).C(db.Collection).Find(query).Count()
}

// GetOne, get one entity by query in collection
func (db *MongoDb) GetOne(query interface{}, result interface{}) error {
	// Copy mongo session (thread safe) and close after function
	conn := db.Connection.Copy()
	defer conn.Close()

	// Find one with query in collection
	return conn.DB(db.Name).C(db.Collection).Find(query).One(result)
}

// Update, update one matched entity by query in collection
func (db *MongoDb) Update(query interface{}, data interface{}) error {
	// Copy mongo session (thread safe) and close after function
	conn := db.Connection.Copy()
	defer conn.Close()

	// Update one with query in collection
	return conn.DB(db.Name).C(db.Collection).Update(query, bson.M{"$set": data})
}

// UpdateAll, update all matched entities by query in collection
func (db *MongoDb) UpdateAll(query interface{}, data interface{}) (ChangeInfo, error) {
	// Copy mongo session (thread safe) and close after function
	conn := db.Connection.Copy()
	defer conn.Close()

	// Update all with query in collection
	info, err := conn.DB(db.Name).C(db.Collection).UpdateAll(query, bson.M{"$set": data})
	changeInfo := ChangeInfo{
		Updated: info.Updated,
		Removed: info.Removed,
		Matched: info.Matched,
	}

	return changeInfo, err
}

// Delete, delete one matched entity by query in collection
func (db *MongoDb) Delete(query interface{}) error {
	// Copy mongo session (thread safe) and close after function
	conn := db.Connection.Copy()
	defer conn.Close()

	// Delete one with query in collection
	return conn.DB(db.Name).C(db.Collection).Remove(query)
}

// DeleteAll, delete all matched entities by query in collection
func (db *MongoDb) DeleteAll(query interface{}) (ChangeInfo, error) {
	// Copy mongo session (thread safe) and close after function
	conn := db.Connection.Copy()
	defer conn.Close()

	// Remove all with query in collection
	info, err := conn.DB(db.Name).C(db.Collection).RemoveAll(query)
	changeInfo := ChangeInfo{
		Updated: info.Updated,
		Removed: info.Removed,
		Matched: info.Matched,
	}

	return changeInfo, err
}
