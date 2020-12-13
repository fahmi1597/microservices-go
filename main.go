package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/fahmi1597/microservices-go/handlers"
)

func main() {

	// Logging
	l := log.New(os.Stdout, "product-api : ", log.LstdFlags)

	// The handlers
	ph := handlers.NewProduct(l)

	// Create servemux and register the handlers
	sm := http.NewServeMux()
	sm.Handle("/", ph)

	s := &http.Server{
		Addr:         ":3000",           // bind address
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
			l.Printf("Error starting server %s \n", err.Error())
			os.Exit(1)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, os.Kill)

	sig := <-ch
	l.Printf("Application terminated (reason : %s)", sig.String())

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)

}
