package rules

import (
	"fmt"
	"log"

	"github.com/ct-fiuba/ct-contagion-updater/pkg/utils/mongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Rule struct {
	Index                      int    `bson:"index"`
	ContagionRisk              int    `bson:"contagionRisk"`
	DurationCmp                string `bson:"durationCmp,omitempty"`
	DurationValue              int    `bson:"durationValue,omitempty"`
	M2Cmp                      string `bson:"m2Cmp,omitempty"`
	M2Value                    int    `bson:"m2Value,omitempty"`
	OpenSpace                  bool   `bson:"openSpace,omitempty"`
	N95Mandatory               bool   `bson:"n95Mandatory,omitempty"`
	Vaccinated                 int    `bson:"vaccinated,omitempty"`
	VaccinatedDaysAgoMin       int    `bson:"vaccinatedDaysAgoMin,omitempty"`
	VaccineReceived            string `bson:"vaccineReceived,omitempty"`
	IllnessRecovered           bool   `bson:"illnessRecovered,omitempty"`
	IllnessRecoveredDaysAgoMax int    `bson:"illnessRecoveredDaysAgoMax,omitempty"`
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

	findOptions := options.Find().SetSort(bson.D{{"index", 1}})
	cursor, err := rules.Collection.Find(rules.Database.Context, bson.D{}, findOptions)
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
	return documents, nil
}
