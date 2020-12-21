package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/fahmi1597/microservices-go/data"
	"github.com/fahmi1597/microservices-go/handlers"
	"github.com/gorilla/mux"
)

func main() {

	// Logging
	l := log.New(os.Stdout, "product-api: ", log.LstdFlags)
	v := data.NewValidator()

	// The handlers
	ph := handlers.NewProduct(l, v)

	// Create servemux
	sm := mux.NewRouter()

	// Create and register the handlers
	// GET
	getProduct := sm.Methods(http.MethodGet).Subrouter()
	getProduct.HandleFunc("/products", ph.GetListProduct)
	getProduct.HandleFunc("/products/{id:[0-9]+}", ph.GetProduct)
	// POST
	addProduct := sm.Methods(http.MethodPost).Subrouter()
	addProduct.HandleFunc("/products", ph.AddProduct)
	addProduct.Use(ph.MiddlewareValidation)
	// PUT
	updateProduct := sm.Methods(http.MethodPut).Subrouter()
	updateProduct.HandleFunc("/products", ph.UpdateProduct)
	updateProduct.Use(ph.MiddlewareValidation)
	// DELETE *not working yet*
	deleteProduct := sm.Methods(http.MethodDelete).Subrouter()
	deleteProduct.HandleFunc("/products/{id:[0-9]+}", ph.DeleteProduct)

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
		l.Println("Server starting on", s.Addr)
		err := s.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			l.Fatalf("[ERROR] starting server: %s\n", err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, os.Kill)

	sig := <-ch
	l.Printf("Signal %s received, shutting down", sig)

	// Graceful shutdown the server, waiting for max of 30 seconds until current operations is completed
	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer func() {
		cancel()
	}()

	err := s.Shutdown(tc)
	if err != nil {
		l.Fatalf("Shutdown failed: %+v", err)
	}
	os.Exit(0)
}
