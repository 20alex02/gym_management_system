package db

type EntryRepository interface {
	CreateEntry(e *Entry) (*Entry, error)
	GetEntryById(id int) (*Entry, error)
	GetAllEntriesByAccountId(id int) (*[]Entry, error)
	DeleteEntry(id int) error
}

func getMembershipId(m *Membership) interface{} {
	if m == nil {
		return nil
	}
	return m.Id
}

// TODO
func (s *PostgresStore) CreateEntry(e *Entry) (*Entry, error) {
	//tx, err := s.Db.Begin()
	//if err != nil {
	//	return nil, err
	//}
	//defer commitOrRollback(tx, &err)

	query := `insert into entry (
                   account_id,
                   event_id,
                   membership_id
	) values ($1, $2, $3)`

	row := s.Db.QueryRow(
		query,
		e.Account.Id,
		e.Event.Id,
		getMembershipId(e.Membership))

	entry := &Entry{}
	if err := scanRow(row, entry); err != nil {
		return nil, err
	}
	return entry, nil
}

func (s *PostgresStore) GetEntryById(id int) (*Entry, error) {
	query := `select * from entry where id = $1`
	row := s.Db.QueryRow(query, id)
	entry := &Entry{}
	if err := scanRow(row, entry); err != nil {
		return nil, err
	}
	return entry, nil
}

func (s *PostgresStore) GetAllEntriesByAccountId(id int) (*[]Entry, error) {
	query := `select * from entry where account_id = $1`
	rows, err := s.Db.Query(query, id)
	if err != nil {
		return nil, err
	}
	entries := &[]Entry{}
	if err := scanRows(rows, entries); err != nil {
		return nil, err
	}
	return entries, nil
}

// TODO refund
func (s *PostgresStore) DeleteEntry(id int) error {
	query := `update entry set deleted_at = current_timestamp where id = $1 and deleted_at is null`
	_, err := s.Db.Exec(query, id)
	return err
}
