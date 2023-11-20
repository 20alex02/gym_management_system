package api

import (
	"encoding/json"
	"gym_management_system/db"
	"gym_management_system/errors"
	"log"
	"net/http"
)

func (s *Server) handleCreateAccountMembership(w http.ResponseWriter, r *http.Request) error {
	claims, ok := r.Context().Value("claims").(*Claims)
	if !ok {
		log.Println("cannot get claims")
		return errors.PermissionDenied{}
	}

	// TODO validation, valid from < today
	req := new(CreateAccountMembershipRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
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
	memberships, err := s.store.GetAllAccountMemberships(claims.Id)
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, map[string][]db.AccountMembership{"memberships": *memberships})
}
