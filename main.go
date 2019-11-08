package main

import (
	"encoding/json"
	"fmt"
	"github.com/benjaminkomen/SchipholFlightAnalyser/outbound/schipholClient"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	http.HandleFunc("/", flightHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf("127.0.0.1:%s", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func flightHandler(w http.ResponseWriter, r *http.Request) {
	var client = schipholClient.New()

	flights, err := client.GetFlights(r.Context(), "D")
	if err != nil {
		log.Printf("Error obtaining flights")
		w.WriteHeader(500)
		return
	}

	err = json.NewEncoder(w).Encode(flights)
}
