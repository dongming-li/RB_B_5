package database

import (
	"log"

	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/config"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/model"
	"github.com/google/uuid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Database represents the current Database connection
type Database struct {
	session *mgo.Session
	url     string
}

// DatastoreSession is a session connection to the Database
type DatastoreSession interface { //TODO: not yet used
	Close()
	Copy() DatastoreSession
	DB(name string)
}

// Datastore represents a store for the data (Database session)
type Datastore struct {
	session *mgo.Session
	name    string
}

// Close closes a datastore's session
func (store *Datastore) Close() {
	store.session.Close()
	store.session = nil
	log.Printf("Closed Datastore: %s\n", store.name)
}

// GetUsers returns a collection of users
func (store *Datastore) GetUsers() model.Collection {
	log.Println("From Datastore: Getting users...")
	return GetMongoCollection(store.session.DB(config.DBNameMain).C("user"))
}

// GetCollection returns a collection with the given [name] from the store
func (store *Datastore) GetCollection(name string) model.Collection {
	log.Printf("From Datastore: Getting %s ...", name)
	return GetMongoCollection(store.session.DB(config.DBNameMain).C(name))
}

// UserIsPresent returns true if a user is present
func (store *Datastore) UserIsPresent(username string) (bool, error) {
	log.Println("From Datastore: Checking for user:", username, " ...")
	c, err := store.session.DB(config.DBNameMain).C("user").Find(bson.M{"username": username}).Count()
	return c != 0, err
}

// Database methods

// NewStoreFactory initializes a Database with the url without connecting it
func (d *Database) NewStoreFactory() model.StoreFactory {
	return func() model.Store {
		log.Println("Creating new store...")
		return &Datastore{
			name:    uuid.New().String(),
			session: d.session.Copy(),
		}
	}

}

// ConnectWithInfo creates a new Database session based on [ConnectWithInfo] and returns it
func ConnectWithInfo(url string) (*Database, error) {
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{url},
		Database: config.AuthDatabase,
		Username: config.AuthUserName,
		Password: config.AuthPassword,
	}
	log.Printf("Connecting to Database on: %s", url)
	sess, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		return nil, err
	}
	log.Println("Connected to Database")
	return &Database{session: sess}, nil
}

// Connect creates a new Database session and returns it
func Connect(url string) (*Database, error) {
	log.Printf("Connecting to Database on: %s", url)
	sess, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}
	log.Println("Connected to Database")
	return &Database{session: sess}, nil
}

// Close kills the current session and ends the Database connection
func (d *Database) Close() {
	if d.session != nil {
		d.session.Close()
	}
	d.session = nil
	log.Println("Database closed")
}
