package middlewares

import (
	"net.vikesh/goshop/config"
	"net.vikesh/goshop/db"
	"net/http"
	"strings"
)

var userService = &db.UserService{}
//Security filter to configure security rules on path access
func Secured(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isOpenUrl(r) {
			next.ServeHTTP(w, r)
		} else {
			cookies := r.Cookies()
			var valid bool
			for _, cookie := range cookies {
				cookieName := cookie.Name
				if cookieName == config.Get().GetString(config.ApiCookieName) {
					valid = userService.IsValidToken(cookie.Value)
				}
			}
			if valid {
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, "not authorized", http.StatusUnauthorized)
			}
		}
	})
}

func isOpenUrl(r *http.Request) bool {
	uri := r.RequestURI
	return strings.Contains(uri, "/uploads/") || strings.Index(uri, "/api/") == -1 || isOpenApi(r)
}

func isOpenApi(r *http.Request) bool {
	uri := r.RequestURI
	return strings.Index(uri, "/api/") == 0 &&
		(strings.Contains(uri, "/api/authenticate") ||
			strings.Contains(uri, "/api/.whoami") ||
			strings.Contains(uri, "/api/authenticate") ||
			strings.Contains(uri, "/api/register"))
}
