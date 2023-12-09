package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"sync"
	"time"
)

type RequestPayload struct {
	ToSort [][]int `json:"to_sort"`
}

type ResponsePayload struct {
	SortedArrays [][]int `json:"sorted_arrays"`
	TimeNS       int64   `json:"time_ns"`
}

func processSingle(w http.ResponseWriter, r *http.Request) {
	var reqPayload RequestPayload
	if err := json.NewDecoder(r.Body).Decode(&reqPayload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	startTime := time.Now()
	sortedArrays := make([][]int, len(reqPayload.ToSort))
	for i, arr := range reqPayload.ToSort {
		sort.Ints(arr)
		sortedArrays[i] = arr
	}
	timeTaken := time.Since(startTime).Nanoseconds()

	response := ResponsePayload{
		SortedArrays: sortedArrays,
		TimeNS:       timeTaken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func processConcurrent(w http.ResponseWriter, r *http.Request) {
	var reqPayload RequestPayload
	if err := json.NewDecoder(r.Body).Decode(&reqPayload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	startTime := time.Now()
	var wg sync.WaitGroup
	sortedArrays := make([][]int, len(reqPayload.ToSort))
	for i, arr := range reqPayload.ToSort {
		wg.Add(1)
		go func(i int, arr []int) {
			defer wg.Done()
			sort.Ints(arr)
			sortedArrays[i] = arr
		}(i, arr)
	}
	wg.Wait()
	timeTaken := time.Since(startTime).Nanoseconds()

	response := ResponsePayload{
		SortedArrays: sortedArrays,
		TimeNS:       timeTaken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/process-single", processSingle)
	http.HandleFunc("/process-concurrent", processConcurrent)

	port := ":8000"
	http.ListenAndServe(port, nil)
}
