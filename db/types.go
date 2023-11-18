package db

import "time"

type EventType string

const (
	OPEN_GYM EventType = "open_gym"
	LECTURE  EventType = "lecture"
	ALL      EventType = "all"
)

type Account struct {
	Id                int    `json:"id"`
	FirstName         string `json:"firstName"`
	LastName          string `json:"lastName"`
	EncryptedPassword string `json:"-"`
	Email             string `json:"email"`
	Credit            int    `json:"credit"`

	CreatedAt time.Time  `json:"createdAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}

type Event struct {
	Id       int       `json:"id"`
	Type     EventType `json:"type"`
	Title    string    `json:"title"`
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	Capacity int       `json:"capacity"`
	Price    int       `json:"price"`

	CreatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
}

type Membership struct {
	Id           int       `json:"id"`
	Type         EventType `json:"type"`
	DurationDays int       `json:"durationDays"`
	Entries      int       `json:"entries"`
	Price        int       `json:"price"`

	CreatedAt time.Time  `json:"createdAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}

type AccountMembership struct {
	Id           int       `json:"id"`
	AccountId    int       `json:"accountId"`
	MembershipId int       `json:"membershipId"`
	ValidFrom    time.Time `json:"validFrom"`
	ValidTo      time.Time `json:"validTo"`
	Entries      int       `json:"entries"`

	CreatedAt time.Time  `json:"createdAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}

type Entry struct {
	Id int `json:"id"`
	*Account
	*Membership
	*Event

	CreatedAt time.Time  `json:"createdAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}
