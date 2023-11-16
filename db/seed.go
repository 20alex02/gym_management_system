package db

import (
	"time"
)

func (s *PostgresStore) Seed() error {
	tables := []string{"entry", "membership", "account", "event"}
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
	if err = s.createEntryTable(); err != nil {
		return err
	}

	accounts := []Account{
		{1, "john", "doe", "12345", "jDoe@mail.com", 0, time.Now(), nil},
		{2, "jack", "johnson", "ashfsaf", "jj@mail.com", 0, time.Now(), nil},
		{3, "john", "jackson", "lengogn", "johnJ@mail.com", 0, time.Now(), nil},
	}
	event := Event{0, OPEN_GYM, "open", time.Now(), time.Now().Add(time.Hour), 10, 100, time.Now(), nil}
	membership := Membership{0, ALL, time.Now().AddDate(0, 0, -1), time.Now().AddDate(0, 1, 0), 30, 1500, &accounts[0], time.Now(), nil}
	for _, account := range accounts {
		if _, err := s.CreateAccount(&account); err != nil {
			return err
		}
	}
	//acc, errors := s.GetAllAccounts()
	//if errors != nil {
	//	return errors
	//}
	//fmt.Println(acc)
	_, err = s.CreateEvent(&event)
	if err != nil {
		return err
	}
	_, err = s.CreateMembership(&membership)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStore) cleanup(tables []string) error {
	for i := 0; i < len(tables); i++ {
		query := `drop table if exists $1`
		if _, err := s.Db.Exec(query, tables[i]); err != nil {
			return err
		}
	}
	_, err := s.Db.Exec("drop type if exists event_type")
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
		account_id int,
		type event_type,
		valid_from timestamp,
		valid_to timestamp,
		entries_left int,
		price int,
		created_at timestamp default current_timestamp,
		deleted_at timestamp default null,
		foreign key (account_id) references account(id)
	)`
	_, err := s.Db.Exec(query)
	return err
}

func (s *PostgresStore) createEntryTable() error {
	query := `create table if not exists entry (
		id serial primary key,
		account_id int,
		event_id int,
		membership_id int null,
		created_at timestamp default current_timestamp,
		deleted_at timestamp default null,
		foreign key (account_id) references account(id),
		foreign key (event_id) references event(id),
		foreign key (membership_id) references membership(id)
	)`
	_, err := s.Db.Exec(query)
	return err
}

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
