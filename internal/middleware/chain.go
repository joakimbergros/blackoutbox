package middleware

import (
	"net/http"
	"slices"
)

func ChainMiddleware(handler http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
	for _, mw := range slices.Backward(middleware) {
		handler = mw(handler)
	}

	return handler
}
