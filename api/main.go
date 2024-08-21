// main.go
package main

import (
    "fmt"
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

func main() {
    r := mux.NewRouter()
    
    // Define routes
    r.HandleFunc("/", HomeHandler).Methods("GET")
    r.HandleFunc("/api/greet/{name:[a-zA-Z0-9]+}", GreetHandler).Methods("GET")
    
    // Start the server
    fmt.Println("Starting server on :8080")
    http.ListenAndServe(":8080", r)
}
