package db

import (
	"fmt"
	"time"
)

func (s *PostgresStore) Seed() error {
	tables := []string{"entry", "account_membership", "membership", "account", "event"}
	var err error
	if err = s.cleanup(tables); err != nil {
		return err
	}

	if err = s.createEventType(); err != nil {
		return err
	}
	if err = s.createAccountTable(); err != nil {
		return err
	}
	if err = s.createEventTable(); err != nil {
		return err
	}
	if err = s.createMembershipTable(); err != nil {
		return err
	}
	if err = s.createAccountMembershipTable(); err != nil {
		return err
	}
	if err = s.createEntryTable(); err != nil {
		return err
	}

	hours := []int{9, 12, 18}
	now := time.Now()
	startOfWeek := now.AddDate(0, 0, 1-int(now.Weekday()))
	var events []Event
	for i := 0; i < 7; i++ {
		dayStart := startOfWeek.AddDate(0, 0, i)
		for _, hour := range hours {
			startTime := time.Date(dayStart.Year(), dayStart.Month(), dayStart.Day(), hour, 0, 0, 0, time.Local)
			endTime := startTime.Add(time.Hour)
			event := Event{
				Type:     OPEN_GYM,
				Title:    fmt.Sprintf("Event at %s", startTime),
				Start:    startTime,
				End:      endTime,
				Capacity: 10,
				Price:    100,
			}
			events = append(events, event)
		}
	}
	for _, event := range events {
		if _, err = s.CreateEvent(&event); err != nil {
			return err
		}
	}

	memberships := []Membership{
		{1, OPEN_GYM, 30, 30, 1000, time.Now(), nil},
		{2, LECTURE, 30, 30, 1500, time.Now(), nil},
		{3, ALL, 30, 30, 2000, time.Now(), nil},
	}
	for _, membership := range memberships {
		if _, err = s.CreateMembership(&membership); err != nil {
			return err
		}
	}

	//accounts := []Account{
	//	{1, "john", "doe", "12345", "jDoe@mail.com", 0, time.Now(), nil},
	//	{2, "jack", "johnson", "ashfsaf", "jj@mail.com", 0, time.Now(), nil},
	//	{3, "john", "jackson", "lengogn", "johnJ@mail.com", 0, time.Now(), nil},
	//}
	//for _, account := range accounts {
	//	if _, err := s.CreateAccount(&account); err != nil {
	//		return err
	//	}
	//}

	//accountMembership := AccountMembership{1, 1, 3, time.Now(), time.Now().AddDate(0, 1, 0), 30, time.Now(), nil}
	//if _, err = s.CreateAccountMembership(&accountMembership); err != nil {
	//	return err
	//}
	return nil
}

func (s *PostgresStore) cleanup(tables []string) error {
	for i := 0; i < len(tables); i++ {
		query := fmt.Sprintf(`drop table if exists %s`, tables[i])
		if _, err := s.Db.Exec(query); err != nil {
			return err
		}
	}
	_, err := s.Db.Exec(`drop type if exists event_type`)
	return err
}

func (s *PostgresStore) createAccountTable() error {
	query := `create table if not exists account (
		id serial primary key,
		first_name varchar(255),
		last_name varchar(255),
		encrypted_password varchar(255),
		email varchar(255) unique,
		credit int,
		created_at timestamp default current_timestamp,
		deleted_at timestamp default null
	)`
	_, err := s.Db.Exec(query)
	return err
}

func (s *PostgresStore) createMembershipTable() error {
	query := `create table if not exists membership (
		id serial primary key,
		type event_type,
		duration_days int,
		entries int,
		price int,
		created_at timestamp default current_timestamp,
		deleted_at timestamp default null
	)`
	_, err := s.Db.Exec(query)
	return err
}

func (s *PostgresStore) createAccountMembershipTable() error {
	query := `create table if not exists account_membership (
		id serial primary key,
		account_id int,
		membership_id int,
		valid_from timestamp,
		valid_to timestamp,
		entries int,
		created_at timestamp default current_timestamp,
		deleted_at timestamp default null,
		foreign key (account_id) references account(id),
		foreign key (membership_id) references membership(id)
	)`
	_, err := s.Db.Exec(query)
	return err
}

func (s *PostgresStore) createEntryTable() error {
	query := `create table if not exists entry (
		id serial primary key,
		account_id int,
		event_id int,
		account_membership_id int null,
		created_at timestamp default current_timestamp,
		deleted_at timestamp default null,
		foreign key (account_id) references account(id),
		foreign key (event_id) references event(id),
		foreign key (account_membership_id) references account_membership(id)
	)`
	_, err := s.Db.Exec(query)
	return err
}

// TODO capacity and capacity_left
func (s *PostgresStore) createEventTable() error {
	query := `create table if not exists event (
		id serial primary key,
		type event_type,
		title varchar(255),
		start_time timestamp,
		end_time timestamp,
		capacity int,
		price int,
		created_at timestamp default current_timestamp,
		deleted_at timestamp default null
	)`
	_, err := s.Db.Exec(query)
	return err
}

func (s *PostgresStore) createEventType() error {
	query := `create type event_type as enum (
		'open_gym',
		'lecture',
		'all'
	)`
	_, err := s.Db.Exec(query)
	return err
}
