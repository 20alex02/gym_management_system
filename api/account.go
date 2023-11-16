package api

import (
	"fmt"
	"net/http"
)

func (s *Server) handleAccountById(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetAccount(w, r)
	case "DELETE":
		return s.handleDeleteAccount(w, r)
	default:
		return fmt.Errorf("method not allowed %s", r.Method)
	}
}

func (s *Server) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getId(r)
	if err != nil {
		return err
	}
	account, err := s.store.GetAccountById(id)
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, account)
}

func (s *Server) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getId(r)
	if err != nil {
		return err
	}
	if err := s.store.DeleteAccount(id); err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, map[string]int{"deleted": id})
}

func (s *Server) handleGetAccounts(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "GET" {
		return fmt.Errorf("method not allowed %s", r.Method)
	}
	accounts, err := s.store.GetAllAccounts()
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, accounts)
}

func (s *Server) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
	//if r.Method != "POST" {
	//	return fmt.Errorf("method not allowed %s", r.Method)
	//}
	//req := new(CreateAccountRequest)
	//if errors := json.NewDecoder(r.Body).Decode(req); errors != nil {
	//	return errors
	//}
	//
	//account, errors := NewAccount(req.FirstName, req.LastName, req.Password)
	//if errors != nil {
	//	return errors
	//}
	//if id, errors := s.store.CreateAccount(account); errors != nil {
	//	return errors
	//}
	//
	//return writeJSON(w, http.StatusOK, map[string]int{"created": id})
}
