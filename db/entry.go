package db

import (
	"gym_management_system/errors"
	"time"
)

type EntryRepository interface {
	CreateEntry(e *Entry) (int, error)
	//GetEntryById(id int) (*Entry, error)
	GetAccountEntries(accountId int) (*[]Entry, error)
	GetEventEntries(eventId int) (*[]Entry, error)
	DeleteEntry(id int) error
}

func (s *PostgresStore) CreateEntry(e *Entry) (int, error) {
	tx, err := s.Db.Begin()
	if err != nil {
		return 0, err
	}
	defer commitOrRollback(tx, &err)

	event := Event{}
	err = checkRecord(tx, EVENT, e.EventId, &event)
	if err != nil {
		return 0, err
	}
	if event.Start.Before(time.Now()) {
		return 0, errors.InvalidRequest{Message: "event already started"}
	}
	if event.Capacity < 1 {
		return 0, errors.InvalidRequest{Message: "full capacity"}
	}
	query := `update event set capacity = $1 where id = $2`
	_, err = tx.Exec(query, event.Capacity+1, event.Id)
	if err != nil {
		return 0, err
	}

	account := Account{}
	err = checkRecord(tx, ACCOUNT, e.AccountId, &account)
	if err != nil {
		return 0, err
	}

	query = `select deleted_at from entry where account_id = $1 and event_id = $2`
	var deletedAt *time.Time
	err = tx.QueryRow(query, account.Id, event.Id).Scan(&deletedAt)

	if e.AccountMembershipId != nil {
		if account.Credit < event.Price {
			return 0, errors.InsufficientResources{}
		}
		query := `update account set credit = $1 where id = $2`
		_, err = tx.Exec(query, account.Credit-event.Price, account.Id)
		if err != nil {
			return 0, err
		}
	} else {
		accountMembership := AccountMembership{}
		err = checkRecord(tx, ACCOUNT_MEMBERSHIP, *e.AccountMembershipId, &accountMembership)
		if err != nil {
			return 0, err
		}
		if event.Start.Before(accountMembership.ValidFrom) || event.End.After(accountMembership.ValidTo) {
			return 0, errors.InvalidRequest{Message: "event does not occur within membership validity"}
		}
		if accountMembership.Entries < 1 {
			return 0, errors.InsufficientResources{}
		}

		membership := Membership{}
		err = checkRecord(tx, MEMBERSHIP, accountMembership.MembershipId, &membership)
		if err != nil {
			return 0, err
		}
		if event.Type != membership.Type && membership.Type != ALL {
			return 0, errors.InvalidRequest{Message: "invalid membership type"}
		}
		query := `update account_membership set entries = $1 where id = $2`
		_, err = tx.Exec(query, accountMembership.Entries-1, accountMembership.Id)
		if err != nil {
			return 0, err
		}
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

//func (s *PostgresStore) GetEntryById(id int) (*Entry, error) {
//	query := `select * from entry where id = $1`
//	row := s.Db.QueryRow(query, id)
//	entry := &Entry{}
//	if err := scanRow(row, entry); err != nil {
//		return nil, err
//	}
//	return entry, nil
//}

func (s *PostgresStore) GetAccountEntries(accountId int) (*[]Entry, error) {
	query := `select * from entry where account_id = $1 and deleted_at is null`
	rows, err := s.Db.Query(query, accountId)
	if err != nil {
		return nil, err
	}
	entries := &[]Entry{}
	if err := scanRows(rows, entries); err != nil {
		return nil, err
	}
	return entries, nil
}

func (s *PostgresStore) GetEventEntries(eventId int) (*[]Entry, error) {
	query := `select * from entry where event_id = $1 and deleted_at is null`
	rows, err := s.Db.Query(query, eventId)
	if err != nil {
		return nil, err
	}
	entries := &[]Entry{}
	if err := scanRows(rows, entries); err != nil {
		return nil, err
	}
	return entries, nil
}

func (s *PostgresStore) DeleteEntry(id int) error {
	tx, err := s.Db.Begin()
	if err != nil {
		return err
	}
	defer commitOrRollback(tx, &err)

	entry := Entry{}
	err = checkRecord(tx, ENTRY, id, &entry)
	if err != nil {
		return err
	}
	event := Event{}
	err = checkRecord(tx, EVENT, entry.EventId, &event)
	if err != nil {
		return err
	}
	if event.Start.Before(time.Now()) {
		return errors.InvalidRequest{Message: "event already started"}
	}
	account := Account{}
	err = checkRecord(tx, ACCOUNT, entry.AccountId, &account)
	if err != nil {
		return err
	}
	membership := AccountMembership{}
	err = checkRecord(tx, ACCOUNT_MEMBERSHIP, entry.AccountId, &membership)
	if err != nil {
		return err
	}

	query := `update event set capacity = $1 where id = $2`
	_, err = s.Db.Exec(query, event.Capacity+1, event.Id)
	if err != nil {
		return err
	}
	if entry.AccountMembershipId == nil {
		query := `update account set credit = $1 where id = $2`
		_, err = s.Db.Exec(query, account.Credit+event.Price, account.Id)
		if err != nil {
			return err
		}
	} else {
		query := `update account_membership set entries = $1 where id = $2`
		_, err = s.Db.Exec(query, membership.Entries+1, membership.Id)
		if err != nil {
			return err
		}
	}
	query = `update entry set deleted_at = current_timestamp where id = $1 and deleted_at is null`
	_, err = s.Db.Exec(query, id)
	return err
}
