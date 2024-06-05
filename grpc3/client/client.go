package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"grpc3/aaa_cz"
	"log/slog"
	"os"
	"time"
)

const (
	defaultName = "world"
	servername  = "backend.cubyte.space"
)

var (
	addr = flag.String("addr", servername+":50051", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func main() {
	flag.Parse()

	// get client certificate
	data, _ := os.ReadFile("my-cert.pem")
	block, _ := pem.Decode(data)
	clientCertificate, _ := x509.ParseCertificate(block.Bytes)
	slog.Info("client.certificate common name: " + clientCertificate.Subject.CommonName)
	slog.Info("client.certificate email: " + clientCertificate.EmailAddresses[0])

	// get client private key
	aaa, _ := os.ReadFile("my-private_key.pem")
	bbb, _ := pem.Decode(aaa)
	privateKey, _ := x509.ParseECPrivateKey(bbb.Bytes)

	tlsCertificate := tls.Certificate{
		Certificate: [][]byte{clientCertificate.Raw},
		PrivateKey:  privateKey,
		Leaf:        clientCertificate,
	}

	// client TLS config
	config := &tls.Config{
		//GetClientCertificate: clientCertificate,
		Certificates: []tls.Certificate{tlsCertificate},
	}
	creds := credentials.NewTLS(config)
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		slog.Error("did not connect: " + err.Error())
		return
	}
	defer conn.Close()

	c := aaa_cz.NewGreeterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &aaa_cz.HelloRequest{Name: *name})
	if err != nil {
		slog.Error("could not greet: " + err.Error())
		return
	}

	slog.Info("Greeting: " + r.GetMessage())
}
