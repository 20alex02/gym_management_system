package api

import (
	"encoding/json"
	"gym_management_system/db"
	"net/http"
)

//func getWeekInterval(offset int) (from time.Time, to time.Time) {
//	now := time.Now()
//	daysToSubtract := int(now.Weekday()) + (offset * 7)
//	startOfWeek := now.AddDate(0, 0, -daysToSubtract)
//	endOfWeek := startOfWeek.AddDate(0, 0, 6)
//
//	return startOfWeek, endOfWeek
//}

func (s *Server) handleGetEvents(w http.ResponseWriter, r *http.Request) error {
	req := new(GetEventsRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}
	events, err := s.store.GetAllEvents(req.From, req.To)
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, map[string][]db.Event{"events": *events})
}
