package api

import (
	"encoding/json"
	"gym_management_system/db"
	"gym_management_system/errors"
	"log"
	"net/http"
)

func (s *Server) handleCreateEntry(w http.ResponseWriter, r *http.Request) error {
	claims, ok := r.Context().Value("claims").(*Claims)
	if !ok {
		log.Println("cannot get claims")
		return errors.PermissionDenied{}
	}

	eventId, err := getId(r, "eventId")
	if err != nil {
		return err
	}

	req := new(CreateEntryRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	entry := &db.Entry{
		AccountId:           claims.Id,
		EventId:             eventId,
		AccountMembershipId: req.AccountMembershipId,
	}
	id, err := s.store.CreateEntry(entry)
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusCreated, map[string]int{"createdId": id})
}

func (s *Server) handleDeleteEntry(w http.ResponseWriter, r *http.Request) error {
	_, ok := r.Context().Value("claims").(*Claims)
	if !ok {
		log.Println("cannot get claims")
		return errors.PermissionDenied{}
	}

	entryId, err := getId(r, "entryId")
	if err != nil {
		return err
	}
	if err := s.store.DeleteEntry(entryId); err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, map[string]string{"message": "success"})
}

func (s *Server) handleGetEventEntries(w http.ResponseWriter, r *http.Request) error {
	_, ok := r.Context().Value("claims").(*Claims)
	if !ok {
		log.Println("cannot get claims")
		return errors.PermissionDenied{}
	}

	eventId, err := getId(r, "eventId")
	if err != nil {
		return err
	}
	entries, err := s.store.GetEventEntries(eventId)
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, map[string][]db.Entry{"entries": *entries})
}

func (s *Server) handleGetAccountEntries(w http.ResponseWriter, r *http.Request) error {
	claims, ok := r.Context().Value("claims").(*Claims)
	if !ok {
		log.Println("cannot get claims")
		return errors.PermissionDenied{}
	}
	entries, err := s.store.GetAccountEntries(claims.Id)
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, map[string][]db.Entry{"entries": *entries})
}
