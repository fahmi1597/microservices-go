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

	// Loging
	l := log.New(os.Stdout, "product-api : ", log.LstdFlags)

	// Handlers
	hh := handlers.NewHello(l)
	gh := handlers.NewGoodbye(l)

	// ServeMux
	sm := http.NewServeMux()
	sm.Handle("/", hh)
	sm.Handle("/goodbye", gh)

	s := &http.Server{
		Addr:         ":3000",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil {
			l.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Application terminated, performing graceful shutdown", "[ Signal :", sig, "]")

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)

	// log.Fatal(http.ListenAndServe(":3000", sm))
}
