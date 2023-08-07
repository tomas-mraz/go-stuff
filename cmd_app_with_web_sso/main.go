package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
	"log"
	"net/http"
)

var (
	state         = "random"
	ctx           context.Context
	verifier      *oidc.IDTokenVerifier
	oauth2Config  oauth2.Config
	receivedToken chan struct{}
)

func handleRequest(responseWriter http.ResponseWriter, req *http.Request) {
	log.Printf("Accept request %s to %s\n", req.Method, req.URL)

	if req.URL.Query().Get("state") != state {
		log.Println("State is", state)
		http.Error(responseWriter, "state did not match", http.StatusBadRequest)
		return
	}

	oauth2Token, err := oauth2Config.Exchange(ctx, req.URL.Query().Get("code"))
	if err != nil {
		http.Error(responseWriter, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		http.Error(responseWriter, "No id_token field in oauth2 token.", http.StatusInternalServerError)
		return
	}
	log.Println("rawIDToken", rawIDToken)

	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		http.Error(responseWriter, "Failed to verify ID Token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := struct {
		OAuth2Token   *oauth2.Token
		IDTokenClaims *json.RawMessage // ID Token payload is just JSON.
	}{oauth2Token, new(json.RawMessage)}

	err = idToken.Claims(&resp.IDTokenClaims)
	data, _ := json.MarshalIndent(resp, "", "    ")
	log.Println(string(data[:]))

	_, err = responseWriter.Write([]byte("mam token"))
	if err != nil {
		panic("error #3")
	}

	accessToken := resp.OAuth2Token.AccessToken
	log.Println(accessToken)
	close(receivedToken)
}

func main() {
	log.Println("Starting")
	configURL := "https://login.cubyte.online:8443/realms/cubyte"
	ctx = context.Background()
	provider, err := oidc.NewProvider(ctx, configURL)
	if err != nil {
		panic(err)
	}
	clientID := "demo-client"
	clientSecret := "3u8UoWHe5N6zjO5WYbFk1pSdgHLPfp5g"
	listenOn := "127.0.0.1:8432"
	redirectURL := "http://" + listenOn + "/callback"

	// Configure an OpenID Connect aware OAuth2 client.
	oauth2Config = oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		// Discovery returns the OAuth2 endpoints.
		Endpoint: provider.Endpoint(),
		// "openid" is a required scope for OpenID Connect flows.
		Scopes: []string{oidc.ScopeOpenID, "profile", "email", "roles"},
	}
	oidcConfig := &oidc.Config{
		ClientID: clientID,
	}
	verifier = provider.Verifier(oidcConfig)

	log.Printf("listening on http://%s/", listenOn)
	server := &http.Server{
		Addr: listenOn,
	}
	handler := http.HandlerFunc(handleRequest)
	http.Handle("/callback", handler)
	receivedToken = make(chan struct{})
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				log.Println("server closed")
			} else {
				log.Println("error #1", err)
			}
		}
	}()

	log.Println("open browser")
	url := oauth2Config.AuthCodeURL(state)
	fmt.Println(url)
	//TODO open browser

	log.Println("waiting for token")
	<-receivedToken

	log.Println("shut local server")
	err = server.Close()
	if err != nil {
		panic("error #2")
	}

	log.Println("done")
}
