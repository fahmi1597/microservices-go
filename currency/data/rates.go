package data

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-hclog"
)

// ExchangeRates receives the exchange rates that has been processed from RateResponse
type ExchangeRates struct {
	log   hclog.Logger
	rates map[string]float64
}

// RateReponse holds exchange rates data from API Call
type RateReponse struct {
	Rates map[string]float64 `json:"rates,omitempty"`
}

// NewExchangeRates creates a new instance of ExchangeRates
func NewExchangeRates(log hclog.Logger) (*ExchangeRates, error) {
	er := &ExchangeRates{
		log:   log,
		rates: map[string]float64{},
	}
	// to do later: env based exchange rate endpoint
	// err := er.fetchRates(endpoint)
	err := er.fetchRates()
	if err != nil {
		hclog.Default().Debug("Failed to fetch latest exchange rates", "error", err)
		return nil, err
	}

	return er, nil
}

// GetRateRatio calculates the rate ratio between two currencies, base and destination
func (er *ExchangeRates) GetRateRatio(base, dest string) (float64, error) {

	// Check if the currency rate is available on our database.
	baseRate, isExist := er.rates[base]
	if !isExist {
		return 0, fmt.Errorf("Unknown rate: %s", base)
	}
	destRate, isExist := er.rates[dest]
	if !isExist {
		return 0, fmt.Errorf("Unknown rate: %s", dest)
	}
	return destRate / baseRate, nil

}
func (er *ExchangeRates) fetchRates() error {

	// fetch latest exchange rate from api endpoint
	resp, err := http.Get("https://api.exchangeratesapi.io/latest")
	if err != nil {
		er.log.Error("Failed to reach exchange rates endpoint", "error", err)
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		er.log.Error("Failed to reach exchange rates endpoint", "error", resp.Status)
		return fmt.Errorf("Exchange rates server is unavailable")
	}
	defer resp.Body.Close()

	rateResponse := RateReponse{}
	err = json.NewDecoder(resp.Body).Decode(&rateResponse)
	if err != nil {
		return err
	}

	// fill ExchangeRates with the latest exchange rate from api call
	// {"curr" : rate} = {"eur": 1.0}
	for curr, rate := range rateResponse.Rates {
		er.rates[curr] = rate
	}

	er.rates["EUR"] = 1

	return nil
}
