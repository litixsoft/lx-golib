package lxDb

import (
	"github.com/globalsign/mgo"
)

// Db struct for mongodb
type mongoDb struct {
	Conn *mgo.Session
	Name string
	Collection string
}

func NewMongoDb(connection *mgo.Session, dbName, collection string) *mongoDb {
	return &mongoDb{
		Conn: connection,
		Name: dbName,
		Collection: collection,
	}
}

// Setup create indexes for user collection.
func (db *mongoDb) Setup(indexes []mgo.Index) error {
	// Copy mongo session (thread safe) and close after function
	conn := db.Conn.Copy()
	defer conn.Close()

	// Ensure indexes
	col := conn.DB(db.Name).C(db.Collection)

	for _, i := range indexes {
		if err := col.EnsureIndex(i); err != nil {
			return err
		}
	}

	return nil
}