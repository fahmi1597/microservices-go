package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	protogc "github.com/fahmi1597/microservices-go/currency/protos/currency"
	"github.com/fahmi1597/microservices-go/product-api/data"
	"github.com/fahmi1597/microservices-go/product-api/handlers"
	"github.com/go-openapi/runtime/middleware"
	ghandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
)

func main() {

	// Create a new instance of logger and validator
	// log := log.New(os.Stdout, "product-api: ", log.LstdFlags)
	log := hclog.New(&hclog.LoggerOptions{
		Name:  "product-api",
		Level: hclog.LevelFromString("DEBUG"),
	})

	validator := data.NewValidator()

	// Create grpc connection
	grpcConn, err := grpc.Dial("localhost:9002", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer grpcConn.Close()

	// New instance of grpc currency client
	cc := protogc.NewCurrencyClient(grpcConn)

	// New instance of ProductsDB
	pdb := data.NewProductsDB(log, cc)

	// New instance of Products handler
	ph := handlers.NewProductHandler(log, validator, pdb)

	// Create ServeMux
	sm := mux.NewRouter()

	// Create and register the handlers
	routerGet := sm.Methods(http.MethodGet).Subrouter()
	routerGet.HandleFunc("/products", ph.GetProductList).Queries("currency", "{[A-Z]{3}}")
	routerGet.HandleFunc("/products", ph.GetProductList)
	routerGet.HandleFunc("/products/{id:[0-9]+}", ph.GetProductByID)

	routerPost := sm.Methods(http.MethodPost).Subrouter()
	routerPost.HandleFunc("/products", ph.AddProduct)
	routerPost.Use(ph.MiddlewareValidation)

	routerPut := sm.Methods(http.MethodPut).Subrouter()
	routerPut.HandleFunc("/products", ph.UpdateProduct)
	routerPut.Use(ph.MiddlewareValidation)

	routerDel := sm.Methods(http.MethodDelete).Subrouter()
	routerDel.HandleFunc("/products/{id:[0-9]+}", ph.DeleteProduct)

	// Handler for swagger documentation
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)

	routerGet.Handle("/docs", sh)
	routerGet.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	// Allow cross-origin resource sharing
	ch := ghandlers.CORS(ghandlers.AllowedOrigins([]string{"*"}))

	// HTTP Server with custom configuration
	s := &http.Server{
		Addr:         "localhost:9000",                                   // bind address
		Handler:      ch(sm),                                             // default handler wrapped with
		ErrorLog:     log.StandardLogger(&hclog.StandardLoggerOptions{}), // set the logger for the server
		ReadTimeout:  time.Second * 15,                                   // max time to read request from client
		WriteTimeout: time.Second * 15,                                   // max time to write request to client
		IdleTimeout:  time.Second * 120,                                  // max time for connection using TCP Keep-Alive
	}

	// Start the server inside go routines
	go func() {
		log.Info("Server starting on", "address", s.Addr)
		err := s.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Error("Failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	sig := <-c
	log.Info("Shutting down", "signal", sig)

	// Graceful shutdown the server, waiting for max of 30 seconds until current operations is completed
	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = s.Shutdown(tc)
	if err != nil {
		log.Error("Shutdown failed", "error", err)
	}
	os.Exit(0)
}
