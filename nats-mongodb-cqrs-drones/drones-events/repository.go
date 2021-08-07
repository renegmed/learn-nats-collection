package main

import (
	"context"

	"cqrs-drones/events/common"
)

type Repository interface {
	UpdateLastTelemetryEvent(telemetryEvent common.TelemetryUpdatedEvent) (err error)
	UpdateLastAlertEvent(alertEvent common.AlertSignalledEvent) (err error)
	UpdateLastPositionEvent(positionEvent common.PositionChangedEvent) (err error)
	GetTelemetryEvent(ctx context.Context, droneID string) (event common.TelemetryUpdatedEvent, err error)
	GetAlertEvent(ctx context.Context, droneID string) (event common.AlertSignalledEvent, err error)
	GetPositionEvent(ctx context.Context, droneID string) (event common.PositionChangedEvent, err error)
}
