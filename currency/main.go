package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	rate "github.com/fahmi1597/microservices-go/currency/data"
	protogc "github.com/fahmi1597/microservices-go/currency/protos/currency"
	"github.com/fahmi1597/microservices-go/currency/server"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log := hclog.New(&hclog.LoggerOptions{
		Name:  "grpc currency",
		Level: hclog.LevelFromString("DEBUG"),
	})

	// Create a new instance of gRPC server
	grpcServer := grpc.NewServer()

	// Generate currency exchange rates
	currencyRates, err := rate.NewExchangeRates(log)
	if err != nil {
		log.Error("Unable to generate currency exchange rates")
		os.Exit(1)
	}

	// Create new instance of Currency server
	currencyServer := server.NewCurrencyServer(log, currencyRates)

	// Register the Currency server.
	protogc.RegisterCurrencyServer(grpcServer, currencyServer)

	// Don't use in production
	reflection.Register(grpcServer)

	nl, err := net.Listen("tcp", "localhost:9002")
	if err != nil {
		log.Error("Error", err)
		os.Exit(1)
	}

	go func() {
		log.Info("Server starting", "address", nl.Addr().String())
		err := grpcServer.Serve(nl)
		if err != grpc.ErrServerStopped {
			log.Error("Unclean shutdown", err)
			os.Exit(1)
		}

	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	sig := <-c
	log.Warn("Shutting down", "signal", sig)
	grpcServer.GracefulStop()

}
