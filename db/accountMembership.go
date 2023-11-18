package db

type AccountMembershipRepository interface {
	CreateAccountMembership(m *AccountMembership) (int, error)
	GetAllAccountMemberships(id int) (*[]AccountMembership, error)
}

func (s *PostgresStore) CreateAccountMembership(m *AccountMembership) (int, error) {
	query := `insert into account_membership (
			account_id,
			membership_id,
			valid_from,
			valid_to,
			entries
		) values ($1, $2, $3, $4, $5) returning id`

	var id int
	err := s.Db.QueryRow(
		query,
		m.AccountId,
		m.MembershipId,
		m.ValidFrom,
		m.ValidTo,
		m.Entries).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *PostgresStore) GetAllAccountMemberships(id int) (*[]AccountMembership, error) {
	query := `select * from account_membership 
         where deleted_at is null 
           and account_id = $1
           and valid_to >= current_timestamp`
	rows, err := s.Db.Query(query, id)
	if err != nil {
		return nil, err
	}
	memberships := &[]AccountMembership{}
	if err := scanRows(rows, memberships); err != nil {
		return nil, err
	}
	return memberships, nil
}
