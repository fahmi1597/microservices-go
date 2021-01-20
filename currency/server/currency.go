package server

import (
	"context"
	"io"
	"time"

	"github.com/fahmi1597/microservices-go/currency/data"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	// alias for protoc generated code
	protogc "github.com/fahmi1597/microservices-go/currency/protos/currency"
	"github.com/hashicorp/go-hclog"
)

// CurrencyServer is a gRPC server, it implements the method defined in CurrencyServer interface
type CurrencyServer struct {
	// gctx          context.Context
	log   hclog.Logger
	rates *data.ExchangeRates
	// rateSubscription is used for cache
	// It's a map that returns a collection of client rate requests
	subscriptions map[protogc.Currency_SubscribeRatesServer][]*protogc.RateRequest
	protogc.UnimplementedCurrencyServer
}

// NewCurrencyServer creates a new Currency server
func NewCurrencyServer(l hclog.Logger, r *data.ExchangeRates) *CurrencyServer {
	cs := &CurrencyServer{
		// gctx:          g,
		log:           l,
		rates:         r,
		subscriptions: make(map[protogc.Currency_SubscribeRatesServer][]*protogc.RateRequest),
	}

	// monitor for updates
	go cs.getRateUpdates()

	return cs
}

func (cs *CurrencyServer) getRateUpdates() {
	rateUpdates := cs.rates.MonitorRates(time.Second * 5)

	for range rateUpdates {
		cs.log.Info("Got exchange rates update")

		// loop over subscription server to find subscribed client requests
		for srs, rr := range cs.subscriptions {

			// loop over client requests to get requested rate
			for _, r := range rr {
				rateRatio, err := cs.rates.GetRateRatio(r.Base.String(), r.Destination.String())
				if err != nil {
					cs.log.Error(
						"Unable to get rates update",
						"base", r.Base,
						"dest", r.Destination,
						"error", err,
					)
				}

				err = srs.Send(&protogc.RateResponse{
					Base:        r.Base,
					Destination: r.Destination,
					Rate:        rateRatio,
				})
				if err != nil {
					cs.log.Error(
						"Unable to send updated rates",
						"base", r.Base,
						"dest", r.Destination,
						"error", err,
					)
				}
			}
		}
	}
}

// GetRate implements CurrencyServer GetRate method from the protoc generated code
func (cs *CurrencyServer) GetRate(ctx context.Context, rr *protogc.RateRequest) (*protogc.RateResponse, error) {
	cs.log.Debug(
		"Handle get rate",
		"base", rr.GetBase(),
		"destination",
		rr.GetDestination(),
	)

	if rr.Destination == rr.Base {
		// example of enrichment error in unary rpc
		errStatus := status.Newf(
			codes.InvalidArgument,
			"The base currency %s can not be the same as the destination currency %s", //[0]
			rr.GetBase().String(),
			rr.GetDestination().String(),
		)

		errStatus, errWithDetails := errStatus.WithDetails(rr)
		if errWithDetails != nil {
			return nil, errWithDetails
		}

		// cs.log.Debug("Unary rich error", "error", errStatus.Code())
		return nil, errStatus.Err()
	}
	rateRatio, err := cs.rates.GetRateRatio(rr.GetBase().String(), rr.GetDestination().String())
	if err != nil {
		return nil, err
	}
	// spew.Dump(rateRatio)
	return &protogc.RateResponse{Rate: rateRatio}, nil
}

// SubscribeRates is new thing
func (cs *CurrencyServer) SubscribeRates(srs protogc.Currency_SubscribeRatesServer) error {

	// Read rate request from client standard input stream (stdin)
	for {
		rateRequest, err := srs.Recv()
		if err == io.EOF {
			cs.log.Info("Client has closed connection")

			// then delete the client from subscription
			delete(cs.subscriptions, srs)
			break
		}
		if err != nil {
			cs.log.Error("Unable to read client request", "error", err)

			// same for any kind of error, delete the client
			delete(cs.subscriptions, srs)
			return err
		}
		cs.log.Info("Handle client request", "base", rateRequest.Base, "dest", rateRequest.Destination)

		// If no client subscription in server, create initial rate requests
		rateRequests, ok := cs.subscriptions[srs]
		if !ok {
			rateRequests = []*protogc.RateRequest{}
		}

		// Otherwise cache the client subscription
		rateRequests = append(rateRequests, rateRequest)
		cs.subscriptions[srs] = rateRequests
	}

	return nil
}
