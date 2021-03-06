package spaces

import (
	"fmt"
	"log"

	"github.com/ct-fiuba/ct-contagion-updater/pkg/utils/mongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Space struct {
	Name                   string             `bson:"name"`
	M2                     int                `bson:"m2"`
	EstimatedVisitDuration int                `bson:"estimatedVisitDuration"`
	HasExit                bool               `bson:"hasExit"`
	OpenSpace              bool               `bson:"openSpace"`
	EstablishmentId        primitive.ObjectID `bson:"establishmentId"`
	N95Mandatory           bool               `bson:"n95Mandatory"`
	Enabled                bool               `bson:"enabled"`
}

type SpacesCollection struct {
	Collection *mongo.Collection
	Database   *mongodb.DB
}

func New(db *mongodb.DB) (*SpacesCollection, error) {
	spaces := &SpacesCollection{
		Collection: nil,
		Database:   db,
	}

	spaces.Collection = db.Database.Collection("spaces")

	return spaces, nil
}

func (spaces *SpacesCollection) Find(id string) (*Space, error) {
	var space Space

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Error while parsing space's ID (%s)", id)
		fmt.Println(err)
		return nil, err
	}

	err = spaces.Collection.FindOne(
		spaces.Database.Context,
		bson.D{{"_id", objectId}},
	).Decode(&space)

	if err != nil {
		log.Printf("Error getting a Space with ID=%s", id)
		fmt.Println(err)
		return nil, err
	}

	return &space, nil
}
