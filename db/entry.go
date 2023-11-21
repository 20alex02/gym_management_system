package db

import (
	"database/sql"
	"errors"
	customErr "gym_management_system/errors"
	"time"
)

type EntryRepository interface {
	CreateEntry(e *Entry) (int, error)
	GetEventEntries(eventId int) (*[]EventEntry, error)
	DeleteEntry(id, accountId int) error
}

func (s *PostgresStore) CreateEntry(e *Entry) (int, error) {
	tx, err := s.Db.Begin()
	if err != nil {
		return 0, err
	}
	defer commitOrRollback(tx, &err)

	event := Event{}
	err = getRecord(tx, EVENT, e.EventId, &event)
	if err != nil {
		return 0, err
	}
	account := Account{}
	err = getRecord(tx, ACCOUNT, e.AccountId, &account)
	if err != nil {
		return 0, err
	}

	query := `select 1 from entry where account_id = $1 and event_id = $2 and deleted_at is null`
	var result int
	err = tx.QueryRow(query, account.Id, event.Id).Scan(&result)
	if err == nil {
		err = customErr.InvalidRequest{Message: "already registered for the event"}
		return 0, err
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}

	if event.End.Before(time.Now()) {
		err = customErr.InvalidRequest{Message: "event already ended"}
		return 0, err
	}
	if event.Start.Before(time.Now()) {
		err = customErr.InvalidRequest{Message: "event already started"}
		return 0, err
	}
	if event.Participants == event.Capacity {
		err = customErr.InvalidRequest{Message: "full capacity"}
		return 0, err
	}

	if e.AccountMembershipId == nil {
		if account.Credit < event.Price {
			err = customErr.InsufficientResources{}
			return 0, err
		}
		query = `update account set credit = $1 where id = $2`
		_, err = tx.Exec(query, account.Credit-event.Price, account.Id)
		if err != nil {
			return 0, err
		}
	} else {
		accountMembership := AccountMembership{}
		err = getRecord(tx, ACCOUNT_MEMBERSHIP, *e.AccountMembershipId, &accountMembership)
		if err != nil {
			return 0, err
		}
		if event.Start.Before(accountMembership.ValidFrom) || event.End.After(accountMembership.ValidTo) {
			err = customErr.InvalidRequest{Message: "event does not occur within membership validity"}
			return 0, err
		}
		if accountMembership.Entries < 1 {
			err = customErr.InsufficientResources{}
			return 0, err
		}

		var eventType EventType
		query = `select type from membership where id = $1`
		err = tx.QueryRow(query, accountMembership.MembershipId).Scan(&eventType)
		if err != nil {
			return 0, err
		}
		if eventType != event.Type && eventType != ALL {
			err = customErr.InvalidRequest{Message: "invalid membership type"}
			return 0, err
		}
		query = `update account_membership set entries = $1 where id = $2`
		_, err = tx.Exec(query, accountMembership.Entries-1, accountMembership.Id)
		if err != nil {
			return 0, err
		}
	}
	query = `update event set participants = $1 where id = $2`
	_, err = tx.Exec(query, event.Participants+1, event.Id)
	if err != nil {
		return 0, err
	}

	query = `insert into entry (
                   account_id,
                   event_id,
                   account_membership_id
	) values ($1, $2, $3) returning id`

	var id int
	err = s.Db.QueryRow(
		query,
		e.AccountId,
		e.EventId,
		e.AccountMembershipId).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *PostgresStore) GetEventEntries(eventId int) (*[]EventEntry, error) {
	query := `select entry.id, account.first_name, account.last_name from entry 
    join account on entry.account_id = account.id where entry.event_id = $1 and entry.deleted_at is null`
	rows, err := s.Db.Query(query, eventId)
	if err != nil {
		return nil, err
	}
	entries := &[]EventEntry{}
	if err = scanRows(rows, entries); err != nil {
		return nil, err
	}
	return entries, nil
}

func (s *PostgresStore) DeleteEntry(id, accountId int) error {
	tx, err := s.Db.Begin()
	if err != nil {
		return err
	}
	defer commitOrRollback(tx, &err)

	entry := Entry{}
	err = getRecord(tx, ENTRY, id, &entry)
	if err != nil {
		return err
	}
	if entry.AccountId != accountId {
		err = customErr.InvalidRequest{Message: "can not delete"}
	}
	event := Event{}
	err = getRecord(tx, EVENT, entry.EventId, &event)
	if err != nil {
		return err
	}
	if event.Start.Before(time.Now()) {
		err = customErr.InvalidRequest{Message: "event already started"}
		return err
	}
	if event.End.Before(time.Now()) {
		err = customErr.InvalidRequest{Message: "event already ended"}
		return err
	}
	account := Account{}
	err = getRecord(tx, ACCOUNT, entry.AccountId, &account)
	if err != nil {
		return err
	}
	membership := AccountMembership{}
	err = getRecord(tx, ACCOUNT_MEMBERSHIP, entry.AccountId, &membership)
	if err != nil {
		return err
	}

	query := `update event set participants = $1 where id = $2`
	_, err = s.Db.Exec(query, event.Participants-1, event.Id)
	if err != nil {
		return err
	}
	if entry.AccountMembershipId == nil {
		query = `update account set credit = $1 where id = $2`
		_, err = s.Db.Exec(query, account.Credit+event.Price, account.Id)
		if err != nil {
			return err
		}
	} else {
		query = `update account_membership set entries = $1 where id = $2`
		_, err = s.Db.Exec(query, membership.Entries+1, membership.Id)
		if err != nil {
			return err
		}
	}
	query = `update entry set deleted_at = current_timestamp where id = $1 and deleted_at is null`
	_, err = s.Db.Exec(query, id)
	return err
}
