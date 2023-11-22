package db

type MembershipRepository interface {
	CreateMembership(m *Membership) (int, error)
	GetAllMemberships() (*[]Membership, error)
	DeleteMembership(id int) error
}

func (s *PostgresStore) CreateMembership(m *Membership) (int, error) {
	query := `insert into membership (
						type,
						duration_days,
						entries,
						price
	) values ($1, $2, $3, $4) returning id`

	var id int
	err := s.Db.QueryRow(
		query,
		m.Type,
		m.DurationDays,
		m.Entries,
		m.Price).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *PostgresStore) GetAllMemberships() (*[]Membership, error) {
	query := `select * from membership where deleted_at is null`
	rows, err := s.Db.Query(query)
	if err != nil {
		return nil, err
	}
	memberships := &[]Membership{}
	if err := scanRows(rows, memberships); err != nil {
		return nil, err
	}
	return memberships, nil
}

func (s *PostgresStore) DeleteMembership(id int) error {
	tx, err := s.Db.Begin()
	if err != nil {
		return err
	}
	defer commitOrRollback(tx, &err)

	membership := Membership{}
	err = getRecord(tx, MEMBERSHIP, id, &membership)
	if err != nil {
		return err
	}
	query := `update membership set deleted_at = current_timestamp where id = $1`
	_, err = tx.Exec(query, id)
	return err
}
