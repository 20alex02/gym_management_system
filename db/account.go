package db

type AccountRepository interface {
	CreateAccount(*Account) error
	GetAccountByID(int) (*Account, error)
	GetAllAccounts() (*[]Account, error)
}

func (s *PostgresStore) CreateAccount(a *Account) error {
	query := `insert into account (
		first_name,
		last_name, 
		encrypted_password,
		email,
		credit
	) values ($1, $2, $3, $4, $5)`

	_, err := s.Db.Query(
		query,
		a.FirstName,
		a.LastName,
		a.EncryptedPassword,
		a.Email,
		a.Credit)

	return err
}

func (s *PostgresStore) GetAllAccounts() (*[]Account, error) {
	query := `select * from account`
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

func (s *PostgresStore) GetAccountById(id int) (*Account, error) {
	query := `select * from account where id = $1`
	row := s.Db.QueryRow(query, id)
	account := &Account{}
	if err := scanRow(row, account); err != nil {
		return nil, err
	}
	return account, nil
}
