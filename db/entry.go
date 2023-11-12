package db

type EntryRepository interface{}

func getMembershipId(m *Membership) interface{} {
	if m == nil {
		return nil
	}
	return m.Id
}

func (s *PostgresStore) CreateEntry(e *Entry) error {
	query := `insert into entry (
                   account_id,
                   event_id,
                   membership_id
	) values ($1, $2, $3)`

	_, err := s.Db.Query(
		query,
		e.Account.Id,
		e.Event.Id,
		getMembershipId(e.Membership))

	return err
}
