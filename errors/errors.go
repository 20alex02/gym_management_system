package errors

import (
	"fmt"
)

type RecordNotFound struct {
	Record   string
	Property string
	Value    interface{}
}

func (e RecordNotFound) Error() string {
	return fmt.Sprintf("%s with %s: %v not found", e.Record, e.Property, e.Value)
}

type DeletedRecord struct {
	Record   string
	Property string
	Value    interface{}
}

func (e DeletedRecord) Error() string {
	return fmt.Sprintf("%s with %s: %v is deleted", e.Record, e.Property, e.Value)
}

type ConflictingRecord struct {
	Property string
}

func (e ConflictingRecord) Error() string {
	return fmt.Sprintf("conflicting %s", e.Property)
}

type InsufficientResources struct {
}

func (e InsufficientResources) Error() string {
	return fmt.Sprintf("insufficient resources")
}

type PermissionDenied struct {
	Id int
}

func (e PermissionDenied) Error() string {
	return fmt.Sprintf("permission denied")
}

type InvalidRequest struct {
	Message string
}

func (e InvalidRequest) Error() string {
	return e.Message
}
