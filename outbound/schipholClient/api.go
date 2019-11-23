package schipholClient

import (
	"context"
	"time"
)

type SchipholClient interface {
	GetFlights(c context.Context, flightDirection string, scheduleDate string) ([]Flight, error)
}

func New() SchipholClient {
	return &realSchipholClient{}
}

type schipholSuccessResponse struct {
	Flights []Flight `json:"flights"`
}

type Flight struct {
	LastUpdatedAt      time.Time `json:"lastUpdatedAt"`
	FlightName         string    `json:"flightName"`
	ActualLandingTime  time.Time `json:"actualLandingTime,omitempty"`  // in case of Arriving flights
	ActualOffBlockTime time.Time `json:"actualOffBlockTime,omitempty"` // in case of Departing flights
}
