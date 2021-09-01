package visits

import (
	"fmt"
	"log"

	"github.com/ct-fiuba/ct-contagion-updater/pkg/utils/mongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Visit struct {
	ScanCode           primitive.ObjectID `bson:"scanCode"`
	UserGeneratedCode  string             `bson:"userGeneratedCode"`
	EntranceTimestamp  primitive.DateTime `bson:"entranceTimestamp"`
	ExitTimestamp      primitive.DateTime `bson:"exitTimestamp"`
	Vaccinated         int                `bson:"vaccinated"`
	VaccineReceived    string             `bson:"vaccineReceived,omitempty"`
	VaccinatedDate     primitive.DateTime `bson:"vaccinatedDate,omitempty"`
	CovidRecovered     bool               `bson:"covidRecovered"`
	CovidRecoveredDate primitive.DateTime `bson:"covidRecoveredDate,omitempty"`
	DetectedTimestamp  primitive.DateTime `bson:"detectedTimestamp,omitempty"`
}

type VisitsCollection struct {
	Collection *mongo.Collection
	Database   *mongodb.DB
}

func New(db *mongodb.DB) (*VisitsCollection, error) {
	visits := &VisitsCollection{
		Collection: nil,
		Database:   db,
	}

	visits.Collection = db.Database.Collection("visits")

	return visits, nil
}

func (visits *VisitsCollection) All() ([]Visit, error) {
	var documents []Visit
	cursor, err := visits.Collection.Find(visits.Database.Context, bson.D{})
	if err != nil {
		log.Printf("Error finding")
		fmt.Println(err)
		return nil, err
	}
	err = cursor.All(visits.Database.Context, &documents)
	if err != nil {
		log.Printf("Error in All")
		fmt.Println(err)
		return nil, err
	}
	return documents, nil
}
