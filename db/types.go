package db

import (
	"time"
)

type EventType string

const (
	OPEN_GYM EventType = "open_gym"
	LECTURE  EventType = "lecture"
	ALL      EventType = "all"
)

type Role string

const (
	ADMIN Role = "admin"
	USER  Role = "user"
)

type Table string

const (
	ACCOUNT            Table = "account"
	EVENT              Table = "event"
	MEMBERSHIP         Table = "membership"
	ENTRY              Table = "entry"
	ACCOUNT_MEMBERSHIP Table = "account_membership"
)

type Account struct {
	Id                int    `json:"id"`
	FirstName         string `json:"firstName"`
	LastName          string `json:"lastName"`
	EncryptedPassword string `json:"-"`
	Email             string `json:"email"`
	Credit            int    `json:"credit"`
	Role              Role   `json:"-"`

	CreatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
}

type Event struct {
	Id           int       `json:"id"`
	Type         EventType `json:"type"`
	Title        string    `json:"title"`
	Start        time.Time `json:"start"`
	End          time.Time `json:"end"`
	Capacity     int       `json:"capacity"`
	Participants int       `json:"participants"`
	Price        int       `json:"price"`

	CreatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
}

type Membership struct {
	Id           int       `json:"id"`
	Type         EventType `json:"type"`
	DurationDays int       `json:"durationDays"`
	Entries      int       `json:"entries"`
	Price        int       `json:"price"`

	CreatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
}

type AccountMembership struct {
	Id           int       `json:"id"`
	AccountId    int       `json:"-"`
	MembershipId int       `json:"-"`
	ValidFrom    time.Time `json:"validFrom"`
	ValidTo      time.Time `json:"validTo"`
	Entries      int       `json:"entries"`

	CreatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
}

type Entry struct {
	Id                  int  `json:"id"`
	AccountId           int  `json:"accountId"`
	EventId             int  `json:"eventId"`
	AccountMembershipId *int `json:"accountMembershipId"`

	CreatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
}

type EventWithEntryId struct {
	Id           int       `json:"id"`
	Type         EventType `json:"type"`
	Title        string    `json:"title"`
	Start        time.Time `json:"start"`
	End          time.Time `json:"end"`
	Capacity     int       `json:"capacity"`
	Participants int       `json:"participants"`
	Price        int       `json:"price"`

	CreatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`

	EntryId int `json:"entryId"`
}

type EventEntry struct {
	Id        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type AccountMembershipWithType struct {
	Id           int       `json:"id"`
	AccountId    int       `json:"-"`
	MembershipId int       `json:"-"`
	ValidFrom    time.Time `json:"validFrom"`
	ValidTo      time.Time `json:"validTo"`
	Entries      int       `json:"entries"`

	CreatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`

	Type EventType `json:"type"`
}
