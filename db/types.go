package db

import "time"

type EventType string

const (
	OpenGym EventType = "open_gym"
	Lecture EventType = "lecture"
	All     EventType = "all"
)

type Account struct {
	Id                int
	FirstName         string
	LastName          string
	EncryptedPassword string
	Email             string
	Credit            int

	CreatedAt time.Time
	DeletedAt *time.Time
}

type Event struct {
	Id       int
	Type     EventType
	Title    string
	Start    time.Time
	End      time.Time
	Capacity int
	Price    int

	CreatedAt time.Time
	DeletedAt *time.Time
}

type Membership struct {
	Type        EventType
	ValidFrom   time.Time
	ValidTo     time.Time
	EntriesLeft int
	Price       int
	*Account

	CreatedAt time.Time
	DeletedAt *time.Time
}

type Entry struct {
	*Account
	*Membership
	*Event

	CreatedAt time.Time
	DeletedAt *time.Time
}
