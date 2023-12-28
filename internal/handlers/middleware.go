package handlers

import "net/http"

// OnlyPostAllowed gets Handler, make validation on Post request and returns also Handler.
func OnlyPostAllowed(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// Only POST method allowed
		if req.Method != http.MethodPost {
			res.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(res, req)
	})
}
