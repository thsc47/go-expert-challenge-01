package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type DolarValue struct {
	UsdDolar struct {
		Bid       string `json:"bid"`
		Timestamp string `json:"timestamp"`
	} `json:"USDBRL"`
}

func main() {
	http.HandleFunc("/", serverHandler)
	http.ListenAndServe(":8080", nil)
}

func serverHandler(w http.ResponseWriter, r *http.Request) {
	bid, err := getActualDolarValue()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bid)
}

func getActualDolarValue() (DolarValue, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return DolarValue{}, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return DolarValue{}, err
	}
	defer resp.Body.Close()
	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return DolarValue{}, err
	}
	var bid DolarValue
	err = json.Unmarshal(res, &bid)
	if err != nil {
		return DolarValue{}, err
	}
	return bid, nil
}
