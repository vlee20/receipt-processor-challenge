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
	"time"
	"unicode"

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

// int check
func isIntegral(val float64) bool {
	return val == float64(int(val))
}

// gen_id layout
func insertCharAt(s string, index int, char rune) string {
	// convert string to a slice of runes
	runes := []rune(s)

	// handle out-of-range index cases
	if index < 0 {
		index = 0
	} else if index > len(runes) {
		index = len(runes)
	}

	// insert the character by slicing and concatenating
	new_string := string(runes[:index]) + string(char) + string(runes[index:])

	return new_string
}

func CreateReceipt(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	// randomInt := rand.Intn(1000)
	var receipt Receipt
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&receipt); err != nil || len(receipt.Items) == 0 || len(receipt.Retailer) == 0 {
		http.Error(w, "The receipt is invalid", http.StatusBadRequest)
		return
	}

	// parse string date to golang date
	_, err_date := time.Parse("2006-01-02", receipt.PurchaseDate)
	if err_date != nil {
		http.Error(w, "The receipt is invalid", http.StatusBadRequest)
		return
	}

	// parse string date to golang time
	_, err_time := time.Parse("15:04", receipt.PurchaseTime)
	if err_time != nil {
		http.Error(w, "The receipt is invalid", http.StatusBadRequest)
		return
	}

	// parse total to make sure its a float
	total_check, err_total := strconv.ParseFloat(receipt.Total, 64)
	if err_total != nil || total_check < 0 {
		http.Error(w, "The receipt is invalid", http.StatusBadRequest)
		return
	}

	// parse price to make sure its a float
	length_items := len(receipt.Items)
	for i := 0; i < length_items; i++ {
		item := receipt.Items[i]
		price_check, err_price := strconv.ParseFloat(item["price"], 64)
		if err_price != nil || price_check < 0 {
			http.Error(w, "The receipt is invalid", http.StatusBadRequest)
			return
		}
	}

	var res_id = make(map[string]string)
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	gen_id := make([]byte, 32)
	for i := range gen_id {
		gen_id[i] = letters[rand.Intn(len(letters))]
	}

	dash := '-'
	ind_1 := 8
	ind_2 := 13
	ind_3 := 18
	ind_4 := 23

	gen_id_f := string(gen_id)
	gen_id_f = insertCharAt(gen_id_f, ind_1, dash)
	gen_id_f = insertCharAt(gen_id_f, ind_2, dash)
	gen_id_f = insertCharAt(gen_id_f, ind_3, dash)
	gen_id_f = insertCharAt(gen_id_f, ind_4, dash)

	global_memory[gen_id_f] = &receipt
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	res_id["id"] = gen_id_f
	json.NewEncoder(w).Encode(res_id)
}

func GetPoints(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var points int = 0
	var res = make(map[string]int)
	vars := mux.Vars(r)
	id := vars["id"]
	if global_memory[id] != nil {
		// One point for every alphanumeric character in the retailer name
		retailer := global_memory[id].Retailer
		for _, char := range retailer {
			if unicode.IsLetter(char) || unicode.IsDigit(char) {
				points = points + 1
			}
		}
		// 50 points if the total is a round dollar amount with no cents.
		t2, _ := strconv.ParseFloat(global_memory[id].Total, 64)
		if isIntegral(t2) {
			points = points + 50
		}
		// 25 points if the total is a multiple of 0.25
		if math.Mod(t2, 0.25) == 0 {
			points = points + 25
		}
		// 5 points for every two items on the receipt
		var item_size = len(global_memory[id].Items)
		points = points + (int(item_size/2) * 5)
		// If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer. The result is the number of points earned.
		l := len(global_memory[id].Items)
		item := global_memory[id].Items
		for i := 0; i < l; i++ {
			price := item[i]["price"]
			conv_price, _ := strconv.ParseFloat(price, 64)
			trim_desc := strings.TrimSpace(item[i]["shortDescription"])
			if len(trim_desc)%3 == 0 {
				points = points + int(math.Ceil(conv_price*0.2))
			}
		}
		// 6 points if the day in the purchase date is odd
		date_layout := "2006-01-02"
		parsedDate, errDate := time.Parse(date_layout, global_memory[id].PurchaseDate)
		if errDate != nil {
			fmt.Println("Error parsing date:", errDate)
			return
		}
		day := parsedDate.Day()
		if day%2 == 1 {
			points = points + 6
		}
		// 10 points if the time of purchase is after 2:00pm and before 4:00pm
		purchase_time := global_memory[id].PurchaseTime
		time_layout := "15:04"
		parsedTime, errTime := time.Parse(time_layout, purchase_time)
		if errTime != nil {
			fmt.Println("Error parsing date:", errTime)
			return
		}
		hour := parsedTime.Hour()
		if hour >= 14 && hour <= 16 {
			points = points + 10
		}
		res["points"] = points
		json.NewEncoder(w).Encode(res)
	} else {
		http.Error(w, "No receipt found for that id", http.StatusNotFound)
	}
}

func main() {
	r := mux.NewRouter()

	// routes
	r.HandleFunc("/receipts/process", CreateReceipt).Methods("POST")

	r.HandleFunc("/receipts/{id}/points", GetPoints).Methods("GET")

	// start the server
	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}
