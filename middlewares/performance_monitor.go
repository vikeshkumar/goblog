package middlewares

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

//PerformanceMonitor logs the incoming requests
func PerformanceMonitor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		defer log.Infof("response_time [%v] : [%v ns]", r.RequestURI, time.Now().Sub(start).Nanoseconds())
	})
}
