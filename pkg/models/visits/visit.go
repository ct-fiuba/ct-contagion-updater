package visits

import (
	"fmt"
	"log"
	"time"

	"github.com/ct-fiuba/ct-contagion-updater/pkg/utils/mongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	TIME_WINDOW_DAYS = 21
)

type Visit struct {
	SpaceId              primitive.ObjectID `bson:"spaceId"`
	UserGeneratedCode    string             `bson:"userGeneratedCode"`
	EntranceTimestamp    primitive.DateTime `bson:"entranceTimestamp"`
	ExitTimestamp        primitive.DateTime `bson:"exitTimestamp"`
	Vaccinated           int                `bson:"vaccinated"`
	VaccineReceived      string             `bson:"vaccineReceived,omitempty"`
	VaccinatedDate       primitive.DateTime `bson:"vaccinatedDate,omitempty"`
	IllnessRecovered     bool               `bson:"illnessRecovered"`
	IllnessRecoveredDate primitive.DateTime `bson:"illnessRecoveredDate,omitempty"`
	DetectedTimestamp    primitive.DateTime `bson:"detectedTimestamp,omitempty"`
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

func (visits *VisitsCollection) FindInSpace(spaceId string) ([]Visit, error) {
	var documents []Visit

	objectId, err := primitive.ObjectIDFromHex(spaceId)
	if err != nil {
		log.Printf("Error while parsing space's ID (%s)", spaceId)
		fmt.Println(err)
		return nil, err
	}

	cursor, err := visits.Collection.Find(
		visits.Database.Context,
		bson.M{
			"spaceId":           objectId,
			"entranceTimestamp": bson.M{"$gte": primitive.NewDateTimeFromTime(time.Now().AddDate(0, 0, -TIME_WINDOW_DAYS))},
		},
	)
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

func (visits *VisitsCollection) FindByGeneratedCode(userGeneratedCode string) (*Visit, error) {
	var visit Visit

	err := visits.Collection.FindOne(
		visits.Database.Context,
		bson.D{{"userGeneratedCode", userGeneratedCode}},
	).Decode(&visit)

	if err != nil {
		log.Printf("Error getting a Visit with userGeneratedCode=%s", userGeneratedCode)
		fmt.Println(err)
		return nil, err
	}

	return &visit, nil
}
