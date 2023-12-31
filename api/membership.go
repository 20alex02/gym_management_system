package api

import (
	"encoding/json"
	"gym_management_system/db"
	"gym_management_system/errors"
	"html/template"
	"net/http"
)

func (s *Server) handleGetMemberships(w http.ResponseWriter, _ *http.Request) error {
	memberships, err := s.store.GetAllMemberships()
	if err != nil {
		return err
	}

	tmpl, err := template.New("memberships.html").ParseFiles("static/memberships.html")
	if err != nil {
		return err
	}

	return tmpl.Execute(w, memberships)
	//return writeJSON(w, http.StatusOK, map[string][]db.Membership{"memberships": *memberships})
}

func (s *Server) handleCreateMembership(w http.ResponseWriter, r *http.Request) error {
	req := new(CreateMembershipRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return errors.InvalidRequest{Message: err.Error()}
	}
	if err := s.validator.Validate(req); err != nil {
		return err
	}
	membership := db.Membership{
		Type:         req.Type,
		DurationDays: req.DurationDays,
		Entries:      req.Entries,
		Price:        req.Price,
	}
	id, err := s.store.CreateMembership(&membership)
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, map[string]int{"createdId": id})
}

func (s *Server) handleDeleteMembership(w http.ResponseWriter, r *http.Request) error {
	id, err := getId(r, "membershipId")
	if err != nil {
		return err
	}
	err = s.store.DeleteMembership(id)
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, map[string]string{"createdId": "success"})
}
