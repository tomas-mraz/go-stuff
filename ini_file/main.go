package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"golang.org/x/crypto/acme/autocert"
	"gopkg.in/ini.v1"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"releaseregister/api"
	"releaseregister/database"
	"releaseregister/global"
	"releaseregister/web"
	"time"
)

func MiddlewareOne(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log.Println("middleware one")
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func MiddlewareTwo(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log.Println("middleware two")
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func webPagesRouter() api.Router {
	controller := &web.WebController{}
	return controller
}

func getStoredOrLetsEncryptCert(certManager *autocert.Manager) func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	return func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		dirCache, ok := certManager.Cache.(autocert.DirCache)
		if !ok {
			dirCache = "certs"
		}

		keyFile := filepath.Join(string(dirCache), hello.ServerName+".key")
		crtFile := filepath.Join(string(dirCache), hello.ServerName+".crt")
		certificate, err := tls.LoadX509KeyPair(crtFile, keyFile)
		if err != nil {
			log.Printf("%s\nUsing Letsencrypt certificate\n", err)
			return certManager.GetCertificate(hello)
		}
		log.Println("Loaded stored certificate.")
		return &certificate, err
	}
}

func main() {
	log.Println("Starting...")

	cfg, err := ini.Load("releaseregister.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	log.Println("Configuration load")

	global.Host = cfg.Section("server").Key("host").String()
	global.Port = cfg.Section("server").Key("port").String()
	global.Protocol = cfg.Section("server").Key("protocol").String()
	global.Domain = cfg.Section("server").Key("domain").String()
	global.DataDirectory = cfg.Section("paths").Key("data").String()

	global.Init()

	//TODO udělat proměnnou
	log.Println("- app mode:", cfg.Section("").Key("app_mode").String())
	log.Println("- data path:", global.DataDirectory)
	log.Println("- mask:", global.DomainMask)
	log.Println("- url:", global.Url)
	log.Println("- listening:", global.Host)

	database.Configure(cfg)
	database.Connect()
	log.Println("Database connect")

	/*
		router := mux.NewRouter()

		apiRouter := router.PathPrefix("/api/v1").Subrouter()
		//apiRouter.Use(MiddlewareOne)
		apiRouter.HandleFunc("/", handler)

		webRouter := router.PathPrefix("/").Subrouter()
		//webRouter.Use(MiddlewareTwo)
		webRouter.HandleFunc("/", handler2)

		log.Fatal(http.ListenAndServe(":8080", router))
	*/
	apiService := api.NewArtifactApiService()
	apiRouter := api.NewArtifactApiController(apiService)

	webRouter := webPagesRouter()

	router := api.NewRouter(apiRouter, webRouter)

	// better shutdown
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()
/*
	// TLS
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(global.Domain, "simple.releaseregister.com", "test.releaseregister.com"),
		Cache:      autocert.DirCache("certs"),
	}
	tlsConfig := certManager.TLSConfig()
	//tlsConfig.MinVersion = tls.VersionTLS12
	//tlsConfig.CurvePreferences = []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256}
	tlsConfig.GetCertificate = getStoredOrLetsEncryptCert(&certManager)
*/
	srv := &http.Server{
		Addr:         global.Host + ":" + global.Port,
		Handler:      router,           // Pass our instance of gorilla/mux in.
		WriteTimeout: time.Second * 15, // Good practice to set timeouts to avoid Slowloris attacks.
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	// Run our server in a goroutine so that it doesn't block.
	switch {
	/*	case global.Protocol == "https":
		go http.ListenAndServe(global.Host+":80", certManager.HTTPHandler(nil))
		go func() {
			if err := srv.ListenAndServeTLS("", ""); err != nil {
				log.Println(err)
			}
		}()
	*/
	case global.Protocol == "http":
		go func() {
			if err := srv.ListenAndServe(); err != nil {
				log.Println(err)
			}
		}()
	default:
		panic("Not supported protocol: " + global.Protocol)
	}
	log.Println("Webserver start")

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}
