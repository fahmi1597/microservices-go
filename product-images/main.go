package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/fahmi1597/microservices-go/product-images/files"
	"github.com/fahmi1597/microservices-go/product-images/handlers"
	ghandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
)

// UploadPath is path to store uploaded files
func main() {

	log := hclog.New(&hclog.LoggerOptions{
		Name:  "product-images",
		Level: hclog.LevelFromString("DEBUG"),
		Color: hclog.ForceColor,
	})

	// Initialize file storage
	fs, err := files.NewFileStorage("./uploads", 1024*1000*5)
	if err != nil {
		log.Error("Unable to create file storage", "error", err)

	}
	// Create server mux, handlers, middleware and CORS
	sm := mux.NewRouter()
	cors := ghandlers.CORS(ghandlers.AllowedOrigins([]string{"*"}))
	fh := handlers.NewFileHandler(log, fs)
	mw := handlers.GzipEncoding{}

	// Register the handlers

	// Upload file handlers
	ufh := sm.Methods(http.MethodPost).Subrouter()
	// Handler in REST way
	ufh.HandleFunc("/images/{id:[0-9]+}/{filename:[a-zA-Z.]+(?:jpg|jpeg|png|gif)$}", fh.UploadREST)
	// Handler in Multipart way
	ufh.HandleFunc("/", fh.UploadMultipart)

	// Serve files handlers
	sfh := sm.Methods(http.MethodGet).Subrouter()
	sfh.Handle(
		"/images/{id:[0-9]+}/{filename:[a-zA-Z.]+(?:jpg|jpeg|png|gif)$}",
		http.StripPrefix("/images/", http.FileServer(http.Dir("./uploads"))),
	)
	sfh.Use(mw.GzipMiddleware)

	s := &http.Server{
		Addr:         "localhost:9001",                                   // Listen Address
		Handler:      cors(sm),                                           // Default handler
		ErrorLog:     log.StandardLogger(&hclog.StandardLoggerOptions{}), // Set the logger for the server
		ReadTimeout:  time.Second * 15,                                   // Max time duration to read request
		WriteTimeout: time.Second * 15,                                   // Max time duration to write response
		IdleTimeout:  time.Second * 60,                                   // Max time duration to keep connetion alive
	}

	go func() {
		log.Info("Starting server on", "address", s.Addr)
		err := s.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Error("Unable to start server", "error", err)
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	sig := <-c
	log.Info("Server shutting down", "signal", sig)

	// Graceful shutdown the server, waiting for max of 30 seconds until current operations is completed
	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = s.Shutdown(tc)
	if err != nil {
		log.Error("Unable to shutdown gracefully", "error", err)
		os.Exit(1)
	}
	os.Exit(0)
}
