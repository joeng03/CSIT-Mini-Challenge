package api

import (
	"context"
	"encoding/json"
	"log"
	"mighty-saver-rabbit/constants"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Flight struct {
	Airlinename string
	Price       int
	Srccity     string
	Srccountry  string
	Destcity    string
	Destcountry string
	Date        time.Time
}

type CheapestFlightResponseData struct {
	City             string
	DepartureDate    string `json:"Departure Date"`
	DepartureAirline string `json:"Departure Airline"`
	DeparturePrice   int    `json:"Departure Price"`
	ReturnDate       string `json:"Return Date"`
	ReturnAirline    string `json:"Return Airline"`
	ReturnPrice      int    `json:"Return Price"`
}

func GetFlights(w http.ResponseWriter, r *http.Request, flightsCollection *mongo.Collection) {

	params := r.URL.Query()
	departureDate, err := time.Parse("2006-01-02", params.Get("departureDate"))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}

	source := constants.SOURCE_CITY
	returnDate, err := time.Parse("2006-01-02", params.Get("returnDate"))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}

	destination := strings.Title(strings.ToLower(params.Get("destination")))

	filter := bson.M{"date": departureDate, "srccity": source, "destcity": destination}
	opts := options.Find().SetSort(bson.M{"price": 1}).SetLimit(1)
	cheapestDepartingFlights, err := flightsCollection.Find(context.Background(), filter, opts)

	returningFlights, err := flightsCollection.Find(context.Background(), bson.M{"date": returnDate, "srccity": destination, "destcity": source})

	var cheapestDepartingFlight Flight
	for cheapestDepartingFlights.Next(context.Background()) {
		if err := cheapestDepartingFlights.Decode(&cheapestDepartingFlight); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
		}

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
		}
	}

	var cheapestFlightsResponseSlice []CheapestFlightResponseData = make([]CheapestFlightResponseData, 0)

	for returningFlights.Next(context.Background()) {
		var returningFlight Flight
		if err := returningFlights.Decode(&returningFlight); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
		}
		cheapestFlightsResponseSlice = append(cheapestFlightsResponseSlice, CheapestFlightResponseData{
			City:             destination,
			DepartureDate:    departureDate.Format("2006-01-02"),
			DepartureAirline: cheapestDepartingFlight.Airlinename,
			DeparturePrice:   cheapestDepartingFlight.Price,
			ReturnDate:       returnDate.Format("2006-01-02"),
			ReturnAirline:    returningFlight.Airlinename,
			ReturnPrice:      returningFlight.Price,
		})
	}
	cheapestFlightsResponseJSON, _ := json.Marshal(cheapestFlightsResponseSlice)

	w.Header().Set("Content-Type", "application/json")
	w.Write(cheapestFlightsResponseJSON)
}
