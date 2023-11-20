package api

import (
	"encoding/json"
	"gym_management_system/errors"
	"log"
	"net/http"
	"time"
)

func (s *Server) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	// TODO validation
	req := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	account, err := newAccount(req.FirstName, req.LastName, req.Email, req.Password)
	if err != nil {
		return err
	}
	id, err := s.store.CreateAccount(account)
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusCreated, map[string]int{"createdId": id})
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) error {
	req := new(LoginRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	acc, err := s.store.GetAccountByEmail(req.Email)
	if err != nil {
		return err
	}

	if !validPassword(req.Password, acc.EncryptedPassword) {
		return errors.PermissionDenied{}
	}

	expTime := time.Now().Add(time.Minute * 15)
	token, err := createJWT(acc, expTime)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: expTime,
	})

	return writeJSON(w, http.StatusOK, map[string]string{"message": "success"})
}

func (s *Server) handleLogout(w http.ResponseWriter, _ *http.Request) error {
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Expires: time.Now(),
	})
	return writeJSON(w, http.StatusOK, map[string]string{"message": "success"})
}

func (s *Server) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	claims, ok := r.Context().Value("claims").(*Claims)
	if !ok {
		log.Println("cannot get claims")
		return errors.PermissionDenied{}
	}
	account, err := s.store.GetAccountById(claims.Id)
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, account)
}

func (s *Server) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	claims, ok := r.Context().Value("claims").(*Claims)
	if !ok {
		log.Println("cannot get claims")
		return errors.PermissionDenied{}
	}
	if err := s.store.DeleteAccount(claims.Id); err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, map[string]string{"message": "success"})
}

func (s *Server) handleUpdateAccount(w http.ResponseWriter, r *http.Request) error {
	claims, ok := r.Context().Value("claims").(*Claims)
	if !ok {
		log.Println("cannot get claims")
		return errors.PermissionDenied{}
	}

	req := new(UpdateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}
	account, err := newAccount(req.FirstName, req.LastName, req.Email, req.Password)
	if err != nil {
		return err
	}
	account.Id = claims.Id
	account.Credit = req.Credit

	if err := s.store.UpdateAccount(account); err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, map[string]string{"message": "success"})
}

//func (s *Server) handleGetAccounts(w http.ResponseWriter, r *http.Request) error {
//	if r.Method != "GET" {
//		return fmt.Errorf("method not allowed %s", r.Method)
//	}
//	accounts, err := s.store.GetAllAccounts()
//	if err != nil {
//		return err
//	}
//
//	return writeJSON(w, http.StatusOK, accounts)
//}

//func (s *Server) handleAccountById(w http.ResponseWriter, r *http.Request) error {
//	switch r.Method {
//	case "GET":
//		return s.handleGetAccount(w, r)
//	case "DELETE":
//		return s.handleDeleteAccount(w, r)
//	default:
//		return fmt.Errorf("method not allowed %s", r.Method)
//	}
//}
