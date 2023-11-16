package db

import "time"

type EventRepository interface {
	CreateEvent(e *Event) (int, error)
	GetAllEvents(from time.Time, to time.Time) (*[]Event, error)
	GetEventById(id int) (*Event, error)
	//UpdateEvent(e *Event) error
	DeleteEvent(id int) error
}

func (s *PostgresStore) CreateEvent(e *Event) (int, error) {
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

func (s *PostgresStore) GetEventById(id int) (*Event, error) {
	query := `select * from event where id = $1`
	row := s.Db.QueryRow(query, id)
	event := &Event{}
	if err := scanRow(row, event); err != nil {
		return nil, err
	}
	return event, nil
}

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

// TODO delete all entries linked to the event and do refunds if the event has not started yet
func (s *PostgresStore) DeleteEvent(id int) error {
	query := `update event set deleted_at = current_timestamp 
             where id = $1 and deleted_at is null`

	_, err := s.Db.Exec(query, id)
	return err
}
