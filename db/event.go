package db

import (
	customErr "gym_management_system/errors"
	"time"
)

type EventRepository interface {
	CreateEvent(e *Event) (int, error)
	GetAllEvents(from time.Time, to time.Time) (*[]Event, error)
	//GetEventById(id int) (*Event, error)
	//UpdateEvent(e *Event) error
	DeleteEvent(id int) error
}

func (s *PostgresStore) CreateEvent(e *Event) (int, error) {
	// TODO check there are no overlaps
	query := `insert into event (
		type,
		title,
		start_time,
		end_time,
		capacity,
        price
	) values ($1, $2, $3, $4, $5, $6) returning id`

	var id int
	err := s.Db.QueryRow(
		query,
		e.Type,
		e.Title,
		e.Start,
		e.End,
		e.Capacity,
		e.Price).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
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

//func (s *PostgresStore) GetEventById(id int) (*Event, error) {
//	query := `select * from event where id = $1`
//	row := s.Db.QueryRow(query, id)
//	event := &Event{}
//	if err := scanRow(row, event); err != nil {
//		if errors.Is(err, sql.ErrNoRows) {
//			return nil, customErr.RecordNotFound{Record: "event", Property: "id", Value: id}
//		}
//		return nil, err
//	}
//	if event.DeletedAt != nil {
//		return nil, customErr.DeletedRecord{Record: "event", Property: "id", Value: id}
//	}
//	return event, nil
//}

//func (s *PostgresStore) UpdateEvent(e *Event) error {
//	query := `update event set
//                 type = $1,
//                 title = $2,
//            	 start_time = $3,
//            	 end_time = $4,
//            	 capacity = $5,
//            	 price = $6
//             where id = $7`
//
//	_, errors := s.Db.Exec(query, e.Type, e.Title, e.Start, e.End, e.Capacity, e.Price, e.Id)
//	return errors
//}

func (s *PostgresStore) DeleteEvent(id int) error {
	tx, err := s.Db.Begin()
	if err != nil {
		return err
	}
	defer commitOrRollback(tx, &err)
	event := Event{}
	err = getRecord(tx, EVENT, id, &event)
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
	query = `update event set deleted_at = current_timestamp
             where id = $1 and deleted_at is null`
	_, err = tx.Exec(query, id)
	return err
}
