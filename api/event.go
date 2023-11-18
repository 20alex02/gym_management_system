package api

import (
	"gym_management_system/db"
	"net/http"
	"time"
)

func getWeekInterval(offset int) (from time.Time, to time.Time) {
	now := time.Now()
	daysToSubtract := int(now.Weekday()) + (offset * 7)
	startOfWeek := now.AddDate(0, 0, -daysToSubtract)
	endOfWeek := startOfWeek.AddDate(0, 0, 6)

	return startOfWeek, endOfWeek
}

func (s *Server) handleViewEvents(w http.ResponseWriter, r *http.Request) error {
	from, to := getWeekInterval(0)
	events, err := s.store.GetAllEvents(from, to)
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, map[string][]db.Event{"events": *events})
}
