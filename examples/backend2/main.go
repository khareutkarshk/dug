package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
	Service string `json:"service"`
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := Response{
		Message: "Hello from backend",
		Service: "backend-service-2",
	}

	json.NewEncoder(w).Encode(response)

}

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/health", health)

	log.Println("Backend listening on :3002")
	log.Fatal(http.ListenAndServe(":3002", nil))
}
