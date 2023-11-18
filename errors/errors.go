package errors

import "fmt"

type RecordNotFound struct {
	Entity   string
	Property string
	Value    interface{}
}

func (e RecordNotFound) Error() string {
	return fmt.Sprintf("%s with %s: %v not found", e.Entity, e.Property, e.Value)
}

type AlreadyDeleted struct {
	Entity   string
	Property string
	Value    interface{}
}

func (e AlreadyDeleted) Error() string {
	return fmt.Sprintf("%s with %s: %v is already deleted", e.Entity, e.Property, e.Value)
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
