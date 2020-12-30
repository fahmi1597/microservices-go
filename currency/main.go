package main

import (
	"net"
	"os"

	protogc "github.com/fahmi1597/microservices-go/currency/protos/currency"
	"github.com/fahmi1597/microservices-go/currency/server"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log := hclog.Default()

	grpcServer := grpc.NewServer()
	currencyServer := server.NewCurrencyServer(log)

	protogc.RegisterCurrencyServer(grpcServer, currencyServer)

	// don't use in production
	reflection.Register(grpcServer)

	nl, err := net.Listen("tcp", "localhost:3000")
	if err != nil {
		log.Error("Error", err)
		os.Exit(1)
	}

	grpcServer.Serve(nl)
}
