package db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	customErr "gym_management_system/errors"
	"log"
	"os"
	"reflect"
)

type Storage interface {
	AccountRepository
	EntryRepository
	MembershipRepository
	AccountMembershipRepository
	EventRepository
}

type PostgresStore struct {
	Db *sql.DB
}

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

func (s *PostgresStore) Init() error {
	var err error
	//if _, err = s.Db.Exec(`drop type if exists event_type`); err != nil {
	//	return err
	//}
	//if err = s.createEventType(); err != nil {
	//	return err
	//}
	if err = s.createAccountTable(); err != nil {
		return err
	}
	if err = s.createEventTable(); err != nil {
		return err
	}
	if err = s.createMembershipTable(); err != nil {
		return err
	}
	if err = s.createAccountMembershipTable(); err != nil {
		return err
	}
	if err = s.createEntryTable(); err != nil {
		return err
	}
	return nil
}

func Close(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Fatal(err)
	}
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
	columns := createColumns(&elem)
	defer closeRows(rows)
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

func getRecord[T any](tx *sql.Tx, table Table, id int, record *T) error {
	query := fmt.Sprintf(`select * from %s where id = $1`, table)
	row := tx.QueryRow(query, id)
	if err := scanRow(row, record); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return customErr.RecordNotFound{Record: string(table), Property: "id", Value: id}
		}
		return err
	}

	v := reflect.ValueOf(record).Elem()
	deletedAtField := v.FieldByName("DeletedAt")
	if deletedAtField.IsValid() && !deletedAtField.IsNil() {
		return customErr.DeletedRecord{Record: string(table), Property: "id", Value: id}
	}
	return nil
}

func isDuplicateKeyError(err error) bool {
	var pqErr *pq.Error
	ok := errors.As(err, &pqErr)
	return ok && pqErr.Code == "23505" // PostgreSQL error code for unique violation
}
