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
		AuthMethod: vault.ConstantToken("hvs.CAESIND8NZ82wSJwSfljQxXdwDNsgoh2Vr86cBtZRke1V8v5Gh4KHGh2cy5aaHllQnl2cTJMbUl3TTc3ZTBiNTZmVmo"),
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
