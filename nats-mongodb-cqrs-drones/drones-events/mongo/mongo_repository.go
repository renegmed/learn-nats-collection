package mongo

import (
	"context"

	"cqrs-drones/events/common"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

// NewEventRollupRepository creates a new mongoDB event rollup repository with the supplied collections.
func NewEventRollupRepository(db *mongo.Database) (repo *EventRollupRepository) {
	repo = &EventRollupRepository{
		PositionsCollection: db.Collection("positions"),
		AlertsCollection:    db.Collection("alerts"),
		TelemetryCollection: db.Collection("telemetry"),
	}
	return
}

// UpdateLastTelemetryEvent updates the most recent telemetry event, or creates a new one if one has never been received.
func (repo *EventRollupRepository) UpdateLastTelemetryEvent(telemetryEvent common.TelemetryUpdatedEvent) (err error) {
	var recordID bson.ObjectId
	var newRecord *mongoTelemetryRecord
	foundEvent, err := repo.getTelemetryRecord(context.TODO(), telemetryEvent.DroneID)
	if err != nil {
		recordID = bson.NewObjectId()
		newRecord = convertTelemetryEventToRecord(telemetryEvent, recordID)
		_, err = repo.TelemetryCollection.InsertOne(context.TODO(), newRecord)
	} else {
		recordID = foundEvent.RecordID
		newRecord = convertTelemetryEventToRecord(telemetryEvent, recordID)
		_, err = repo.TelemetryCollection.ReplaceOne(context.TODO(), bson.D{{"_id", recordID}}, newRecord)
	}
	return
}

// UpdateLastAlertEvent updates the most recent alert event, or creates a new one if one has never been received.
func (repo *EventRollupRepository) UpdateLastAlertEvent(alertEvent common.AlertSignalledEvent) (err error) {
	var recordID bson.ObjectId
	var newRecord *mongoAlertRecord
	foundEvent, err := repo.getAlertRecord(context.TODO(), alertEvent.DroneID)
	if err != nil {
		recordID = bson.NewObjectId()
		newRecord = convertAlertEventToRecord(alertEvent, recordID)
		_, err = repo.AlertsCollection.InsertOne(context.TODO(), newRecord)
	} else {
		recordID = foundEvent.RecordID
		newRecord = convertAlertEventToRecord(alertEvent, recordID)
		_, err = repo.AlertsCollection.ReplaceOne(context.TODO(), bson.D{{"_id", recordID}}, newRecord)
	}

	return
}

// UpdateLastPositionEvent updates the last position event, or creates a new one if one has never been received.
func (repo *EventRollupRepository) UpdateLastPositionEvent(positionEvent common.PositionChangedEvent) (err error) {
	var recordID bson.ObjectId
	var newRecord *mongoPositionRecord
	foundEvent, err := repo.getPositionRecord(context.TODO(), positionEvent.DroneID)
	if err != nil {
		recordID = bson.NewObjectId()
		newRecord = convertPositionEventToRecord(positionEvent, recordID)
		_, err = repo.PositionsCollection.InsertOne(context.TODO(), newRecord)
	} else {
		recordID = foundEvent.RecordID
		newRecord = convertPositionEventToRecord(positionEvent, recordID)
		_, err = repo.PositionsCollection.ReplaceOne(context.TODO(), bson.D{{"_id", recordID}}, newRecord)
	}

	return
}

// GetTelemetryEvent retrieves the most recent telemetry event for a given drone.
func (repo *EventRollupRepository) GetTelemetryEvent(droneID string) (event common.TelemetryUpdatedEvent, err error) {
	record, err := repo.getTelemetryRecord(context.TODO(), droneID)
	if err == nil {
		event = convertTelemetryRecordToEvent(record)
	}

	return
}

// GetAlertEvent retrieves the most recent alert event for a given drone.
func (repo *EventRollupRepository) GetAlertEvent(droneID string) (event common.AlertSignalledEvent, err error) {
	record, err := repo.getAlertRecord(context.TODO(), droneID)
	if err == nil {
		event = convertAlertRecordToEvent(record)
	}
	return
}

// GetPositionEvent retrieves the most recent position event for a given drone.
func (repo *EventRollupRepository) GetPositionEvent(droneID string) (event common.PositionChangedEvent, err error) {
	record, err := repo.getPositionRecord(context.TODO(), droneID)
	if err == nil {
		event = convertPositionRecordToEvent(record)
	}
	return
}

func (repo *EventRollupRepository) getTelemetryRecord(ctx context.Context, droneID string) (record mongoTelemetryRecord, err error) {
	query := bson.M{"drone_id": droneID}
	err = repo.TelemetryCollection.FindOne(ctx, query).Decode(&record)
	return

}

func (repo *EventRollupRepository) getAlertRecord(ctx context.Context, droneID string) (record mongoAlertRecord, err error) {
	query := bson.M{"drone_id": droneID}
	err = repo.AlertsCollection.FindOne(ctx, query).Decode(&record)
	return

}

func (repo *EventRollupRepository) getPositionRecord(ctx context.Context, droneID string) (record mongoPositionRecord, err error) {
	query := bson.M{"drone_id": droneID}
	err = repo.PositionsCollection.FindOne(ctx, query).Decode(&record)
	return

}
