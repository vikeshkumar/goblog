package api

import (
	"net/http"
	"strings"
)

const userCookieName string = `session`

func whoamiHandler(w http.ResponseWriter, r *http.Request) {
	var userCookieValue string
	cookies := r.Cookies()
	for i := 0; i < len(cookies) && len(userCookieValue) == 0; i++ {
		if strings.EqualFold(userCookieName, cookies[i].Name) {
			userCookieValue = cookies[i].Value
		}
	}
	user, _ := userService.FindUserForCookieValue(userCookieValue)
	toJson(success(&user), w, failed(nil))
}
