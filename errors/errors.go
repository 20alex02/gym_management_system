package errors

import (
	"fmt"
	"gym_management_system/db"
)

type RecordNotFound struct {
	Record   db.Table
	Property string
	Value    interface{}
}

func (e RecordNotFound) Error() string {
	return fmt.Sprintf("%s with %s: %v not found", e.Record, e.Property, e.Value)
}

type DeletedRecord struct {
	Record   db.Table
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

type EntryError struct {
	Message string
}

func (e EntryError) Error() string {
	return e.Message
}

type PermissionDenied struct {
	Id int
}

func (e PermissionDenied) Error() string {
	return fmt.Sprintf("permission denied")
}

type InvalidRequestFormat struct {
	Message string
}

func (e InvalidRequestFormat) Error() string {
	return fmt.Sprintf("invalid request format: %s", e.Message)
}
