package schipholClient

import "context"

type mockSchipholClient struct{}

func (msc *mockSchipholClient) GetFlights(c context.Context, flightDirection string) ([]Flight, error) {
	return []Flight{}, nil
}
