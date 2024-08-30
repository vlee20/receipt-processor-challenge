// main.go
package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/gorilla/mux"
)

// Handler for the root endpoint
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Go web API!!")
}

// Handler for the /api/greet endpoint
func GreetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	fmt.Fprintf(w, "Hello, %s!", name)
}

func IdGenerator(w http.ResponseWriter, r *http.Request) {
	randomInt := rand.Intn(1000)
	// return map[string]int{"id": randomInt}
	fmt.Fprintf(w, `{"id:"`+string(randomInt)+`}`)
}

type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchasedate"`
	PurchaseTime string `json:"purchasetime"`
	Items        []map[string]string
	Total        string `json:"total"`
}

func CreateReceipt(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	randomInt := rand.Intn(1000)
	var receipt Receipt
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&receipt); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	// Respond with the received data
	response := fmt.Sprintf("Receipt created: %s", receipt.Retailer)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message":"` + response + `"}`))
	w.Write([]byte(`{"id":"` + string(randomInt) + `"}`))
}

func main() {
	// r := mux.NewRouter()

	// Define routes
	// r.HandleFunc("/", HomeHandler).Methods("GET")
	// r.HandleFunc("/api/greet/{name:[a-zA-Z0-9]+}", GreetHandler).Methods("GET")

	// r.HandleFunc("/receipts/process", IdGenerator).Methods("POST")

	http.HandleFunc("/receipts/process", CreateReceipt)

	fmt.Println("Server is listening on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server failed to start:", err)
	}
	// Start the server
	// fmt.Println("Starting server on :8080")
	// http.ListenAndServe(":8080", r)
}
