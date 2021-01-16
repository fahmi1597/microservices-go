package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/fahmi1597/microservices-go/currency/data"
	protogc "github.com/fahmi1597/microservices-go/currency/protos/currency"
	"github.com/fahmi1597/microservices-go/currency/server"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	l := hclog.New(&hclog.LoggerOptions{
		Name:  "grpc currency",
		Level: hclog.LevelFromString("DEBUG"),
	})

	// Create a new instance of gRPC Server
	gs := grpc.NewServer()

	er, err := data.NewExchangeRates(l)
	if err != nil {
		l.Error("Unable to create currency exchange rates")
		os.Exit(1)
	}

	// Create new instance of CurrencyServer
	cs := server.NewCurrencyServer(l, er)

	// Register the gRPC currency server
	protogc.RegisterCurrencyServer(gs, cs)

	// Don't use in production
	reflection.Register(gs)

	nl, err := net.Listen("tcp", "localhost:9002")
	if err != nil {
		l.Error("Error", err)
		os.Exit(1)
	}

	go func() {
		l.Info("Server starting", "address", nl.Addr().String())
		err := gs.Serve(nl)
		if err != grpc.ErrServerStopped {
			l.Error("Unclean shutdown", err)
			os.Exit(1)
		}

	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	sig := <-c
	l.Info("Shutting down", "signal", sig)
	gs.GracefulStop()

}
