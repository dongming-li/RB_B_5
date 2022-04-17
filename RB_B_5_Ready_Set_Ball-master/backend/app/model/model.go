package model

import "time"

// StoreFactory is a closure function that returns a new instance of a store
type StoreFactory func() Store

// Store represents a database that holds data
type Store interface {
	// Ueful?
	// This is equivalent to GetCollection('user') but it's
	// safer since we'll ensure this method always points to
	// a collection for the users even if the name changes
	GetUsers() Collection //change to collection interface
	UserIsPresent(username string) (bool, error)
	// GetCollection returns a collection with the given [name]
	// NOTE: it creates a collection with that[name]
	// if one doesn't currently exist
	GetCollection(name string) Collection
	Close()
}

// Collection represents a row of data
type Collection interface {
	Find(query interface{}) QueryResult
	//TODO remove this method if it doesn't get used alot
	// It is supposed to remove the overhead fromFind(bson.M{"_id": id})
	// FindId(query interface{}) QueryResult

	// FindAndModify searches for a document using [selector]  and updates it with [update] if found
	// It returns the modified object if [ReturnNew] is true
	FindAndModify(selector interface{}, update interface{}, result interface{}, returnNew bool) (ChangeInfo, error)

	Insert(docs ...interface{}) error

	// Update performs an [update] on [selector] in the collection
	Update(selector interface{}, update interface{}) error

	//Remove will remove the document(s) that matches the [selector]
	Remove(selector interface{}) error

	//Will create an index key on a document
	EnsureIndex(key []string, time time.Duration) error
}

// QueryResult represnts the result from a query to [Collection]
type QueryResult interface {
	All(result interface{}) error
	One(result interface{}) error

	// Deprecated 0.0.1 Do not use
	Apply(change Change, result interface{}) (ChangeInfo, error)
}

// Change represents an update based on [mgo.Change]
// Deprecated 0.0.1 Do not use
type Change struct {
	Update    map[string]interface{}
	ReturnNew bool
}

// ChangeInfo holds details about the outcome of an update operation.
type ChangeInfo interface {
	Removed() int
	Updated() int
}

// entry represnts a possible input to the database that can be validated
type entry interface {
	OK() error
}

// Identity represents a unique element that can be used to access an entry from the [Store]
type Identity interface { //TODO: change name to Identifier
	Identify() map[string]interface{}
}
