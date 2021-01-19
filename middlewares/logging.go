package middlewares

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

//LoggingMiddleware logs the incoming requests
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Infof("request [url = %v] [method = %v], [remote = %v], [protocol = %v]",
			r.RequestURI,
			r.Method,
			r.RemoteAddr,
			r.Proto)
		next.ServeHTTP(w, r)
	})
}
