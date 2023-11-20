package api

import (
	"gym_management_system/db"
	"net/http"
)

func (s *Server) handleGetMemberships(w http.ResponseWriter, _ *http.Request) error {
	memberships, err := s.store.GetAllMemberships()
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, map[string][]db.Membership{"memberships": *memberships})
}
