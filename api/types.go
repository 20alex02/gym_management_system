package api

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type CreateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Id    int    `json:"id"`
	Token string `json:"token"`
}

type GetEventsRequest struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

type Claims struct {
	Id int `json:"id"`
	jwt.RegisteredClaims
}
