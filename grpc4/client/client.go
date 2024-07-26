package main

import (
	"context"
	"crypto/tls"
	"flag"
	"github.com/johanbrandhorst/certify"
	"github.com/johanbrandhorst/certify/issuers/vault"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"grpc3/aaa_cz"
	"log"
	"net/url"
	"time"
)

const (
	defaultName = "world2"
	servername  = "backend.cubyte.space"
)

var (
	addr = flag.String("addr", servername+":50051", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func main() {
	flag.Parse()

	// Certify
	issuer := &vault.Issuer{
		URL: &url.URL{
			Scheme: "http",
			Host:   "security.cubyte.online:8200",
		},
		AuthMethod: vault.ConstantToken("Bearer eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJxYTJKX0w4VkRGbkFGZWxmY1RWN0JiZTNmNkY1cmpGT2ExWGo1cUdsVW1RIn0.eyJleHAiOjE3MTY3MzY4MDEsImlhdCI6MTcxNjczNjUwMSwiYXV0aF90aW1lIjoxNzE2NzM2NTAxLCJqdGkiOiIwMmEzOThlZS0yYzgyLTQ2MDEtOGRiYi0wNzY4OGE5MzczMGEiLCJpc3MiOiJodHRwczovL2xvZ2luLmN1Ynl0ZS5vbmxpbmUvcmVhbG1zL2N1Ynl0ZSIsImF1ZCI6ImFjY291bnQiLCJzdWIiOiJiM2ZjZTY4NC04YjQxLTRhNzQtOWFkMC04ZmMzZDBlZWEyMzkiLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiJkZW1vLWNsaWVudCIsInNlc3Npb25fc3RhdGUiOiI1ODg2OGRkNy1iN2ViLTRmZjAtYTVhZi01ZWE3ZTVlNDcwNGEiLCJhY3IiOiIxIiwiYWxsb3dlZC1vcmlnaW5zIjpbImh0dHA6Ly8xMjcuMC4wLjEiXSwicmVhbG1fYWNjZXNzIjp7InJvbGVzIjpbIm9mZmxpbmVfYWNjZXNzIiwiZGVmYXVsdC1yb2xlcy1jdWJ5dGUiLCJ1bWFfYXV0aG9yaXphdGlvbiJdfSwicmVzb3VyY2VfYWNjZXNzIjp7ImFjY291bnQiOnsicm9sZXMiOlsibWFuYWdlLWFjY291bnQiLCJtYW5hZ2UtYWNjb3VudC1saW5rcyIsInZpZXctcHJvZmlsZSJdfX0sInNjb3BlIjoib3BlbmlkIGVtYWlsIHByb2ZpbGUiLCJzaWQiOiI1ODg2OGRkNy1iN2ViLTRmZjAtYTVhZi01ZWE3ZTVlNDcwNGEiLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwicHJlZmVycmVkX3VzZXJuYW1lIjoiZGVtbyIsImdpdmVuX25hbWUiOiIiLCJmYW1pbHlfbmFtZSI6IiIsImVtYWlsIjoiemlpa0B6b25hNjQuY3oifQ.cBB71h6NHzJmvzCCu-EeJkrgunqMMNn9wAVZQEk31iEx5qSQpId8CUFYHLd88jvr1WzQ8qAbzLBIXdgGZNN-a8mtWl_uw_v14PnJeH6c-mlSliGRW41soxG0H4CuJ_KtfU4O79k-DcnwnsQz2lDfIOJnHkky8p3muOi1EVXLQo2VzW7bscU399xoSgFCtC3xPS1fCdoIJQmBs3nqt8gL2w_y6W7PdGY0osYOlsA_-Y_btECHYIpNeZOlhWLwV7YPxE4yIr8GIaDy3zQZsohexpDQj1t5Z-daZcIF-x7Fqrzv-_1fx2C0nacHevd9A3hh7X5MZG84jAQ2Hr_q0nDZ6Q"),
		Role:       "cubyte-dot-space",
		TimeToLive: 8 * time.Hour,
		Mount:      "pki_int",
	}
	cert := &certify.Certify{
		CommonName:  "ziik.user.cubyte.space",
		Issuer:      issuer,
		Cache:       certify.NewMemCache(),
		RenewBefore: 8 * time.Hour,
	}

	// Let's Encrypt
	config := &tls.Config{
		GetClientCertificate: cert.GetClientCertificate,
	}
	creds := credentials.NewTLS(config)
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := aaa_cz.NewGreeterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &aaa_cz.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("Greeting: %s", r.GetMessage())
}
