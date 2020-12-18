package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/fahmi1597/microservices-go/data"
)

// MiddlewareValidation is a midleware handler for validation
func (p *Products) MiddlewareValidation(next http.Handler) http.Handler {

	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {

		prod := &data.Product{}

		err := data.FromJSON(prod, req.Body)
		if err != nil {
			p.l.Println("[ERROR] Deserialize product", http.StatusBadRequest)
			http.Error(resp, "Error reading product", http.StatusBadRequest)
			return
		}

		// Validate products
		vErr := p.v.Validate(prod)
		if len(vErr) != 0 {
			p.l.Println("[ERROR] Validating product", http.StatusBadRequest)
			http.Error(
				resp,
				fmt.Sprintf("Error validating product %s", vErr),
				http.StatusBadRequest,
			)
			return
		}

		// Create context of validated product
		ctx := context.WithValue(req.Context(), KeyProduct{}, *prod)

		// Send the validated product request to the next handler
		req = req.WithContext(ctx)
		next.ServeHTTP(resp, req)
	})
}
