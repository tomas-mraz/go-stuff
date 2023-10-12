package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log/slog"
	"net/http"
	"os"
	"sort"
	"sync"
	"time"
)

type endpoint struct {
	lastUpdate time.Time
	pointer    int
	requests   [maxRequests]string
}

const (
	maxRequests = 30

	htmlPageStart = "<html><head><title>BlackHole endpoint</title></head>" +
		"<body><h1>Tady žijí testovací endpointy</h1><ul>" +
		"<li><b>/</b> ... tahle stránka</li>" +
		"<li><b>/log</b> ... již neukazuje souhrnné logy (funguje jako přijímací endpoint viz /* níže)</li>" +
		"<li><b>/*</b> ... jakýkoli název kontextu přijímá requesty a vrací 200:{\"message\":\"Status OK\"}</li>" +
		"<li><a href=\"/*/log\">/*<b>/log</b></a> ... <KONTEXT>/log zobrazuje posledních 30 přijatých requestů pro daný kontext (poslední je nahoře)</li>" +
		"</ul>" +
		"Příklady:<br>" +
		"ferda.zona64.cz:8080<b>/magda</b> ... bude mít logy na /magda/log<br>" +
		"ferda.zona64.cz:8080<b>/magda/123</b> ... bude mít logy na /magda/123/log<br>" +
		"<br><br>" +
		"<b><u>Používané endpointy</u></b>:" +
		"<ul>"
	htmlPageEnd = "</ul></body></html>"
)

var (
	endpoints  map[string]*endpoint
	requestsMu sync.Mutex
)

func pingHandler(writer http.ResponseWriter, request *http.Request) {
	// skip requests to favicon
	if request.URL.Path == "/favicon.ico" {
		return
	}

	requestsMu.Lock()
	defer requestsMu.Unlock()

	params := mux.Vars(request)
	name := params["name"]

	// logging body
	buf, err := io.ReadAll(request.Body)
	if err != nil {
		slog.Error("Body read Error: " + err.Error())
		return
	}
	formattedBody := ""
	if len(buf) != 0 {
		// body je na nové řádce
		formattedBody = "\n" + string(buf)
	}

	e, ok := endpoints[name]
	if !ok {
		f := endpoint{time.Now(), 0, [maxRequests]string{}}
		endpoints[name] = &f
		e = &f
	}
	e.lastUpdate = time.Now()
	e.requests[e.pointer] = time.Now().Format(time.RFC3339) + " - " + request.Method + " " + request.URL.Path + formattedBody + "\n"

	e.pointer++
	if e.pointer >= maxRequests {
		e.pointer = 0
	}

	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["message"] = "Status OK"
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		slog.Error("Error happened in JSON marshal: " + err.Error())
	}
	_, err = writer.Write(jsonResp)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	return
}

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug(r.URL.Path)
	if r.URL.Path == "/" {
		keys := make([]string, 0, len(endpoints))
		for key := range endpoints {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		a := ""
		for _, name := range keys {
			a += "<li><a href=\"" + name + "/log\">" + name + "</a> (" + endpoints[name].lastUpdate.Format("15:04:05 02.01.2006") + ")</li>\n"
		}

		b := htmlPageStart + a + htmlPageEnd
		_, err := w.Write([]byte(b))

		if err != nil {
			slog.Error("writing homepage error: " + err.Error())
			return
		}
	}
}

func logPageHandler(writer http.ResponseWriter, request *http.Request) {
	requestsMu.Lock()
	defer requestsMu.Unlock()

	params := mux.Vars(request)
	name := params["name"]

	e := endpoints[name]
	if e == nil {
		_, err := fmt.Fprintln(writer, "tento endpoint jeste nikdo nezavolal")
		if err != nil {
			slog.Error(err.Error())
		}
		return
	}

	// Print requests in reverse order
	for i := 0; i < maxRequests; i++ {
		a := e.pointer - i - 1 // pointer target to index of next item (not the last one)
		if a < 0 {
			a = maxRequests + a
		}
		_, err := fmt.Fprintln(writer, e.requests[a])
		if err != nil {
			slog.Error(err.Error())
			return
		}
	}
}

func main() {
	endpoints = make(map[string]*endpoint)
	port := os.Args[1]
	slog.Info("Start listening on port " + port)

	rtr := mux.NewRouter()
	rtr.HandleFunc("/{name:.+}/log", logPageHandler)
	rtr.HandleFunc("/{name:.+}", pingHandler)
	rtr.HandleFunc("/", homePageHandler)
	http.Handle("/", rtr)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		slog.Error(err.Error())
		return
	}
}
