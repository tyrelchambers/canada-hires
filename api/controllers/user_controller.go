package controllers

import (
	"canada-hires/helpers"
	"canada-hires/services"
	"encoding/json"
	"net/http"
)

type UserController interface {
	GetUser(w http.ResponseWriter, r *http.Request)
}

type userController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) UserController {
	return &userController{userService: userService}
}

func (uc *userController) GetUser(w http.ResponseWriter, r *http.Request) {
	user := helpers.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
