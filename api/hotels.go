package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Hotel struct {
	HotelName  string
	TotalPrice int
}

type CheapestHotelResponseData struct {
	City         string
	CheckInDate  string `json:"Check In Date"`
	CheckOutDate string `json:"Check Out Date"`
	Hotel        string
	Price        int
}

func GetHotels(w http.ResponseWriter, r *http.Request, hotelsCollection *mongo.Collection) {
	params := r.URL.Query()
	checkInDate, err := time.Parse("2006-01-02", params.Get("checkInDate"))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}

	checkOutDate, err := time.Parse("2006-01-02", params.Get("checkOutDate"))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}

	destination := strings.Title(strings.ToLower(params.Get("destination")))

	cheapestHotelsAggregationPipeline := mongo.Pipeline{
		bson.D{
			{"$match", bson.D{
				{"city", destination},
				{"date", bson.D{
					{"$gte", checkInDate},
					{"$lte", checkOutDate},
				}},
			}},
		}, // Filter by destination, checkInDate, and checkOutDate
		bson.D{
			{"$group", bson.D{
				{"_id", "$hotelName"}, // Group by the "hotelName" field
				{"hotelName", bson.D{
					{"$first", "$hotelName"}, // Retrieve the hotelName using $first operator
				}},
				{"totalPrice", bson.D{
					{"$sum", "$price"}, // Calculate the total price of staying in each hotel for the period
				}},
			}},
		},
		bson.D{
			{"$sort", bson.D{
				{"totalPrice", 1}, // Sort by totalPrice in ascending order
			}},
		},
	}

	cheapestHotels, err := hotelsCollection.Aggregate(context.Background(), cheapestHotelsAggregationPipeline)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}

	var cheapestHotelsResponseSlice []CheapestHotelResponseData
	checkInDateString := checkInDate.Format("2006-01-02")
	checkOutDateString := checkOutDate.Format("2006-01-02")

	for cheapestHotels.Next(context.Background()) {
		var hotel Hotel
		if err := cheapestHotels.Decode(&hotel); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
		}
		cheapestHotelsResponseSlice = append(cheapestHotelsResponseSlice, CheapestHotelResponseData{
			City:         destination,
			CheckInDate:  checkInDateString,
			CheckOutDate: checkOutDateString,
			Hotel:        hotel.HotelName,
			Price:        hotel.TotalPrice,
		})

	}
	cheapestHotelsResponseJSON, _ := json.Marshal(cheapestHotelsResponseSlice)

	w.Header().Set("Content-Type", "application/json")
	w.Write(cheapestHotelsResponseJSON)

}
