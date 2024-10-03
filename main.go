package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Car struct {
	ID     int    `json:"id"`
	Make   string `json:"make"`
	Model  string `json:"model"`
	Year   int    `json:"year"`
	Status string `json:"status"` // "available", "in maintenance", "rented"
	Owner  string `json:"owner"`
}

var cars []Car
var nextID = 1

// Create a new car
func createCar(w http.ResponseWriter, r *http.Request) {
	var car Car
	json.NewDecoder(r.Body).Decode(&car)
	car.ID = nextID
	nextID++
	car.Status = "available"
	cars = append(cars, car)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(car)
}

// Get all cars
func getCars(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cars)
}

// Get a car by ID
func getCarByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	for _, car := range cars {
		if car.ID == id {
			json.NewEncoder(w).Encode(car)
			return
		}
	}
	http.Error(w, "Car not found", http.StatusNotFound)
}

// Update a car
func updateCar(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var updatedCar Car
	json.NewDecoder(r.Body).Decode(&updatedCar)
	for i, car := range cars {
		if car.ID == id {
			cars[i] = updatedCar
			cars[i].ID = id
			json.NewEncoder(w).Encode(cars[i])
			return
		}
	}
	http.Error(w, "Car not found", http.StatusNotFound)
}

// Delete a car
func deleteCar(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	for i, car := range cars {
		if car.ID == id {
			cars = append(cars[:i], cars[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "Car not found", http.StatusNotFound)
}

func main() {
	r := mux.NewRouter()

	// Routes
	r.HandleFunc("/cars", createCar).Methods("POST")
	r.HandleFunc("/cars", getCars).Methods("GET")
	r.HandleFunc("/cars/{id}", getCarByID).Methods("GET")
	r.HandleFunc("/cars/{id}", updateCar).Methods("PUT")
	r.HandleFunc("/cars/{id}", deleteCar).Methods("DELETE")

	log.Println("Server Running On Port 8080")
	http.ListenAndServe(":8080", r)
}
