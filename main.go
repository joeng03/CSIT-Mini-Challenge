package main

import (
	"context"
	"fmt"
	"net/http"
	"mighty-saver-rabbit/api"
	"mighty-saver-rabbit/constants"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	// MongoDB connection
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(constants.MONGO_URI))
	if err != nil {
		panic(err)
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to MongoDB")

	// Database collections
	db := client.Database("minichallenge")
	flightsCollection := db.Collection("flights")
	hotelsCollection := db.Collection("hotels")

	// Routers
	router := mux.NewRouter()

	getFlightsHandler := func(w http.ResponseWriter, r *http.Request) {
		api.GetFlights(w, r, flightsCollection)
	}

	getHotelsHandler := func(w http.ResponseWriter, r *http.Request) {
		api.GetHotels(w, r, hotelsCollection)
	}

	router.HandleFunc("/flight", getFlightsHandler).Methods("GET")
	router.HandleFunc("/hotel", getHotelsHandler).Methods("GET")

	if err := http.ListenAndServe(constants.BASE_URL, router); err != nil {
		fmt.Println(err)
	}
}
