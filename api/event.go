package api

import (
	"encoding/json"
	"gym_management_system/db"
	"gym_management_system/errors"
	"html/template"
	"log"
	"net/http"
)

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
	tmpl, err := template.New("events.html").ParseFiles("static/events.html")
	if err != nil {
		return err
	}
	return tmpl.Execute(w, events)
	//return writeJSON(w, http.StatusOK, map[string][]db.Event{"events": *events})
}

func (s *Server) handleCreateEvent(w http.ResponseWriter, r *http.Request) error {
	req := new(CreateEventRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return errors.InvalidRequest{Message: err.Error()}
	}
	if err := s.validator.Validate(req); err != nil {
		return err
	}

	event := db.Event{
		Type:     req.Type,
		Title:    req.Title,
		Start:    req.Start,
		End:      req.End,
		Capacity: req.Capacity,
		Price:    req.Price,
	}
	id, err := s.store.CreateEvent(&event)
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, map[string]int{"createdId": id})
}

func (s *Server) handleDeleteEvent(w http.ResponseWriter, r *http.Request) error {
	id, err := getId(r, "eventId")
	if err != nil {
		return err
	}
	err = s.store.DeleteEvent(id)
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, map[string]string{"status": "success"})
}
