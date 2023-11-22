package db

import (
	"database/sql"
	"errors"
	customErr "gym_management_system/errors"
	"time"
)

type EventRepository interface {
	CreateEvent(e *Event) (int, error)
	GetAllEvents(from time.Time, to time.Time) (*[]Event, error)
	GetAccountEvents(accountId int) (*[]EventWithEntryId, error)
	DeleteEvent(id int) error
}

func (s *PostgresStore) CreateEvent(e *Event) (int, error) {
	tx, err := s.Db.Begin()
	if err != nil {
		return 0, err
	}
	defer commitOrRollback(tx, &err)

	overlapQuery := `select id from event
		where ($1, $2) overlaps (start_time, end_time) and deleted_at is null`
	var overlappingEventID int
	err = tx.QueryRow(overlapQuery, e.Start, e.End).Scan(&overlappingEventID)
	if err == nil {
		err = customErr.InvalidRequest{Message: "overlapping events"}
		return 0, err
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}

	query := `insert into event (
		type,
		title,
		start_time,
		end_time,
		capacity,
        price
	) values ($1, $2, $3, $4, $5, $6) returning id`
	var id int
	err = tx.QueryRow(
		query,
		e.Type,
		e.Title,
		e.Start,
		e.End,
		e.Capacity,
		e.Price).Scan(&id)
	return id, err
}

func (s *PostgresStore) GetAllEvents(from time.Time, to time.Time) (*[]Event, error) {
	query := `select * from event 
         where deleted_at is null 
           and start_time >= $1
           and end_time <= $2`

	rows, err := s.Db.Query(query, from, to)
	if err != nil {
		return nil, err
	}

	events := &[]Event{}
	if err := scanRows(rows, events); err != nil {
		return nil, err
	}
	return events, nil
}

func (s *PostgresStore) GetAccountEvents(accountId int) (*[]EventWithEntryId, error) {
	query := `select event.*, entry.id as entry_id from entry 
    inner join event on entry.event_id = event.id
    where entry.account_id = $1 and entry.deleted_at is null`
	rows, err := s.Db.Query(query, accountId)
	if err != nil {
		return nil, err
	}
	events := &[]EventWithEntryId{}
	err = scanRows(rows, events)
	return events, err
}

func (s *PostgresStore) DeleteEvent(id int) error {
	tx, err := s.Db.Begin()
	if err != nil {
		return err
	}
	defer commitOrRollback(tx, &err)

	now := time.Now()
	event := Event{}
	err = getRecord(tx, EVENT, id, &event)
	if err != nil {
		return err
	}
	if event.Start.Before(now) && event.End.After(now) {
		err = customErr.InvalidRequest{Message: "event in progress"}
		return err
	}
	if event.End.Before(now) {
		query := `update entry set deleted_at = current_timestamp
             where event_id = $1 and deleted_at is null`
		_, err = tx.Exec(query, id)
		query = `update event set deleted_at = current_timestamp
             where id = $1 and deleted_at is null`
		_, err = tx.Exec(query, id)
		return err
	}
	query := `select * from entry where event_id = $1 and deleted_at is null`
	rows, err := tx.Query(query, id)
	if err != nil {
		return err
	}
	var entries []Entry
	err = scanRows(rows, &entries)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.AccountMembershipId == nil {
			account := Account{}
			err = getRecord(tx, ACCOUNT, entry.AccountId, &account)
			if err != nil {
				return err
			}
			query = `update account set credit = $1 where id = $2`
			_, err = tx.Exec(query, account.Credit+event.Price, account.Id)
			if err != nil {
				return err
			}
		} else {
			membership := AccountMembership{}
			err = getRecord(tx, ACCOUNT_MEMBERSHIP, *entry.AccountMembershipId, &membership)
			if err != nil {
				return err
			}
			query = `update account_membership set entries = $1 where id = $2`
			_, err = tx.Exec(query, membership.Entries+1, membership.Id)
			if err != nil {
				return err
			}
		}
	}
	query = `update entry set deleted_at = current_timestamp
             where event_id = $1 and deleted_at is null`
	_, err = tx.Exec(query, id)
	query = `update event set deleted_at = current_timestamp
             where id = $1 and deleted_at is null`
	_, err = tx.Exec(query, id)
	return err
}
