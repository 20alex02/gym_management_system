package api

import (
	"encoding/json"
	"gym_management_system/db"
	"gym_management_system/errors"
	"html/template"
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
	if err = json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}
	if err = s.validator.Validate(req); err != nil {
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
	claims, ok := r.Context().Value("claims").(*Claims)
	if !ok {
		log.Println("cannot get claims")
		return errors.PermissionDenied{}
	}
	entryId, err := getId(r, "entryId")
	if err != nil {
		return err
	}
	if err = s.store.DeleteEntry(entryId, claims.Id); err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

func (s *Server) handleGetEventEntries(w http.ResponseWriter, r *http.Request) error {
	eventId, err := getId(r, "eventId")
	if err != nil {
		return err
	}
	entries, err := s.store.GetEventEntries(eventId)
	if err != nil {
		return err
	}
	tmpl, err := template.New("event_entries.html").ParseFiles("static/event_entries.html")
	if err != nil {
		return err
	}

	return tmpl.Execute(w, entries)
	//return writeJSON(w, http.StatusOK, map[string][]db.EventEntry{"entries": *entries})
}

func (s *Server) handleGetAccountEvents(w http.ResponseWriter, r *http.Request) error {
	claims, ok := r.Context().Value("claims").(*Claims)
	if !ok {
		log.Println("cannot get claims")
		return errors.PermissionDenied{}
	}
	events, err := s.store.GetAccountEvents(claims.Id)
	if err != nil {
		return err
	}
	tmpl, err := template.New("my_events.html").ParseFiles("static/my_events.html")
	if err != nil {
		return err
	}

	return tmpl.Execute(w, events)
	//return writeJSON(w, http.StatusOK, map[string][]db.EventWithEntryId{"events": *events})
}
