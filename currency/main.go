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

	er, err := data.NewExchangeRates(l)
	if err != nil {
		l.Error("Unable to create currency exchange rates")
		os.Exit(1)
	}

	// pctx := context.Background()
	// globalCtx, globalCancel := context.WithCancel(pctx)

	// Create new instance of CurrencyServer
	cs := server.NewCurrencyServer(l, er)

	// Create a new instance of gRPC Server
	gs := grpc.NewServer()

	// Register the gRPC currency server
	protogc.RegisterCurrencyServer(gs, cs)

	// Don't use in production
	reflection.Register(gs)

	// Graceful
	// taken from https://stackoverflow.com/questions/55797865/behavior-of-server-gracefulstop-in-golang
	c := make(chan os.Signal, 1)
	cErr := make(chan error, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	nl, err := net.Listen("tcp", "localhost:9002")
	if err != nil {
		l.Error("Error", err)
		os.Exit(1)
	}

	go func() {
		l.Info("Server starting", "address", nl.Addr().String())
		if err := gs.Serve(nl); err != nil {
			cErr <- err
		}
	}()
	defer func() {
		gs.GracefulStop()
	}()

	select {
	case err := <-cErr:
		l.Error("Fatal error", "error", err)
		os.Exit(1)
	case sig := <-c:
		l.Info("Shutting down", "signal", sig)
	}

}
