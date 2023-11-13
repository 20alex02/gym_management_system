package db

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
	"reflect"
)

type Storage interface {
	AccountRepository
	EntryRepository
	MembershipRepository
	EventRepository
}

type PostgresStore struct {
	Db *sql.DB
}

//type dbConfig struct {
//	host     string
//	user     string
//	dbName   string
//	password string
//	port     string
//	driver   string
//}

//func getPostgresConfig() (*dbConfig, error) {
//	err := godotenv.Load(".env")
//	if err != nil {
//		return nil, err
//	}
//	return &dbConfig{
//		host:     os.Getenv("DB_HOST"),
//		user:     os.Getenv("DB_USER"),
//		dbName:   os.Getenv("DB_NAME"),
//		password: os.Getenv("DB_PASSWORD"),
//		port:     os.Getenv("DB_PORT"),
//		driver:   os.Getenv("DB_DRIVER"),
//	}, nil
//}

func NewPostgresStore() (*PostgresStore, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, err
	}
	var (
		host     = os.Getenv("DB_HOST")
		user     = os.Getenv("DB_USER")
		dbName   = os.Getenv("DB_NAME")
		password = os.Getenv("DB_PASSWORD")
		port     = os.Getenv("DB_PORT")
		driver   = os.Getenv("DB_DRIVER")
		connStr  = fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s",
			host,
			user,
			dbName,
			password,
			port,
		)
	)
	db, err := sql.Open(driver, connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		Db: db,
	}, nil
}

func (s *PostgresStore) cleanup(tables []string) error {
	for i := 0; i < len(tables); i++ {
		query := "drop table if exists " + tables[i]
		if _, err := s.Db.Query(query); err != nil {
			return err
		}
	}
	_, err := s.Db.Query("drop type if exists event_type")
	return err
}

func (s *PostgresStore) createAccountTable() error {
	query := `create table if not exists account (
		id serial primary key,
		first_name varchar(255),
		last_name varchar(255),
		encrypted_password varchar(255),
		email varchar(255) unique,
		credit int,
		created_at timestamp default current_timestamp,
		deleted_at timestamp default null
	)`
	_, err := s.Db.Exec(query)
	return err
}

func (s *PostgresStore) createMembershipTable() error {
	query := `create table if not exists membership (
		id serial primary key,
		account_id int,
		type event_type,
		valid_from timestamp,
		valid_to timestamp,
		entries_left int,
		price int,
		created_at timestamp default current_timestamp,
		deleted_at timestamp default null,
		foreign key (account_id) references account(id)
	)`
	_, err := s.Db.Exec(query)
	return err
}

func (s *PostgresStore) createEntryTable() error {
	query := `create table if not exists entry (
		id serial primary key,
		account_id int,
		event_id int,
		membership_id int null,
		created_at timestamp default current_timestamp,
		deleted_at timestamp default null,
		foreign key (account_id) references account(id),
		foreign key (event_id) references event(id),
		foreign key (membership_id) references membership(id)
	)`
	_, err := s.Db.Exec(query)
	return err
}

func (s *PostgresStore) createEventTable() error {
	query := `create table if not exists event (
		id serial primary key,
		type event_type,
		title varchar(255),
		start_time timestamp,
		end_time timestamp,
		capacity int,
		price int,
		created_at timestamp default current_timestamp,
		deleted_at timestamp default null
	)`
	_, err := s.Db.Exec(query)
	return err
}

func (s *PostgresStore) createEventType() error {
	query := `create type event_type as enum (
		'open_gym',
		'lecture',
		'all'
	)`
	_, err := s.Db.Exec(query)
	return err
}

func (s *PostgresStore) Init() error {
	tables := []string{"entry", "membership", "account", "event"}
	var err error
	if err = s.cleanup(tables); err != nil {
		return err
	}

	if err = s.createEventType(); err != nil {
		return err
	}
	if err = s.createAccountTable(); err != nil {
		return err
	}
	if err = s.createEventTable(); err != nil {
		return err
	}
	if err = s.createMembershipTable(); err != nil {
		return err
	}
	if err = s.createEntryTable(); err != nil {
		return err
	}
	return nil
}

func closeRows(rows *sql.Rows) {
	err := rows.Close()
	if err != nil {
		log.Println(err)
	}
}

func createColumns(target interface{}) []interface{} {
	s := reflect.ValueOf(target).Elem()
	numCols := s.NumField()
	columns := make([]interface{}, numCols)
	for i := 0; i < numCols; i++ {
		field := s.Field(i)
		columns[i] = field.Addr().Interface()
	}
	return columns
}

func scanRows[T any](rows *sql.Rows, target *[]T) error {
	var elem T
	//s := reflect.ValueOf(&elem).Elem()
	//numCols := s.NumField()
	//columns := make([]interface{}, numCols)
	//for i := 0; i < numCols; i++ {
	//	field := s.Field(i)
	//	columns[i] = field.Addr().Interface()
	//}
	columns := createColumns(&elem)
	for rows.Next() {
		err := rows.Scan(columns...)
		if err != nil {
			return err
		}
		*target = append(*target, elem)
	}
	return nil
}

func scanRow[T any](row *sql.Row, target *T) error {
	//s := reflect.ValueOf(target).Elem()
	//numCols := s.NumField()
	//columns := make([]interface{}, numCols)
	//for i := 0; i < numCols; i++ {
	//	field := s.Field(i)
	//	columns[i] = field.Addr().Interface()
	//}
	columns := createColumns(target)
	return row.Scan(columns...)
}

func commitOrRollback(tx *sql.Tx, err *error) {
	if *err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Println(rollbackErr)
		}
		return
	}
	if commitErr := tx.Commit(); commitErr != nil {
		log.Println(commitErr)
		*err = commitErr
	}
}
