package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type DB struct {
	Client   *mongo.Client
	Database *mongo.Database
	Context  context.Context
	Cancel   context.CancelFunc
}

func New(dbUri string, dbName string) (*DB, error) {
	db := &DB{
		Client:   nil,
		Database: nil,
		Context:  nil,
		Cancel:   nil,
	}

	var err error

	log.Printf("Trying to Connect to the MongoDB %s, at address %s", dbName, dbUri)

	db.Context, db.Cancel = context.WithCancel(context.Background())
	db.Client, err = mongo.Connect(db.Context, options.Client().ApplyURI(dbUri))
	if err != nil {
		log.Printf("Error connecting to DB")
		return nil, err
	}

	log.Printf("Connected to the DB!")

	db.Database = db.Client.Database(dbName)

	return db, nil
}

func (db *DB) Shutdown() error {
	defer db.Client.Disconnect(db.Context)
	defer db.Cancel()
	return nil
}
