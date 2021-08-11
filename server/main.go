package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

//const SortDataingAlgorithms = []string{"bubble", "insertion", "merge"}

func setupCorsResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")
}

type SortData struct {
	Elements  []int  `json:"elements"`
	Algorithm string `json:"algorithm"`
}

func newSortData(reader io.Reader) (*SortData, error) {
	sort := new(SortData)
	err := json.NewDecoder(reader).Decode(sort)
	if err != nil {
		return nil, err
	}
	return sort, nil
}

func swap(arr []int, i int, j int) {
	temp := arr[i]
	arr[i] = arr[j]
	arr[j] = temp
}

type Step struct {
	ID   int   `json:"id"`
	List []int `json:"list"`
}

type SortAlgorithm func(elements []int) []Step

var id int = 0

func bubbleSort(elements []int) []Step {
	steps := []Step{}
	for i := 0; i < len(elements)-1; i++ {
		for j := 0; j < len(elements)-i-1; j++ {
			if elements[j] > elements[j+1] {
				swap(elements, j+1, j)
				temp := make([]int, len(elements))
				// create a copy of the original slice so the order is preserved through each iteration
				copy(temp, elements)
				steps = append(steps, Step{ID: id, List: temp})
				id++
			}
		}
	}
	return steps
}

func sort(data *SortData) []Step {
	algorithms := map[string]SortAlgorithm{
		"bubble": bubbleSort,
	}
	steps := algorithms[data.Algorithm](data.Elements)
	return steps
}

func handleSortData(w http.ResponseWriter, r *http.Request) {
	setupCorsResponse(&w, r)
	if r.Method == http.MethodOptions {
		return
	}
	sortData, err := newSortData(r.Body)
	if err != nil {
		log.Printf("Unable to decode input. Error: %+v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	steps := sort(sortData)
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
