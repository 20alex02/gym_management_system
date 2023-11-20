package db

import (
	"database/sql"
	"errors"
	customErr "gym_management_system/errors"
)

type AccountRepository interface {
	CreateAccount(a *Account) (int, error)
	GetAccountById(id int) (*Account, error)
	GetAccountByEmail(email string) (*Account, error)
	//GetAllAccounts() (*[]Account, error)
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
		if isDuplicateKeyError(err) {
			return 0, customErr.ConflictingRecord{Property: "email"}
		}
		return 0, err
	}
	return id, nil
}

func (s *PostgresStore) GetAccountById(id int) (*Account, error) {
	query := `select * from account where id = $1`
	row := s.Db.QueryRow(query, id)
	account := &Account{}
	if err := scanRow(row, account); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, customErr.RecordNotFound{Record: string(ACCOUNT), Property: "id", Value: id}
		}
		return nil, err
	}
	if account.DeletedAt != nil {
		return nil, customErr.DeletedRecord{Record: string(ACCOUNT), Property: "id", Value: id}
	}
	return account, nil
}

func (s *PostgresStore) GetAccountByEmail(email string) (*Account, error) {
	query := `select * from account where email = $1`
	row := s.Db.QueryRow(query, email)
	account := &Account{}
	if err := scanRow(row, account); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, customErr.RecordNotFound{Record: string(ACCOUNT), Property: "email", Value: email}
		}
		return nil, err
	}
	if account.DeletedAt != nil {
		return nil, customErr.DeletedRecord{Record: string(ACCOUNT), Property: "email", Value: email}
	}
	return account, nil
}

/*
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
*/

func (s *PostgresStore) UpdateAccount(a *Account) error {
	tx, err := s.Db.Begin()
	if err != nil {
		return err
	}
	defer commitOrRollback(tx, &err)

	account := Account{}
	err = checkRecord(tx, ACCOUNT, a.Id, &account)
	if err != nil {
		return err
	}
	query := `update account set
                   first_name = $1,
                   last_name = $2,
                   encrypted_password = $3,
                   email = $4,
                   credit = $5
               where id = $6`
	_, err = tx.Exec(query, a.FirstName, a.LastName, a.EncryptedPassword, a.Email, a.Credit, a.Id)
	if isDuplicateKeyError(err) {
		return customErr.ConflictingRecord{Property: "email"}
	}
	return err
}

func (s *PostgresStore) DeleteAccount(id int) error {
	tx, err := s.Db.Begin()
	if err != nil {
		return err
	}
	defer commitOrRollback(tx, &err)

	account := Account{}
	err = checkRecord(tx, ACCOUNT, id, &account)
	if err != nil {
		return err
	}

	query := `update account set deleted_at = current_timestamp where id = $1`
	_, err = tx.Exec(query, id)
	if err != nil {
		return err
	}

	query = `update entry set deleted_at = current_timestamp 
             where account_id = $1 and deleted_at is null`
	_, err = tx.Exec(query, id)
	if err != nil {
		return err
	}

	query = `update account_membership set deleted_at = current_timestamp 
                          where account_id = $1 and deleted_at is null`
	_, err = tx.Exec(query, id)
	return err
}
