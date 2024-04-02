package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"golang.org/x/crypto/acme/autocert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"grpc3/aaa_cz"
	"log"
	"net"
	"net/http"
	"os"
)

var (
	port = flag.Int("port", 50051, "The server port")
	addr = flag.String("addr", "[2001:470:6f:53b::1710]", "The server address")
)

const (
	servername = "backend.cubyte.space"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	aaa_cz.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *aaa_cz.HelloRequest) (*aaa_cz.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &aaa_cz.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func loadCaPool() *x509.CertPool {
	// Load certificate of the CA who signed server's certificate
	pemServerCA, err := os.ReadFile("certs/ca-cert.pem")
	if err != nil {
		fmt.Println("failed to read file with CA certificate")
		panic("aaa")
	}
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		fmt.Println("failed to add server CA's certificate")
		panic("bbb")
	}
	return certPool
}

func main() {
	flag.Parse()

	bind := fmt.Sprintf("%s:%d", *addr, *port)
	log.Println(bind)
	lis, err := net.Listen("tcp6", bind)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Let's Encrypt
	certManager := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(servername),
		Cache:      autocert.DirCache("certs"),
	}
	tlsConfig := certManager.TLSConfig()
	tlsConfig.MinVersion = tls.VersionTLS12
	tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
	tlsConfig.ClientCAs = loadCaPool()

	tlsCredentials := credentials.NewTLS(tlsConfig)

	s := grpc.NewServer(grpc.Creds(tlsCredentials))
	aaa_cz.RegisterGreeterServer(s, &server{})

	// is needed for ACME challenge
	go func() {
		err := http.ListenAndServe(*addr+":http", certManager.HTTPHandler(nil))
		log.Println("ukoncen acme listener")
		if err != nil {
			log.Println(err)
		}
	}()

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
