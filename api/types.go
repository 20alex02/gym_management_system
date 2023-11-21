package api

import (
	"github.com/golang-jwt/jwt/v4"
	"gym_management_system/db"
	"time"
)

type CreateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UpdateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Credit    int    `json:"credit"`
}

type CreateAccountMembershipRequest struct {
	ValidFrom time.Time `json:"validFrom"`
}

type CreateEntryRequest struct {
	AccountMembershipId *int `json:"accountMembershipId"`
}

type CreateEventRequest struct {
	Type     db.EventType `json:"type"`
	Title    string       `json:"title"`
	Start    time.Time    `json:"start"`
	End      time.Time    `json:"end"`
	Capacity int          `json:"capacity"`
	Price    int          `json:"price"`
}

type CreateMembershipRequest struct {
	Type         db.EventType `json:"type"`
	DurationDays int          `json:"durationDays"`
	Entries      int          `json:"entries"`
	Price        int          `json:"price"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claims struct {
	Id int `json:"id"`
	jwt.RegisteredClaims
}
