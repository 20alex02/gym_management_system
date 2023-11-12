package db

type MembershipRepository interface {
}

func (s *PostgresStore) CreateMembership(m *Membership) error {
	query := `insert into membership (
						type,
						valid_from, 
						valid_to,
						entries_left,
						price,
						account_id
	) values ($1, $2, $3, $4, $5, $6)`

	_, err := s.Db.Query(
		query,
		m.Type,
		m.ValidFrom,
		m.ValidTo,
		m.EntriesLeft,
		m.Price,
		m.Account.Id)

	return err
}
