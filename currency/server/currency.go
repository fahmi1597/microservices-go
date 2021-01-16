package server

import (
	"context"

	"github.com/fahmi1597/microservices-go/currency/data"

	// alias for protoc generated code
	protogc "github.com/fahmi1597/microservices-go/currency/protos/currency"
	"github.com/hashicorp/go-hclog"
)

// CurrencyServer is a gRPC server, it implements the method defined in CurrencyServer interface
type CurrencyServer struct {
	log   hclog.Logger
	rates *data.ExchangeRates
	protogc.UnimplementedCurrencyServer
}

// NewCurrencyServer creates a new Currency server
func NewCurrencyServer(l hclog.Logger, r *data.ExchangeRates) *CurrencyServer {
	return &CurrencyServer{
		log:   l,
		rates: r,
	}
}

// GetRate implements CurrencyServer GetRate method from the protoc generated code
func (cs *CurrencyServer) GetRate(ctx context.Context, rr *protogc.RateRequest) (*protogc.RateResponse, error) {
	cs.log.Debug(
		"GetRate",
		"base", rr.GetBase(),
		"destination",
		rr.GetDestination(),
	)

	rateRatio, err := cs.rates.GetRateRatio(rr.GetBase().String(), rr.GetDestination().String())
	if err != nil {
		return nil, err
	}
	// spew.Dump(rateRatio)
	return &protogc.RateResponse{Rate: rateRatio}, nil
}
