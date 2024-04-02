package main

import (
	"context"
	"crypto/tls"
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"grpc2/aaa_cz"
	"log"
	"time"
)

const (
	defaultName = "world"
	servername  = "backend.mydomain.com"
)

var (
	addr = flag.String("addr", servername+":50051", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func main() {
	flag.Parse()

	creds := credentials.NewTLS(&tls.Config{})
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := aaa_cz.NewGreeterClient(conn)

	// Contact the server and print out its response.
    // timeout is needed to wait for cert generation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &aaa_cz.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("Greeting: %s", r.GetMessage())
}
