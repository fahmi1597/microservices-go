package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/fahmi1597/microservices-go/handlers"
	"github.com/gorilla/mux"
)

func main() {

	// Logging
	l := log.New(os.Stdout, "product-api: ", log.LstdFlags)

	// The handlers
	ph := handlers.NewProduct(l)

	// Create servemux
	sm := mux.NewRouter()

	// Create the handlers
	getProduct := sm.Methods(http.MethodGet).Subrouter()
	getProduct.HandleFunc("/", ph.GetProduct)

	addProduct := sm.Methods(http.MethodPost).Subrouter()
	addProduct.HandleFunc("/", ph.AddProduct)
	addProduct.Use(ph.MiddlewareValidation)

	updateProduct := sm.Methods(http.MethodPut).Subrouter()
	updateProduct.HandleFunc("/{id:[0-9]+}", ph.UpdateProduct)
	updateProduct.Use(ph.MiddlewareValidation)

	s := &http.Server{
		Addr:         "localhost:3000",  // bind address
		Handler:      sm,                // default handler
		ErrorLog:     l,                 // set the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from client
		WriteTimeout: 10 * time.Second,  // max time to write request to client
		IdleTimeout:  120 * time.Second, // max time for connection using TCP Keep-Alive
	}

	// Start the server inside go routines
	go func() {
		l.Println("Server starting on port", s.Addr)
		if err := s.ListenAndServe(); err != nil {
			l.Println(err)
			os.Exit(1)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, os.Kill)

	sig := <-ch
	l.Printf("Application terminated: %s", sig.String())

	// Graceful shutdown the server, waiting for max of 30 seconds until current operations is completed
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)

}
