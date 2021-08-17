package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/algorithm-visualizer/server/sort"
	"github.com/gorilla/mux"
)

//const SortDataingAlgorithms = []string{"bubble", "insertion", "merge"}

func setupCorsResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")
}

func handleSortData(w http.ResponseWriter, r *http.Request) {
	setupCorsResponse(&w, r)
	if r.Method == http.MethodOptions {
		return
	}
	input, err := sort.DecodeInput(r.Body)
	if err != nil {
		log.Printf("Unable to decode input. Error: %+v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	steps := sort.Do(input)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(steps)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/sort", handleSortData).Methods(http.MethodPost, http.MethodOptions)

	port := ":8080"
	log.Printf("Listening on port %s", port)
	err := http.ListenAndServe(port, router)
	if err != nil {
		log.Fatal(err)
	}
}
