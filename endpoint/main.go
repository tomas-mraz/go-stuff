package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["message"] = "Status OK"
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
	return
}

func main() {
	fmt.Println("Start listening on port 8080 ...")
	handler := http.HandlerFunc(handleRequest)
	http.Handle("/ping", handler)
	http.ListenAndServe(":8080", nil)
}
