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

func IdGenerator(w http.ResponseWriter, r *http.Request) {
	// var uppercase_letters string = "ABCDEFGHJKLMNPQRSTUVWXYZ"
	// var lowercase_letters string = "abcdefghijkmnopqrstuvwxyz"
	// var digits string = "0123456789"
	randomInt := rand.Intn(1000)
	// const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	// b := make([]byte, 16)
	// for i := range b {
	//     b[i] = letters[rand.Intn(len(letters))]
	// }
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

var global_memory = make(map[string]*Receipt)

func CreateReceipt(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	// randomInt := rand.Intn(1000)
	var receipt Receipt
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&receipt); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	// Respond with the received data
	response := fmt.Sprintf("Receipt created: %s", receipt.Retailer)
	// rand := fmt.Sprintf("%d", randomInt)

	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 16)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	global_memory[string(b)] = &receipt
	fmt.Println(global_memory)
	fmt.Println(global_memory[string(b)])
	fmt.Println(receipt.Items[0]["price"])
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message":"` + response + `"}, {"id":"` + string(b) + `"}`))
}

func GetPoints(w http.ResponseWriter, r *http.Request) {
	// if r.Method != http.MethodGet {
	// 	http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	// 	return
	// }
	vars := mux.Vars(r)
	id := vars["id"]
	// id := r.URL.Query().Get("id")
	fmt.Println(id)
	w.Write([]byte(`{"id":"` + id + `"}`))
}

func main() {
	r := mux.NewRouter()

	// Define routes

	// r.HandleFunc("/receipts/process", IdGenerator).Methods("POST")

	r.HandleFunc("/receipts/process", CreateReceipt).Methods("POST")

	r.HandleFunc("/receipts/{id}/points", GetPoints).Methods("GET")

	// fmt.Println("Server is listening on port 8080")
	// if err := http.ListenAndServe(":8080", nil); err != nil {
	// 	fmt.Println("Server failed to start:", err)
	// }
	// Start the server
	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}
