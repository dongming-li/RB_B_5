package database

import (
	"log"
	"time"

	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/model"
	mgo "gopkg.in/mgo.v2"
)

// MongoCollection represents a mongogb implementation of [model.Collection]
type MongoCollection struct {
	Collection *mgo.Collection
}

// MongoQueryResult represents a mongogb implementation of [model.QueryResult]
type MongoQueryResult struct {
	query *mgo.Query
}

// MongoChangeInfo holds details about the outcome of an mGo update operation.
type MongoChangeInfo struct {
	changeInfo *mgo.ChangeInfo
}

// All returns all the values from a query. Mongodb implementation from [model.aLL]
func (rq MongoQueryResult) All(result interface{}) error {
	err := rq.query.All(result) //TODO: mGo-switch
	if err != nil {
		log.Println("Database error: ", err)
	}
	return err
}

// One unmarshals the first contained doc into [result] and returns an error if any
func (rq MongoQueryResult) One(result interface{}) error {
	err := rq.query.One(result) //TODO: mGo-switch
	//mGo-switch  - Currently, mgo unmarshals the data even when there is an error which sometimes leads to a partial or no unmarshal
	// We could manually always return an empty struct if there's an error
	// like
	// if err != nil {
	// 	result = struct{}{}
	// 	return err
	// }
	if err != nil {
		log.Println("Database error: ", err)
	}
	return err
}

// Apply performs an update operation on a database find op
func (rq MongoQueryResult) Apply(change model.Change, result interface{}) (info model.ChangeInfo, err error) {
	chInfo, err := rq.query.Apply(mgo.Change{
		Update:    change.Update,
		ReturnNew: change.ReturnNew,
	}, result)
	if err != nil {
		return nil, err
	}
	return MongoChangeInfo{changeInfo: chInfo}, nil
}

// GetMongoCollection returns a mongo collection with the [name]
func GetMongoCollection(mC *mgo.Collection) MongoCollection {
	c := MongoCollection{mC}
	if c.Collection.Name == "game" {
		err := c.EnsureIndex([]string{"endTime"}, time.Duration(1)*time.Second)
		if err != nil {
			panic(err)
		}
	}

	return c
}

// Find queries the database with the  [query]. Mongodb implementation from [model.Find]
func (c MongoCollection) Find(query interface{}) model.QueryResult {
	return MongoQueryResult{c.Collection.Find(query)}
}

// Insert inserts one or more documents into the collection
func (c MongoCollection) Insert(docs ...interface{}) error {
	return c.Collection.Insert(docs...)
}

// Update updates one or more documents into the collection
func (c MongoCollection) Update(selector interface{}, update interface{}) error {
	return c.Collection.Update(selector, update)
}

// Remove removes one document in the collection
func (c MongoCollection) Remove(selector interface{}) error {
	return c.Collection.Remove(selector)
}

// EnsureIndex will create an index for a document
func (c MongoCollection) EnsureIndex(key []string, time time.Duration) error {
	index := mgo.Index{
		Key:         key,
		ExpireAfter: time,
		Background:  true,
	}
	return c.Collection.EnsureIndex(index)
}

// FindAndModify finds a document and performs an [update] operation on the document
func (c MongoCollection) FindAndModify(selector interface{}, update interface{}, result interface{}, returnNew bool) (model.ChangeInfo, error) {

	chInfo, err := c.Collection.Find(selector).Apply(mgo.Change{
		Update:    update,
		ReturnNew: returnNew,
	}, result)

	if err != nil {
		return nil, err
	}
	return MongoChangeInfo{changeInfo: chInfo}, nil
}

//Removed returns the number of documents removed
func (c MongoChangeInfo) Removed() int {
	return c.changeInfo.Removed
}

//Updated returns the number of documents updated
func (c MongoChangeInfo) Updated() int {
	return c.changeInfo.Updated
}
