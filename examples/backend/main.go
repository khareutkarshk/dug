package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type Response struct {
	Message string `json:"message"`
	Service string `json:"service"`
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Backend", "3001")
	response := Response{
		Message: "Hello from backend",
		Service: "backend-service",
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/health", health)
	http.HandleFunc("/slow", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Second)
		_, _ = fmt.Fprintln(w, os.Getenv("SERVICE_NAME"))
	})

	log.Println("Backend listening on :3001")
	log.Fatal(http.ListenAndServe(":3001", nil))
}
