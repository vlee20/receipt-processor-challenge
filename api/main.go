// main.go
package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// Handler for the root endpoint

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

func isIntegral(val float64) bool {
	return val == float64(int(val))
}

func GetPoints(w http.ResponseWriter, r *http.Request) {
	// if r.Method != http.MethodGet {
	// 	http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	// 	return
	// }
	var points int = 0
	var res = make(map[string]int)
	vars := mux.Vars(r)
	id := vars["id"]
	if global_memory[id] != nil {
		fmt.Println(global_memory[id])
		// One point for every alphanumeric character in the retailer name
		points = points + len(global_memory[id].Retailer)
		// 50 points if the total is a round dollar amount with no cents.
		t2, _ := strconv.ParseFloat(global_memory[id].Total, 8)
		if isIntegral(t2) {
			points = points + 50
		}
		// 25 points if the total is a multiple of 0.25
		if math.Mod(t2, 0.25) == 0 {
			points = points + 25
		}
		// 5 points for every two items on the receipt
		var item_size = len(global_memory[id].Items)
		points = points + int(item_size/2)
		// If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer. The result is the number of points earned.
		l := len(global_memory[id].Items)
		item := global_memory[id].Items
		for i := 0; i < l; i++ {
			price := item[i]["price"]
			conv_price, _ := strconv.ParseFloat(price, 8)
			trim_desc := strings.TrimSpace(item[i]["shortDescription"])
			fmt.Println(trim_desc)
			if len(trim_desc)%3 == 0 {
				points = points + int(math.Ceil(conv_price*0.2))
			}
		}
		res["points"] = points
		json.NewEncoder(w).Encode(res)
	} else {
		http.Error(w, "Item not found", http.StatusNotFound)
	}
	fmt.Println(id)
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
