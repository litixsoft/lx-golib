package lxDb

import (
	"errors"
	"github.com/globalsign/mgo"
)

// Db struct for mongodb
type MongoBaseDb struct {
	Connection *mgo.Session
	Name     string
	Collection string
}

func NewMongoBaseDb(connection *mgo.Session, dbName, collection string) *MongoBaseDb {
	return &MongoBaseDb{
		Connection: connection,
		Name:       dbName,
		Collection: collection,
	}
}

// Setup create indexes for user collection.
func (db *MongoBaseDb) Setup(config interface{}) error {
	// Copy mongo session (thread safe) and close after function
	conn := db.Connection.Copy()
	defer conn.Close()

	idx, ok := config.([]mgo.Index)
	if !ok {
		return errors.New("lxDb.mongoDb.Setup config interface is not []mgo.index")
	}

	// Ensure indexes
	col := conn.DB(db.Name).C(db.Collection)

	for _, i := range idx {
		if err := col.EnsureIndex(i); err != nil {
			return err
		}
	}

	return nil
}

// Create, create new entity in collection
func (db *MongoBaseDb) Create(data interface{}) error {
	// Copy mongo session (thread safe) and close after function
	conn := db.Connection.Copy()
	defer conn.Close()

	// Insert data
	return conn.DB(db.Name).C(db.Collection).Insert(data)
}

// GetAll, get all entities by query in collection
func (db *MongoBaseDb) GetAll(query interface{}, result interface{}, opts *Options) (int, error) {
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
func (db *MongoBaseDb) GetCount(query interface{}) (int, error) {
	// Copy mongo session (thread safe) and close after function
	conn := db.Connection.Copy()
	defer conn.Close()

	// Find all with query in collection
	return conn.DB(db.Name).C(db.Collection).Find(query).Count()
}

// GetOne, get one entity by query in collection
func (db *MongoBaseDb) GetOne(query interface{}, result interface{}) error {
	// Copy mongo session (thread safe) and close after function
	conn := db.Connection.Copy()
	defer conn.Close()

	// Find one with query in collection
	return conn.DB(db.Name).C(db.Collection).Find(query).One(result)
}

// Update, update one matched entity by query in collection
func (db *MongoBaseDb) Update(query interface{}, data interface{}) error {
	// Copy mongo session (thread safe) and close after function
	conn := db.Connection.Copy()
	defer conn.Close()

	// Update one with query in collection
	return conn.DB(db.Name).C(db.Collection).Update(query, data)
}

// UpdateAll, update all matched entities by query in collection
func (db *MongoBaseDb) UpdateAll(query interface{}, data interface{}) (ChangeInfo, error) {
	// Copy mongo session (thread safe) and close after function
	conn := db.Connection.Copy()
	defer conn.Close()

	// Update all with query in collection
	info, err := conn.DB(db.Name).C(db.Collection).UpdateAll(query, data)
	changeInfo := ChangeInfo{
		Updated: info.Updated,
		Removed: info.Removed,
		Matched: info.Matched,
	}

	return changeInfo, err
}

// Delete, delete one matched entity by query in collection
func (db *MongoBaseDb) Delete(query interface{}) error {
	// Copy mongo session (thread safe) and close after function
	conn := db.Connection.Copy()
	defer conn.Close()

	// Delete one with query in collection
	return conn.DB(db.Name).C(db.Collection).Remove(query)
}

// DeleteAll, delete all matched entities by query in collection
func (db *MongoBaseDb) DeleteAll(query interface{}) (ChangeInfo, error) {
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
