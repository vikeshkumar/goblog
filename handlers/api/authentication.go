package api

import (
	"errors"
	"net.vikesh/goshop/config"
	"net.vikesh/goshop/dto"
	"net/http"
	"time"
)

func authenticationHandler(w http.ResponseWriter, r *http.Request) {
	type authentication struct {
		Token string `json:"token"`
	}
	loginForm := &dto.LoginForm{}
	decodeError := decode(loginForm, r)
	if decodeError != nil {
		toJson(success(nil), w, failedWithStatus(decodeError, http.StatusUnauthorized))
		return
	}
	id, user, findError := userService.FindUserByUserName(loginForm.Username)
	if findError != nil || id == 0 || len(user.UserName) == 0 {
		toJson(success(nil), w, failedWithStatus(errors.New("user not found"), http.StatusUnauthorized))
		return
	}
	token, tokenError := userService.CreateTokenForUser(id)
	if tokenError != nil {
		toJson(success(nil), w, failedWithStatus(tokenError, http.StatusUnauthorized))
		return
	}
	if cfg == nil {
		cfg = config.Get()
	}
	expiration := time.Now().Add(config.Get().GetDuration(config.ApiCookieValidity) * time.Minute)
	cookie := http.Cookie{
		Name:     config.Get().GetString(config.ApiCookieName),
		Value:    token,
		Domain:   config.Get().GetString(config.ApiCookieDomain),
		Expires:  expiration,
		Secure:   config.Get().GetBoolean(config.ApiCookieSecure),
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	toJson(success(&authentication{token}), w, nil)
}
