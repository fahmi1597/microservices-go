package data

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-hclog"
)

type (
	// ExchangeRates receives the exchange rates that has been processed from RateResponse
	ExchangeRates struct {
		log   hclog.Logger
		rates map[string]float64
	}

	// RateReponse holds exchange rates data from API Call
	RateReponse struct {
		Rates map[string]float64 `json:"rates"`
	}
)

// NewExchangeRates creates a new instance of ExchangeRates
func NewExchangeRates(l hclog.Logger) (*ExchangeRates, error) {
	er := &ExchangeRates{
		log:   l,
		rates: map[string]float64{},
	}

	err := er.fetchRates()

	return er, err
}

// GetRateRatio calculates the exchange rate ratio between two currencies, base and destination
func (er *ExchangeRates) GetRateRatio(base, dest string) (float64, error) {
	// {"USD":1.0}

	// Check if the currency rate is available on our database.
	br, IsExist := er.rates[base]
	if !IsExist {
		return 0, fmt.Errorf("Unknown rate: %s", base)
	}
	dr, isExist := er.rates[dest]
	if !isExist {
		return 0, fmt.Errorf("Unknown rate: %s", dest)
	}

	return br / dr, nil

}
func (er *ExchangeRates) fetchRates() error {
	resp, err := http.Get("https://api.exchangeratesapi.io/latest")

	if err != nil {
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		er.log.Error("[ERROR] Rate request failed", "status", resp.StatusCode)
		return fmt.Errorf("Failed getting exchange rates")
	}
	defer resp.Body.Close()
	// {
	// 	"rates": {
	// 	  "CAD": 1.5543,
	// 	  "HKD": 9.4982,
	// 	}
	//  "unused": "asd"
	// }

	rr := RateReponse{}
	json.NewDecoder(resp.Body).Decode(&rr)

	for k, v := range rr.Rates {
		er.rates[k] = v
	}

	er.rates["USD"] = 1

	return nil
}
