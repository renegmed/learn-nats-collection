package mongo

import (
	"context"
	"log"
	"testing"
	"time"

	"cqrs-drones/events/common"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const TEST_PORT = "27018"

func getTestDbase() (*mongo.Database, error) {

	clientOptions := options.Client().ApplyURI("mongodb://localhost:" + TEST_PORT)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}
	//Set up a context required by mongo.Connect
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//To close the connection at the end
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal("Couldn't connect to the database", err)
		return nil, err
	} else {
		log.Println("Connected!")
	}
	db := client.Database("go_mongo")
	return db, nil
}

func TestCreate(t *testing.T) {
	dbase, err := getTestDbase()
	if err != nil {
		t.Fail()
	}

	event := common.PositionChangedEvent{}
	event.DroneID = "123456"
	event.Latitude = 123.2334425
	event.Longitude = 13.4443565
	event.Altitude = 1250.33
	event.CurrentSpeed = 654.33
	event.HeadingCardinal = 2
	event.ReceivedOn = 12321332 //time.Now()

	repo := NewEventRollupRepository(dbase)
	err = repo.UpdateLastPositionEvent(event)
	if err != nil {
		t.Fail()
	}

	posChangedEvent, err := repo.GetPositionEvent(context.TODO(), event.DroneID)
	if err != nil {
		t.Fail()
	}

	if posChangedEvent.DroneID != event.DroneID {
		t.Logf("DroneID should be %s not %s", event.DroneID, posChangedEvent)
		t.Fail()
	}

}
