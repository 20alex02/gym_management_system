package db

import (
	customErr "gym_management_system/errors"
)

type AccountMembershipRepository interface {
	CreateAccountMembership(m *AccountMembership) (int, error)
	GetAccountMemberships(id int) (*[]AccountMembershipWithType, error)
}

func (s *PostgresStore) CreateAccountMembership(m *AccountMembership) (int, error) {
	tx, err := s.Db.Begin()
	if err != nil {
		return 0, err
	}
	defer commitOrRollback(tx, &err)

	account := Account{}
	err = getRecord(tx, ACCOUNT, m.AccountId, &account)
	if err != nil {
		return 0, err
	}

	membership := Membership{}
	err = getRecord(tx, MEMBERSHIP, m.MembershipId, &membership)
	if err != nil {
		return 0, err
	}

	if account.Credit < membership.Price {
		err = customErr.InsufficientResources{}
		return 0, err
	}
	account.Credit -= membership.Price
	query := `update account set credit = $1 where id = $2`
	_, err = tx.Exec(query, account.Credit, account.Id)
	if err != nil {
		return 0, err
	}

	m.ValidTo = m.ValidFrom.AddDate(0, 0, membership.DurationDays)
	m.Entries = membership.Entries
	query = `insert into account_membership (
			account_id,
			membership_id,
			valid_from,
			valid_to,
			entries
		) values ($1, $2, $3, $4, $5) returning id`

	var id int
	err = tx.QueryRow(
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

func (s *PostgresStore) GetAccountMemberships(id int) (*[]AccountMembershipWithType, error) {
	query := `select account_membership.*, membership.type from account_membership join membership on account_membership.membership_id = membership.id 
         where account_membership.deleted_at is null 
           and account_membership.account_id = $1
           and account_membership.valid_to >= current_timestamp`
	rows, err := s.Db.Query(query, id)
	if err != nil {
		return nil, err
	}
	memberships := &[]AccountMembershipWithType{}
	if err = scanRows(rows, memberships); err != nil {
		return nil, err
	}
	return memberships, nil
}
