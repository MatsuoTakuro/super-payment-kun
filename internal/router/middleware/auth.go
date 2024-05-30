package middleware

import (
	"net/http"
	"super-payment-kun/internal/handler"
	"super-payment-kun/internal/pkg"
)

func AuthJWT(j *pkg.JWTer) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			req, err := j.Validate(r)
			if err != nil {
				handler.DeriveStatusAndRespond(r.Context(), w, err)
				return
			}
			next.ServeHTTP(w, req)
		})
	}
}
