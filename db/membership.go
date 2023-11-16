package db

type MembershipRepository interface {
	CreateMembership(m *Membership) (int, error)
	GetAllMembershipsByAccountId(id int) (*[]Membership, error)
	GetMembershipById(id int) (*Membership, error)
}

func (s *PostgresStore) CreateMembership(m *Membership) (int, error) {
	query := `insert into membership (
						type,
						valid_from, 
						valid_to,
						entries_left,
						price,
						account_id
	) values ($1, $2, $3, $4, $5, $6) returning id`

	var id int
	err := s.Db.QueryRow(
		query,
		m.Type,
		m.ValidFrom,
		m.ValidTo,
		m.EntriesLeft,
		m.Price,
		m.Account.Id).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *PostgresStore) GetAllMembershipsByAccountId(id int) (*[]Membership, error) {
	query := `select * from membership 
         where account_id = $1 
           and deleted_at is null
           and valid_from <= current_timestamp
           and (valid_to is null or valid_to >= current_timestamp);
    `
	rows, err := s.Db.Query(query, id)
	if err != nil {
		return nil, err
	}
	memberships := &[]Membership{}
	if err := scanRows(rows, memberships); err != nil {
		return nil, err
	}
	return memberships, nil
}

func (s *PostgresStore) GetMembershipById(id int) (*Membership, error) {
	query := `select * from membership where id = $1`
	row := s.Db.QueryRow(query, id)
	membership := &Membership{}
	if err := scanRow(row, membership); err != nil {
		return nil, err
	}
	return membership, nil
}
