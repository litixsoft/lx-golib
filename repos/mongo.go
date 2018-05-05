package lxBaseRepo

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// Repository for mongodb
type MongoDb struct {
	Conn *mgo.Session
	Db   string
	Coll string
}

// Create, create new entity in collection
func (rm *MongoDb) Create(data interface{}) error {
	// Copy mongo session (thread safe) and close after function
	conn := rm.Conn.Copy()
	defer conn.Close()

	// Insert data
	return conn.DB(rm.Db).C(rm.Coll).Insert(data)
}

// GetAll, get all entities by query in collection
func (rm *MongoDb) GetAll(query interface{}, result interface{}, opts *Options) (int, error) {
	// Copy mongo session (thread safe) and close after function
	conn := rm.Conn.Copy()
	defer conn.Close()

	// Set default count
	n := 0

	// Check if activate counter in options
	if opts.Count {
		var err error
		n, err = conn.DB(rm.Db).C(rm.Coll).Find(query).Count()
		if err != nil {
			return n, err
		}
	}

	// Find all with query in collection
	return n, conn.DB(rm.Db).C(rm.Coll).Find(query).Skip(opts.Skip).Limit(opts.Limit).All(result)
}

// GetCount, get count of entities by query in collection
func (rm *MongoDb) GetCount(query interface{}) (int, error) {
	// Copy mongo session (thread safe) and close after function
	conn := rm.Conn.Copy()
	defer conn.Close()

	// Find all with query in collection
	return conn.DB(rm.Db).C(rm.Coll).Find(query).Count()
}

// GetOne, get one entity by query in collection
func (rm *MongoDb) GetOne(query interface{}, result interface{}) error {
	// Copy mongo session (thread safe) and close after function
	conn := rm.Conn.Copy()
	defer conn.Close()

	// Find one with query in collection
	return conn.DB(rm.Db).C(rm.Coll).Find(query).One(result)
}

// Update, update one matched entity by query in collection
func (rm *MongoDb) Update(query interface{}, data interface{}) error {
	// Copy mongo session (thread safe) and close after function
	conn := rm.Conn.Copy()
	defer conn.Close()

	// Update one with query in collection
	return conn.DB(rm.Db).C(rm.Coll).Update(query, bson.M{"$set": data})
}

// UpdateAll, update all matched entities by query in collection
func (rm *MongoDb) UpdateAll(query interface{}, data interface{}) (ChangeInfo, error) {
	// Copy mongo session (thread safe) and close after function
	conn := rm.Conn.Copy()
	defer conn.Close()

	// Update all with query in collection
	info, err := conn.DB(rm.Db).C(rm.Coll).UpdateAll(query, bson.M{"$set": data})
	changeInfo := ChangeInfo{
		Updated: info.Updated,
		Removed: info.Removed,
		Matched: info.Matched,
	}

	return changeInfo, err
}

// Delete, delete one matched entity by query in collection
func (rm *MongoDb) Delete(query interface{}) error {
	// Copy mongo session (thread safe) and close after function
	conn := rm.Conn.Copy()
	defer conn.Close()

	// Delete one with query in collection
	return conn.DB(rm.Db).C(rm.Coll).Remove(query)
}

// DeleteAll, delete all matched entities by query in collection
func (rm *MongoDb) DeleteAll(query interface{}) (ChangeInfo, error) {
	// Copy mongo session (thread safe) and close after function
	conn := rm.Conn.Copy()
	defer conn.Close()

	// Remove all with query in collection
	info, err := conn.DB(rm.Db).C(rm.Coll).RemoveAll(query)
	changeInfo := ChangeInfo{
		Updated: info.Updated,
		Removed: info.Removed,
		Matched: info.Matched,
	}

	return changeInfo, err
}
