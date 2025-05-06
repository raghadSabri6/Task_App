package controllers

import (
	"log"
	"net/http"
	"task2/database"
	"task2/helperFunc"
	"task2/middlewares"
	"task2/models"
	"task2/schemas"
)

type Users struct {
	l *log.Logger
}

func NewUsers(l *log.Logger) *Users {
	return &Users{l}
}

func (u *Users) Signup(w http.ResponseWriter, r *http.Request) {

	val := r.Context().Value(middlewares.BindKey)
	userReq, ok := val.(*models.User)
	if !ok {
		helperFunc.RespondJSON(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	user, err := userReq.CreateUser(database.DB, userReq.Name, userReq.Email, userReq.Password)
	if err != nil {
		if err.Error() == "email already registered" {
			helperFunc.RespondJSON(w, http.StatusConflict, err.Error(), nil)
		} else {
			helperFunc.RespondJSON(w, http.StatusInternalServerError, err.Error(), nil)
		}
		return
	}

	if err := helperFunc.SendSignupEmail(user.Name, user.Email); err != nil {
		helperFunc.RespondJSON(w, http.StatusInternalServerError, "Failed to send email", nil)
		return
	}

	response := map[string]interface{}{
		"user": user,
	}
	helperFunc.RespondJSON(w, http.StatusCreated, "", response)
}

func (u *Users) Login(w http.ResponseWriter, r *http.Request) {

	val := r.Context().Value(middlewares.BindKey)
	loginReq, ok := val.(*schemas.LoginRequest)
	if !ok {
		helperFunc.RespondJSON(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	user, err := models.GetUserByEmail(loginReq.Email)
	if err != nil {
		helperFunc.RespondJSON(w, http.StatusUnauthorized, "Invalid email or password", nil)
		return
	}

	if err := helperFunc.VerifyPassword(user.Password, loginReq.Password); err != nil {
		helperFunc.RespondJSON(w, http.StatusUnauthorized, "Invalid email or password", nil)
		return
	}

	expiryDuration := helperFunc.GetTokenExpiration()
	tokenString, err := helperFunc.GenerateJWTToken(user.UUID, expiryDuration)
	if err != nil {
		helperFunc.RespondJSON(w, http.StatusInternalServerError, "Failed to generate token", nil)
		return
	}

	helperFunc.SetJWTTokenCookie(w, tokenString, expiryDuration)

	helperFunc.RespondJSON(w, http.StatusOK, "", map[string]interface{}{
		"user":  user,
		"token": tokenString,
	})
}

func (u *Users) GetUsers(w http.ResponseWriter, r *http.Request) {

	users, err := models.GetAllUsers()
	if err != nil {
		helperFunc.RespondJSON(w, http.StatusInternalServerError, "Failed to fetch users", nil)
		return
	}

	helperFunc.RespondJSON(w, http.StatusOK, "", map[string]interface{}{"users": users})
}
