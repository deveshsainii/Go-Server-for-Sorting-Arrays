// main.go
package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"sync"
	"time"
)

type SortRequest struct {
	ToSort [][]int `json:"to_sort"`
}

type SortResponse struct {
	SortedArrays [][]int `json:"sorted_arrays"`
	TimeNs       int64   `json:"time_ns"`
}

func sortSequential(input [][]int) ([][]int, int64) {
	startTime := time.Now()

	for i := range input {
		sort.Ints(input[i])
	}

	return input, time.Since(startTime).Nanoseconds()
}

func sortConcurrent(input [][]int) ([][]int, int64) {
	startTime := time.Now()

	var wg sync.WaitGroup
	wg.Add(len(input))

	for i := range input {
		go func(i int) {
			defer wg.Done()
			sort.Ints(input[i])
		}(i)
	}

	wg.Wait()

	return input, time.Since(startTime).Nanoseconds()
}

func processSingleHandler(w http.ResponseWriter, r *http.Request) {
	var req SortRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sortedArrays, timeNs := sortSequential(req.ToSort)

	res := SortResponse{
		SortedArrays: sortedArrays,
		TimeNs:       timeNs,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func processConcurrentHandler(w http.ResponseWriter, r *http.Request) {
	var req SortRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sortedArrays, timeNs := sortConcurrent(req.ToSort)

	res := SortResponse{
		SortedArrays: sortedArrays,
		TimeNs:       timeNs,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func main() {
	http.HandleFunc("/process-single", processSingleHandler)
	http.HandleFunc("/process-concurrent", processConcurrentHandler)

	http.ListenAndServe(":8000", nil)
}
