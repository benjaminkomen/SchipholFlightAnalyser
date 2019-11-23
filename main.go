package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/benjaminkomen/SchipholFlightAnalyser/outbound/schipholClient"
	"github.com/go-echarts/go-echarts/charts"
	"log"
	"net/http"
	"os"
	"sort"
	"time"
)

func main() {

	http.HandleFunc("/api/flights", flightHandler)
	http.HandleFunc("/", flightChartHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf("127.0.0.1:%s", port),
		WriteTimeout: 600 * time.Second,
		ReadTimeout:  600 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func flightChartHandler(w http.ResponseWriter, r *http.Request) {
	data, err := getFlightChartData(r.Context())
	log.Printf("Obtained chart data")
	sortedData := sortMapByKey(data)
	log.Printf("Sorted chart data")
	if err != nil {
		log.Printf("error getting flight chart data: %s", err)
		w.WriteHeader(500)
		return
	}

	var nameItems []int
	var departingFlights []int
	var arrivingFlights []int

	for _, dataEntry := range sortedData {
		nameItems = append(nameItems, dataEntry.Key)
		departingFlights = append(departingFlights, dataEntry.Value.departCount)
		arrivingFlights = append(arrivingFlights, dataEntry.Value.arrivingCount)
	}

	log.Printf("Start drawing chart")
	bar := charts.NewBar()
	bar.SetGlobalOptions(charts.TitleOpts{Title: "Flights per hour of the day"})
	bar.AddXAxis(nameItems)
	bar.AddYAxis("departing", departingFlights, charts.BarOpts{}).
		AddYAxis("arriving", arrivingFlights)

	f, err := os.Create("html/bar.html")
	if err != nil {
		log.Println(err)
	}
	log.Printf("Wrote chart data to html")

	err = bar.Render(w, f)
	if err != nil {
		log.Printf("error rendering chart: %s", err)
		w.WriteHeader(500)
		return
	}
}

type Pair struct {
	Key   int
	Value FlightData
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Key < p[j].Key }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func sortMapByKey(data map[int]FlightData) PairList {
	pairList := make(PairList, len(data))
	i := 0
	for k, v := range data {
		pairList[i] = Pair{k, v}
		i++
	}
	sort.Sort(pairList)
	return pairList
}

func flightHandler(w http.ResponseWriter, r *http.Request) {
	var client = schipholClient.New()

	flights, err := client.GetFlights(r.Context(), "D", "2019-11-22")
	if err != nil {
		log.Printf("error obtaining flights: %s", err)
		w.WriteHeader(500)
		return
	}

	err = json.NewEncoder(w).Encode(flights)
	if err != nil {
		log.Printf("error encoding flights to json: %s", err)
		w.WriteHeader(500)
		return
	}
}

type FlightData struct {
	departCount   int
	arrivingCount int
}

func (fd *FlightData) append(departCount int, arrivingCount int) {

	if departCount != 0 {
		fd.departCount += 1
	}

	if arrivingCount != 0 {
		fd.arrivingCount += 1
	}
}

func getFlightChartData(c context.Context) (map[int]FlightData, error) {
	results := make(map[int]FlightData, 24)
	var client = schipholClient.New()

	departingFlights, err := client.GetFlights(c, "D", "2019-11-22")
	if err != nil {
		log.Printf("error obtaining departing flights: %s", err)
		return nil, err
	}

	for _, flight := range departingFlights {

		if flight.ActualOffBlockTime.IsZero() {
			continue // we skip results with no reliable data
		}

		hour := flight.ActualOffBlockTime.Hour()
		data := results[hour]
		data.append(hour, 0)
		results[hour] = data
	}

	arrivingFlights, err := client.GetFlights(c, "A", "2019-11-22")
	if err != nil {
		log.Printf("error obtaining arriving flights: %s", err)
		return nil, err
	}

	for _, flight := range arrivingFlights {

		if flight.ActualLandingTime.IsZero() {
			continue // we skip results with no reliable data
		}

		hour := flight.ActualLandingTime.Hour()
		data := results[hour]
		data.append(0, hour)
		results[hour] = data
	}
	return results, nil
}
