package api

import (
	"github.com/golang-jwt/jwt/v4"
	"gym_management_system/db"
	"time"
)

type CreateAccountRequest struct {
	FirstName string `json:"firstName" validate:"alpha,min=3,max=30"`
	LastName  string `json:"lastName" validate:"alpha,min=3,max=30"`
	Email     string `json:"email" validate:"email"`
	Password  string `json:"password" validate:"password,min=6,max=30"`
}

type UpdateAccountRequest struct {
	FirstName       *string `json:"firstName,omitempty" validate:"omitempty,alpha,min=3,max=30"`
	LastName        *string `json:"lastName,omitempty" validate:"omitempty,alpha,min=3,max=30"`
	Email           *string `json:"email,omitempty" validate:"omitempty,email"`
	OldPassword     *string `json:"oldPassword,omitempty" validate:"omitempty,required_with=NewPassword"`
	NewPassword     *string `json:"newPassword,omitempty" validate:"omitempty,required_with=OldPassword,password,min=6,max=30"`
	RechargedCredit *int    `json:"rechargedCredit,omitempty" validate:"omitempty,min=0"`
}

type CreateAccountMembershipRequest struct {
	ValidFrom time.Time `json:"validFrom" validate:"gteCurrentDay"`
}

type CreateEntryRequest struct {
	AccountMembershipId *int `json:"accountMembershipId,omitempty"`
}

type CreateEventRequest struct {
	Type     db.EventType `json:"type" validate:"oneof=open_gym lecture all"`
	Title    string       `json:"title" validate:"min=3,max=30"`
	Start    time.Time    `json:"start" validate:"gtNow"`
	End      time.Time    `json:"end" validate:"gtfield=Start"`
	Capacity int          `json:"capacity" validate:"gt=0"`
	Price    int          `json:"price" validate:"gte=0"`
}

type CreateMembershipRequest struct {
	Type         db.EventType `json:"type" validate:"oneof=open_gym lecture all"`
	DurationDays int          `json:"durationDays" validate:"gte=1"`
	Entries      int          `json:"entries" validate:"gte=1"`
	Price        int          `json:"price" validate:"gte=0"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"-"`
}

type Claims struct {
	Id int `json:"id"`
	jwt.RegisteredClaims
}
