package api

import (
	"errors"
	"net.vikesh/goshop/db"
	"net.vikesh/goshop/dto"
	"net/http"
)

var userService = &db.UserService{}
func registerUserHandler(w http.ResponseWriter, r *http.Request) {
	form := &dto.RegistrationForm{}
	decodeError := decode(form, r)
	if decodeError != nil {
		toJson(success(nil), w, failedWithStatus(errors.New("invalid content"), http.StatusNotAcceptable))
		return
	}
	if len(form.DisplayName) == 0 || len(form.Email) == 0 || len(form.Password) == 0 || len(form.Username) == 0 {
		toJson(nil, w, failedWithStatus(errors.New("all required fields have not been entered"), http.StatusNotAcceptable))
		return
	}
	err := userService.RegisterNewUser(form)
	if err != nil {
		toJson(success(nil), w, failed(err))
	}
}
