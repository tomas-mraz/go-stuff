package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"net/url"
)

var (
	EchoBaseURL, _ = url.Parse("https://test.company.com")
	// Hop-by-hop headers. These are removed when sent to the backend.
	// http://www.w3.org/Protocols/rfc2616/rfc2616-sec13.html
	hopHeaders = []string{
		"Connection",
		"Keep-Alive",
		"Proxy-Authenticate",
		"Proxy-Authorization",
		"Te", // canonicalized version of "TE"
		"Trailers",
		"Transfer-Encoding",
		"Upgrade",
	}
)

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func delHopHeaders(header http.Header) {
	for _, h := range hopHeaders {
		header.Del(h)
	}
}

type proxy struct {
}

func (p *proxy) ServeHTTP(responseWriter http.ResponseWriter, req *http.Request) {
	log.Printf("Accept request %s to %s\n", req.Method, req.URL)
	destinationURL := EchoBaseURL.ResolveReference(req.URL).String()

	// create new request with old body and new destination
	newReq, err := http.NewRequest(req.Method, destinationURL, req.Body)
	// copy headers to new request
	delHopHeaders(req.Header)
	copyHeader(newReq.Header, req.Header)

	client := &http.Client{}

	resp, err := client.Do(newReq)
	if err != nil {
		http.Error(responseWriter, "Server Error", http.StatusInternalServerError)
		log.Fatal("ServeHTTP:", err)
	}
	defer resp.Body.Close()
	log.Printf("Received response %s from %s\n", resp.Status, destinationURL)

	delHopHeaders(resp.Header)
	copyHeader(responseWriter.Header(), resp.Header)
	responseWriter.WriteHeader(resp.StatusCode)
	io.Copy(responseWriter, resp.Body)
}

// the base of code is from https://gist.github.com/yowu/f7dc34bd4736a65ff28d?permalink_comment_id=4068010
func main() {
	var addr = flag.String("addr", "0.0.0.0:8484", "The addr of the application.")
	flag.Parse()

	handler := &proxy{}

	log.Println("Starting proxy server on", *addr)
	if err := http.ListenAndServe(*addr, handler); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
