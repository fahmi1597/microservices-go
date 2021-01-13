package server

import (
	"context"

	rate "github.com/fahmi1597/microservices-go/currency/data"

	// alias for protoc generated code
	protogc "github.com/fahmi1597/microservices-go/currency/protos/currency"
	"github.com/hashicorp/go-hclog"
)

// Currency is a gRPC server, it implements the method defined in CurrencyServer interface
type Currency struct {
	log   hclog.Logger
	rates *rate.ExchangeRates
	protogc.UnimplementedCurrencyServer
}

// NewCurrencyServer creates a new Currency server
func NewCurrencyServer(l hclog.Logger, r *rate.ExchangeRates) *Currency {
	return &Currency{log: l}
}

// GetRate implements CurrencyServer GetRate method from the protoc generated code
func (cs *Currency) GetRate(ctx context.Context, rr *protogc.RateRequest) (*protogc.RateResponse, error) {
	cs.log.Info(
		"Handle GetRate",
		"base", rr.GetBase(),
		"destination",
		rr.GetDestination(),
	)

	rateRatio, err := cs.rates.GetRateRatio(rr.GetBase().String(), rr.GetDestination().String())
	if err != nil {
		return nil, err
	}
	return &protogc.RateResponse{Rate: rateRatio}, nil
}
