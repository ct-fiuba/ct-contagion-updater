package compromisedCodes

import (
	"fmt"

	"github.com/ct-fiuba/ct-contagion-updater/pkg/utils/mongodb"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CompromisedCode struct {
	SpaceId           primitive.ObjectID `bson:"spaceId"`
	UserGeneratedCode string             `bson:"userGeneratedCode"`
	DateDetected      primitive.DateTime `bson:"dateDetected"`
	Risk              int                `bson:"risk"`
}

type CompromisedCodesCollection struct {
	Collection *mongo.Collection
	Database   *mongodb.DB
}

func New(db *mongodb.DB) (*CompromisedCodesCollection, error) {
	compromised := &CompromisedCodesCollection{
		Collection: nil,
		Database:   db,
	}

	compromised.Collection = db.Database.Collection("compromisedcodes")

	return compromised, nil
}

func (compromised *CompromisedCodesCollection) Insert(cc CompromisedCode) error {
	res, err := compromised.Collection.InsertOne(compromised.Database.Context, cc)
	if err != nil {
		fmt.Errorf("%w\n", err)
		return err
	}

	fmt.Printf("Inserted new compromised code with ID: %v\n", res.InsertedID)
	return nil
}
