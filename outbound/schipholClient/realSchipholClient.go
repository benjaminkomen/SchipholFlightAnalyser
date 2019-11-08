package schipholClient

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const timeout = 10 * time.Second

type realSchipholClient struct{}

func (rsc *realSchipholClient) GetFlights(c context.Context, flightDirection string) ([]Flight, error) {

	ctxWithAlternativeTimeout, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	configModel, err := GetConfig()

	client := http.Client{}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.schiphol.nl/public-flights/flights?flightDirection=%s", flightDirection), nil)
	if err != nil {
		log.Printf("Error composing new request")
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("ResourceVersion", "v4")
	req.Header.Set("app_id", configModel.AppId)
	req.Header.Set("app_key", configModel.AppKey)

	resp, err := client.Do(req.WithContext(ctxWithAlternativeTimeout))
	if err != nil {
		log.Printf("Sending GET request to api.schiphol.com")
		return nil, err
	}

	if resp.StatusCode <= 300 {
		var response schipholSuccessResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			log.Printf("Error decoding response from Schiphol")
			return nil, err
		}
		log.Printf("Successfully called Schiphol API and received: %+v", response)
		return response.Flights, nil
	}

	defer resp.Body.Close()

	return nil, nil
}
