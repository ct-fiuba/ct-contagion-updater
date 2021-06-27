package rules

import (
	"fmt"
	"log"

	"github.com/ct-fiuba/ct-contagion-updater/pkg/utils/mongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Rule struct {
	Index         int    `bson:"index"`
	ContagionRisk string `bson:"contagionRisk"`
	DurationCmp   string `bson:"durationCmp,omitempty"`
	DurationValue int    `bson:"durationValue,omitempty"`
	M2Cmp         string `bson:"m2Cmp,omitempty"`
	M2Value       int    `bson:"m2Value,omitempty"`
	SpaceValue    string `bson:"spaceValue,omitempty"`
}

type RulesCollection struct {
	Collection *mongo.Collection
	Database   *mongodb.DB
}

func New(db *mongodb.DB) (*RulesCollection, error) {
	rules := &RulesCollection{
		Collection: nil,
		Database:   db,
	}

	rules.Collection = db.Database.Collection("rules")

	return rules, nil
}

func (rules *RulesCollection) All() ([]Rule, error) {
	var documents []Rule
	cursor, err := rules.Collection.Find(rules.Database.Context, bson.D{})
	if err != nil {
		log.Printf("Error finding")
		fmt.Println(err)
		return nil, err
	}
	err = cursor.All(rules.Database.Context, &documents)
	if err != nil {
		log.Printf("Error in All")
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(documents)
	return documents, nil
}
