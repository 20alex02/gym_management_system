package api

import (
	"gym_management_system/db"
	"gym_management_system/errors"
	"log"
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
	from, err := getTime(r, "from")
	if err != nil {
		return errors.InvalidRequest{Message: err.Error()}
	}
	to, err := getTime(r, "to")
	if err != nil {
		log.Println("invalid to format")
		return errors.InvalidRequest{Message: err.Error()}
	}

	log.Println("req parsed")
	log.Println(from, to)
	events, err := s.store.GetAllEvents(from, to)
	log.Println("events retrieved from db")
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, map[string][]db.Event{"events": *events})
}
