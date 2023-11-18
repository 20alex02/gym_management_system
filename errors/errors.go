package errors

import "fmt"

type RecordNotFound struct {
	Id int
}

func (e RecordNotFound) Error() string {
	return fmt.Sprintf("record with Id %d not found", e.Id)
}

type AlreadyDeleted struct {
	Id int
}

func (e AlreadyDeleted) Error() string {
	return fmt.Sprintf("record with Id %d is already deleted", e.Id)
}

type ConflictingRecord struct {
	Property string
}

func (e ConflictingRecord) Error() string {
	return fmt.Sprintf("conflicting %s", e.Property)
}

type PermissionDenied struct {
	Id int
}

func (e PermissionDenied) Error() string {
	return fmt.Sprintf("permission denied")
}
