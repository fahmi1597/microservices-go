package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/fahmi1597/microservices-go/product-api/data"
	"github.com/fahmi1597/microservices-go/product-api/handlers"
	"github.com/go-openapi/runtime/middleware"
	ghandlers "github.com/gorilla/handlers"
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

	// Create and register the handlers for each subrouter
	// GET
	getR := sm.Methods(http.MethodGet).Subrouter()
	getR.HandleFunc("/products", ph.GetProducts)
	getR.HandleFunc("/products/{id:[0-9]+}", ph.GetProduct)
	// POST
	postR := sm.Methods(http.MethodPost).Subrouter()
	postR.HandleFunc("/products", ph.AddProduct)
	postR.Use(ph.MiddlewareValidation)
	// PUT
	putR := sm.Methods(http.MethodPut).Subrouter()
	putR.HandleFunc("/products", ph.UpdateProduct)
	putR.Use(ph.MiddlewareValidation)
	// DELETE *not working yet*
	delR := sm.Methods(http.MethodDelete).Subrouter()
	delR.HandleFunc("/products/{id:[0-9]+}", ph.DeleteProduct)

	// handler for documentation
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)

	getR.Handle("/docs", sh)
	getR.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	ch := ghandlers.CORS(ghandlers.AllowedOrigins([]string{"*"}))
	s := &http.Server{
		Addr:         "localhost:8000",  // bind address
		Handler:      ch(sm),            // default handler
		ErrorLog:     l,                 // set the logger for the server
		ReadTimeout:  time.Second * 15,  // max time to read request from client
		WriteTimeout: time.Second * 15,  // max time to write request to client
		IdleTimeout:  time.Second * 120, // max time for connection using TCP Keep-Alive
	}

	// Start the server inside go routines
	go func() {
		l.Println("Server starting on", s.Addr)
		err := s.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			l.Fatalf("[ERROR] starting server: %s\n", err)
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

	err := s.Shutdown(tc)
	if err != nil {
		l.Fatalf("Shutdown failed: %+v", err)
	}
	os.Exit(0)
}
