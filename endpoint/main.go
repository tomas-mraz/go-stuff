package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

var (
	requests    []string
	requestsMu  sync.Mutex
	port        = "8080"
	maxRequests = 30
	htmlPage    = "<html><head><title>BlackHole endpoint</title></head>" +
		"<body><h1>Tady žijí dva endpointy</h1><ul>" +
		"<li><b>/ping</b> ... který přijímá requesty</li>" +
		"<li><a href=\"/log\"><b>/log</b></a> ... který zobrazuje posledních 30 přijatých requestů</li>" +
		"</ul></body></html>"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
	requestsMu.Lock()
	defer requestsMu.Unlock()
	// logging body
	buf, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Body read Error: " + err.Error())
		return
	}
	formattedBody := ""
	if len(buf) != 0 {
		formattedBody = "\n" + string(buf)
	}

	requests = append(requests, time.Now().Format(time.RFC3339)+" - "+r.Method+" "+r.URL.Path+formattedBody+"\n")

	if len(requests) > maxRequests {
		// Remove the oldest request to keep the list within the limit
		requests = requests[len(requests)-maxRequests:]
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["message"] = "Status OK"
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		slog.Error("Error happened in JSON marshal: " + err.Error())
	}
	w.Write(jsonResp)
	return
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(htmlPage))
	if err != nil {
		slog.Error("writing homepage error: " + err.Error())
		return
	}
}

func logHandler(w http.ResponseWriter, r *http.Request) {
	requestsMu.Lock()
	defer requestsMu.Unlock()

	// Print requests in reverse order
	for i := len(requests) - 1; i >= 0; i-- {
		fmt.Fprintln(w, requests[i])
	}
}

func main() {
	slog.Info("Start listening on port " + port)
	http.Handle("/", http.HandlerFunc(homeHandler))
	http.Handle("/ping", http.HandlerFunc(pingHandler))
	http.Handle("/log", http.HandlerFunc(logHandler))
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		slog.Error(err.Error())
		return
	}
}
