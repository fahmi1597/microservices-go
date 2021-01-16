package data

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-hclog"
)

func TestExchangeRates(t *testing.T) {
	testRate, err := NewExchangeRates(hclog.Default())

	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Rates: %#v", testRate.rates)
}

func TestExchangeRates_GetRateRatio(t *testing.T) {
	type fields struct {
		log   hclog.Logger
		rates map[string]float64
	}
	type args struct {
		base string
		dest string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    float64
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Case 1: Missing currency",
			fields: fields{
				log:   hclog.Default(),
				rates: map[string]float64{"EUR": 1, "JPY": 100},
			},
			args: args{
				base: "EUR",
				dest: "IDR",
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "Case 2: Available currency",
			fields: fields{
				log:   hclog.Default(),
				rates: map[string]float64{"EUR": 1, "JPY": 100},
			},
			args: args{
				base: "EUR",
				dest: "JPY",
			},
			want:    100,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			er := &ExchangeRates{
				log:   tt.fields.log,
				rates: tt.fields.rates,
			}
			got, err := er.GetRateRatio(tt.args.base, tt.args.dest)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExchangeRates.GetRateRatio() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ExchangeRates.GetRateRatio() = %v, want %v", got, tt.want)
			}
		})
	}
}
