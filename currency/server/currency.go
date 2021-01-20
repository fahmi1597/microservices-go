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
	// subscriptions is a map that returns a collection of client rate requests.
	// It is used for cache.
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

// GetRate implements CurrencyServer GetRate method from the protoc generated code
func (cs *CurrencyServer) GetRate(ctx context.Context, rr *protogc.RateRequest) (*protogc.RateResponse, error) {
	cs.log.Debug("Handle get rate", "base", rr.GetBase(), "destination", rr.GetDestination())

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

	// Read client rate request from stream
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
		cs.log.Info("Added subscription for", "base", rateRequest.Base.String(), "dest", rateRequest.Destination.String())

		// Check if the client rate request exist in our subscription
		rateRequests, ok := cs.subscriptions[srs]
		if !ok {
			// If it doesn't, create one.
			rateRequests = []*protogc.RateRequest{}
		}

		// Check for client subscription rate currency as it cannot be the same.
		for _, rrs := range rateRequests {
			if rrs.Base == rateRequest.Base && rrs.Destination == rateRequest.Destination {
				cs.log.Error("Subscription already active for", "base", rateRequest.Base.String(), "destination", rateRequest.Destination.String())

				// If it was, then create a custom grpc error message and send it to the stream,
				// instead returning an error as that will terminate the client connection.
				grpcError := status.Newf(codes.AlreadyExists, "Subscription for the rate already active")
				grpcError, err = grpcError.WithDetails(rateRequest)
				if err != nil {
					cs.log.Error("Unable to create add error metadata", "error", err)
					continue
				}

				// The error later can be handled by Recv().
				streamError := &protogc.ErrorOrRateResponse_ErrorMessage{ErrorMessage: grpcError.Proto()}
				srs.Send(&protogc.ErrorOrRateResponse{Message: streamError})
				// exit the loop to prevent error stacking.
				break
			}
		}

		// If all is ok, append the client subscription to our collection.
		rateRequests = append(rateRequests, rateRequest)
		// then update the subscription list
		cs.subscriptions[srs] = rateRequests
	}

	return nil
}

func (cs *CurrencyServer) getRateUpdates() {

	// Simulated rate update every 15 sec, real data update is once a day
	rateUpdates := cs.rates.MonitorRates(time.Second * 60)

	for range rateUpdates {
		cs.log.Info("Got exchange rates update")
		// loop over subscription server to find subscribed client requests
		for srs, rr := range cs.subscriptions {
			// loop over client requests to get requested rate
			for _, r := range rr {
				rateRatio, err := cs.rates.GetRateRatio(r.Base.String(), r.Destination.String())
				if err != nil {
					cs.log.Error("Unable to get rates update", "error", err)
				}

				// streamUpdates contains updated rate
				streamUpdates := &protogc.ErrorOrRateResponse_RateResponse{
					RateResponse: &protogc.RateResponse{
						Base:        r.Base,
						Destination: r.Destination,
						Rate:        rateRatio,
					},
				}

				err = srs.Send(&protogc.ErrorOrRateResponse{Message: streamUpdates})
				if err != nil {
					cs.log.Error("Unable to send updated rates", "base", r.GetBase().String(), "dest", r.GetDestination().String(), "error", err)
				}
			}
		}
	}
}
