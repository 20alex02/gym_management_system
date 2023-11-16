package db

type EntryRepository interface {
	CreateEntry(e *Entry) (int, error)
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
func (s *PostgresStore) CreateEntry(e *Entry) (int, error) {
	//tx, errors := s.Db.Begin()
	//if errors != nil {
	//	return nil, errors
	//}
	//defer commitOrRollback(tx, &errors)

	query := `insert into entry (
                   account_id,
                   event_id,
                   membership_id
	) values ($1, $2, $3) returning id`

	var id int
	err := s.Db.QueryRow(
		query,
		e.Account.Id,
		e.Event.Id,
		getMembershipId(e.Membership)).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
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
