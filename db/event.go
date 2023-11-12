package db

type EventRepository interface{}

func (s *PostgresStore) CreateEvent(e *Event) error {
	query := `insert into event (
		type,
		title, 
		start_time,
		end_time,
		capacity,
        price
	) values ($1, $2, $3, $4, $5, $6)`

	_, err := s.Db.Query(
		query,
		e.Type,
		e.Title,
		e.Start,
		e.End,
		e.Capacity,
		e.Price)

	return err
}
