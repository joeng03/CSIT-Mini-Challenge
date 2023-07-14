package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"log"
	"time"
	"mighty-saver-rabbit/constants"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
)

type Flight struct {
	Airline string
	Price int
	Srccity string
	Srccountry string 
	Destcity string 
	Destcountry string 
}
/*
{
	"_id": {
	  "$oid": "648095079b6d8b50581b727b"
	},
	"airline": "LH",
	"airlineid": 3320,
	"srcairport": "SIN",
	"srcairportid": 3316,
	"destairport": "FRA",
	"destairportid": 340,
	"codeshare": "",
	"stop": 0,
	"eq": "388",
	"airlinename": "Lufthansa",
	"srcairportname": "Singapore Changi Airport",
	"srccity": "Singapore",
	"srccountry": "Singapore",
	"destairportname": "Frankfurt am Main Airport",
	"destcity": "Frankfurt",
	"destcountry": "Germany",
	"price": 2432,
	"date": {
	  "$date": "2023-12-10T00:00:00.000Z"
	}
  }
  {
    "City": "Frankfurt",
    "Departure Date": "2023-12-10",
    "Departure Airline": "US Airways",
    "Departure Price": 1766,
    "Return Date": "2023-12-16",
    "Return Airline": "US Airways",
    "Return Price": 716
  }
*/


func GetFlights(w http.ResponseWriter, r *http.Request, flightsCollection *mongo.Collection) {

	params := r.URL.Query()
	departureDate, err := time.Parse("2006-01-02", params.Get("departureDate"))
	source := constants.SOURCE_CITY
	returnDate, err := time.Parse("2006-01-02", params.Get("returnDate"))
	destination := strings.Title(strings.ToLower(params.Get("destination")))
	

	
	options := options.FindOneOptions{
		Sort:  bson.M{"price": 1},
		Limit: 1,
	}
	cheapestDepartingFlight := flightsCollection.FindOne(context.Background(), bson.M{"date": departureDate, "srccity": source, "destcity": destination}, &options)
	returningFlights, err := flightsCollection.Find(context.Background(), bson.M{"date": returnDate, "srccity": destination, "destcity": source})

	if err := result.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	defer returningFlights.Close(context.Background())


	// var departingFlightsSlice []Flight
	// for departingFlights.Next(context.Background()) {
	// 	var flight Flight
	// 	err := departingFlights.Decode(&flight)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	departingFlightsSlice = append(departingFlightsSlice, flight)
	// }

	var returningFlightsSlice []Flight
	for returningFlights.Next(context.Background()) {
		var flight Flight
		err := returningFlights.Decode(&flight)
		if err != nil {
			log.Fatal(err)
		}
		returningFlightsSlice = append(returningFlightsSlice, flight)
	}

	flightsJSON, err := json.Marshal(returningFlightsSlice)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(flightsJSON)
}
