package handlers

import (
	"context"
	"net/http"

	"github.com/fahmi1597/microservices-go/product-api/data"
)

// MiddlewareValidation is a middleware handler for product validation before going to the next handler
func (p *Products) MiddlewareValidation(next http.Handler) http.Handler {

	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {

		prod := &data.Product{}

		err := data.FromJSON(prod, req.Body)
		if err != nil {
			p.l.Println("[ERROR] Deserialize product", http.StatusBadRequest)

			resp.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&GenericError{Message: err.Error()}, resp)
			return
		}

		// Validate products
		vErr := p.v.Validate(prod)
		if len(vErr) != 0 {
			p.l.Println("[ERROR] Validating product", http.StatusBadRequest)

			resp.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&ValidationError{Messages: vErr.Errors()}, resp)
			return
		}

		// Create context of validated product
		ctx := context.WithValue(req.Context(), KeyProduct{}, *prod)

		// Send the validated product request to the next handler
		req = req.WithContext(ctx)
		next.ServeHTTP(resp, req)
	})
}
