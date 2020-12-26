package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/fahmi1597/microservices-go/product-images/files"
	"github.com/fahmi1597/microservices-go/product-images/handlers"
	ghandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// UploadPath is path to store uploaded files
func main() {

	l := log.New(os.Stdout, "product-images: ", log.LstdFlags)

	// Initialize file storage
	fs, err := files.New("./uploads", 1024*1000*5)
	if err != nil {
		l.Fatal("[ERROR] Unable to create storage", err)

	}
	// Create server mux, handlers, and CORS
	sm := mux.NewRouter()
	cors := ghandlers.CORS(ghandlers.AllowedOrigins([]string{"*"}))
	fh := handlers.New(l, fs)

	// Register the handlers

	// Upload file handlers
	ufh := sm.Methods(http.MethodPut).Subrouter()
	ufh.HandleFunc("/images/{id:[0-9]+}/{filename:[a-zA-Z.]+(?:jpg|jpeg|png|gif)$}", fh.Upload)
	//fhPUT.HandleFunc("/images/{id:[0-9]+}/{filename:[a-zA-Z.]+[a-z]{3,4}}", fh.ServeHTTP)

	// Serve files handlers
	sfh := sm.Methods(http.MethodGet).Subrouter()
	sfh.Handle(
		"/images/{id:[0-9]+}/{filename:[a-zA-Z.]+(?:jpg|jpeg|png|gif)$}",
		http.StripPrefix("/images/", http.FileServer(http.Dir("./uploads"))),
	)

	s := &http.Server{
		Addr:         ":3030",          // Listen Address
		Handler:      cors(sm),         // Default handler
		ReadTimeout:  time.Second * 15, // Max time duration to read request
		WriteTimeout: time.Second * 15, // Max time duration to write response
		IdleTimeout:  time.Second * 60, // Max time duration to keep connetion alive
	}

	go func() {
		l.Println("[INFO] Starting server on", s.Addr)
		err := s.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			l.Fatalf("[ERROR] Starting server: %s\n", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	sig := <-c
	l.Printf("Signal %s received, shutting down", sig)

	// Graceful shutdown the server, waiting for max of 30 seconds until current operations is completed
	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer func() {
		cancel()
	}()

	err = s.Shutdown(tc)
	if err != nil {
		l.Fatalf("Shutdown failed: %+v", err)
	}
	os.Exit(0)
}