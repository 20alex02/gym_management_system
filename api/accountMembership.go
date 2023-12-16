package api

import (
	"encoding/json"
	"gym_management_system/db"
	"gym_management_system/errors"
	"html/template"
	"log"
	"net/http"
	"time"
)

func (s *Server) handleCreateAccountMembership(w http.ResponseWriter, r *http.Request) error {
	claims, ok := r.Context().Value("claims").(*Claims)
	if !ok {
		log.Println("cannot get claims")
		return errors.PermissionDenied{}
	}

	req := new(CreateAccountMembershipRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return errors.InvalidRequest{Message: err.Error()}
	}
	log.Println(req.ValidFrom)
	log.Println(time.Now())
	if err := s.validator.Validate(req); err != nil {
		return err
	}
	membershipId, err := getId(r, "membershipId")
	if err != nil {
		return err
	}
	membership := &db.AccountMembership{
		AccountId:    claims.Id,
		MembershipId: membershipId,
		ValidFrom:    req.ValidFrom,
	}
	id, err := s.store.CreateAccountMembership(membership)
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusCreated, map[string]int{"createdId": id})
}

func (s *Server) handleGetAccountMemberships(w http.ResponseWriter, r *http.Request) error {
	claims, ok := r.Context().Value("claims").(*Claims)
	if !ok {
		log.Println("cannot get claims")
		return errors.PermissionDenied{}
	}
	memberships, err := s.store.GetAccountMemberships(claims.Id)
	if err != nil {
		return err
	}

	tmpl, err := template.New("my_memberships.html").ParseFiles("static/my_memberships.html")
	if err != nil {
		return err
	}

	// Execute the template and write the result to the response
	return tmpl.Execute(w, memberships)
	//return writeJSON(w, http.StatusOK, map[string][]db.AccountMembershipWithType{"memberships": *memberships})
}
