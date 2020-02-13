package main
/*
	TODO
	- Check error handling
	- Decouple code
*/
import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"time"
)

type Person struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Firstname string `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname string `json:"lastname,omitempty" bson:"lastname,omitempty"`
}

type FoodOrderRequest struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string `json:"name,omitempty" bson:"name,omitempty"`
	Quantity int `json:"quantity,omitempty" bson:"quantity,omitempty"`
}

type FoodOrderResponse struct {
	ID primitive.ObjectID `json:"transactionId,omitempty" bson:"transactionId,omitempty"`
	Name string `json:"name,omitempty" bson:"name,omitempty"`
	Quantity int `json:"quantity,omitempty" bson:"quantity,omitempty"`
}

var client *mongo.Client

func CreateFoodOrderRequest(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	response.Write([]byte ("foodOrderResponse "))
	var foodOrderRequest FoodOrderRequest
	json.NewDecoder(request.Body).Decode(&foodOrderRequest)
	collection := client.Database("foodOrder").Collection("foodRequests")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx, foodOrderRequest)
	json.NewEncoder(response).Encode(result)
}

func CreatePersonEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var person Person
	json.NewDecoder(request.Body).Decode(&person)
	collection := client.Database("thepolyglotdeveloper").Collection("people")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx, person)
	json.NewEncoder(response).Encode(result)
}

func GetPeopleEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var people []Person
	collection := client.Database("thepolyglotdeveloper").Collection("people")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `"}`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var person Person
		cursor.Decode(&person)
		people = append(people, person)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(people)
}

func GetPersonEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	response.Write([]byte ("person "))
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var person Person
	collection := client.Database("thepolyglotdeveloper").Collection("people")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := collection.FindOne(ctx, Person{ID: id}).Decode(&person)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte (`{ "message": "` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(person)
}

func GetFoodOrderRequestByTransactionId(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	//response.Write([]byte ("foodResponseByTransactionId "))
	params := mux.Vars(request)
	transactionId, _ := primitive.ObjectIDFromHex(params["transactionId"])
	var foodOrderResponse FoodOrderResponse
	collection := client.Database("foodOrder").Collection("foodRequests")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := collection.FindOne(ctx, FoodOrderResponse{ID: transactionId}).Decode(&foodOrderResponse)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte (`{ "message": "` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(foodOrderResponse)
}

func main () {
	fmt.Println("Starting the application ...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	//Use Atlas Mongo cluster connection String instead
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	router := mux.NewRouter()
	router.HandleFunc("/person", CreatePersonEndpoint).Methods("POST")
	router.HandleFunc("/food/order/request", CreateFoodOrderRequest).Methods("POST")
	router.HandleFunc("/people", GetPeopleEndpoint).Methods("GET")
	router.HandleFunc("/person/{id}", GetPersonEndpoint).Methods("GET")
	router.HandleFunc("/food/order/transactionId/{transactionId}", GetFoodOrderRequestByTransactionId).Methods("GET")
	http.ListenAndServe(":8080", router)
}