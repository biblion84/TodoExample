package main

import (
	"crypto/tls"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

var templateCache map[string]*template.Template

type application struct {
	templateCache map[string]*template.Template
	todos         []Todo
}

func main() {

	app := application{}
	app.todos = []Todo{{
		Checked: false,
		Text:    "Buy milk",
	}, {
		Checked: true,
		Text:    "Buy eggs",
	}}

	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	logger.Println("Server is starting...")
	var err error

	app.templateCache, err = newTemplateCache("./ui/html/")
	if err != nil {
		log.Fatal(err)
	}

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	//go http.ListenAndServe(":80", redirect())

	server := &http.Server{
		Addr:         ":443",
		Handler:      app.routes(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 120 * time.Second, // Long timeout to handle the challenge taking ages
		IdleTimeout:  time.Minute,
		TLSConfig:    tlsConfig,
	}

	logger.Println("Server is ready to handle requests at", "https://localhost")

	certFile := "./tls/cert.pem"
	keyFile := "./tls/key.pem"

	if err := server.ListenAndServeTLS(certFile, keyFile); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Could not listen on %s: %v\n", ":443", err)
	}

	logger.Println("Server stopped")
}
