package db

type AccountRepository interface {
	CreateAccount(a *Account) (int, error)
	GetAccountById(id int) (*Account, error)
	GetAllAccounts() (*[]Account, error)
	UpdateAccount(a *Account) error
	DeleteAccount(id int) error
}

func (s *PostgresStore) CreateAccount(a *Account) (int, error) {
	query := `insert into account (
		first_name,
		last_name, 
		encrypted_password,
		email,
		credit
	) values ($1, $2, $3, $4, $5) returning id`

	var id int
	err := s.Db.QueryRow(
		query,
		a.FirstName,
		a.LastName,
		a.EncryptedPassword,
		a.Email,
		a.Credit).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *PostgresStore) GetAccountById(id int) (*Account, error) {
	query := `select * from account where id = $1`
	row := s.Db.QueryRow(query, id)
	account := &Account{}
	if err := scanRow(row, account); err != nil {
		return nil, err
	}
	return account, nil
}

func (s *PostgresStore) GetAllAccounts() (*[]Account, error) {
	query := `select * from account where deleted_at is null`
	rows, err := s.Db.Query(query)
	if err != nil {
		return nil, err
	}
	defer closeRows(rows)
	accounts := &[]Account{}
	err = scanRows(rows, accounts)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

// TODO update acc info or update credit
func (s *PostgresStore) UpdateAccount(a *Account) error {
	query := `update account set 
                   first_name = $1,
                   last_name = $2,
                   encrypted_password = $3,
                   email = $4,
                   credit = $5
               where id = $6`
	_, err := s.Db.Exec(query, a.FirstName, a.LastName, a.EncryptedPassword, a.Email, a.Credit)
	return err
}

func (s *PostgresStore) DeleteAccount(id int) error {
	tx, err := s.Db.Begin()
	if err != nil {
		return err
	}
	defer commitOrRollback(tx, &err)

	err = checkDeleted(tx, "account", id)
	if err != nil {
		return err
	}
	query := `update account set deleted_at = current_timestamp where id = $1`
	_, err = tx.Exec(query, id)
	return err
}
