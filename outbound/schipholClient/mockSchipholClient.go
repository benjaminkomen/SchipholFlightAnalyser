package schipholClient

import "context"

type mockSchipholClient struct{}

func (msc *mockSchipholClient) GetFlights(c context.Context, flightDirection string, scheduleDate string) ([]Flight, error) {
	return []Flight{}, nil
}
