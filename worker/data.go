package worker

import "time"

type StopTime struct {
	RouteID       string    `json:"route_id"`
	DirectionID   int       `json:"direction_id"`
	TripID        string    `json:"trip_id"`
	ArrivalTime   time.Time `json:"arrival_time"`
	DepartureTime time.Time `json:"departure_time"`
	StopID        string    `json:"stop_id"`
	StopSequence  int       `json:"stop_sequence"`
}
