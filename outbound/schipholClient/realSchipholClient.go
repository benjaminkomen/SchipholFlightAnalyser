package schipholClient

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/benjaminkomen/SchipholFlightAnalyser/outbound/config"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const timeout = 10 * time.Second

type realSchipholClient struct{}

func (rsc *realSchipholClient) GetFlights(c context.Context, flightDirection string, scheduleDate string) ([]Flight, error) {

	allResultingFlights := []Flight{}

	configModel, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	// first call
	resultingFlights, nextLink, err := obtainFlightsForUrl(c, configModel.AppId, configModel.AppKey, allResultingFlights, fmt.Sprintf("https://api.schiphol.nl/public-flights/flights?flightDirection=%s&scheduleDate=%s", flightDirection, scheduleDate))
	if err != nil {
		return nil, err
	}

	if nextLink != "" {
		// the rest of the calls
		for nextLink != "" {
			resultingFlights, nextLink, err = obtainFlightsForUrl(c, configModel.AppId, configModel.AppKey, resultingFlights, nextLink)
			if err != nil {
				return nil, err
			}
		}
	}

	allResultingFlights = resultingFlights

	return allResultingFlights, nil
}

func obtainFlightsForUrl(c context.Context, appId string, appKey string, inputFlights []Flight, url string) ([]Flight, string, error) {
	ctxWithAlternativeTimeout, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	client := http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error composing new request")
		return nil, "", err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("ResourceVersion", "v4")
	req.Header.Set("app_id", appId)
	req.Header.Set("app_key", appKey)

	resp, err := client.Do(req.WithContext(ctxWithAlternativeTimeout))
	if err != nil {
		log.Printf("Sending GET request to api.schiphol.com")
		return nil, "", err
	}

	if resp.StatusCode > 299 {
		return nil, "", err
	}

	var response schipholSuccessResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Printf("Error decoding response from Schiphol")
		return nil, "", err
	}

	responseLinks := parseResponseLinks(resp.Header.Get("link"))
	log.Printf("Response link is %+v", responseLinks)

	outputFlights := append(inputFlights, response.Flights...)

	defer resp.Body.Close()

	return outputFlights, responseLinks["next"], nil
}

func parseResponseLinks(link string) map[string]string {

	responseLinks := make(map[string]string, 3)

	splitLinkByComma := strings.Split(link, ",")
	r := regexp.MustCompile(`<(.*?)>;\srel="(.*?)"`)

	for _, splitLinkByCommaPart := range splitLinkByComma {
		regexResult := r.FindStringSubmatch(splitLinkByCommaPart)
		responseLinks[regexResult[2]] = regexResult[1] // format: rel = link
	}

	return responseLinks
}
