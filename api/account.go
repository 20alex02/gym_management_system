package api

import (
	"encoding/json"
	"gym_management_system/db"
	"gym_management_system/errors"
	"log"
	"net/http"
	"time"
)

func (s *Server) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	// TODO validation --> names have to be 3-30 characters long, email has to be in valid format, password has to contain at least 6 characters of which at least one is lowercase, uppercase, number, special character
	req := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	pw, err := HashPassword(req.Password)
	account := db.Account{
		FirstName:         req.FirstName,
		LastName:          req.LastName,
		EncryptedPassword: pw,
		Email:             req.Email,
	}
	if err != nil {
		return err
	}
	id, err := s.store.CreateAccount(&account)
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
		Name:     "token",
		Value:    token,
		Expires:  expTime,
		HttpOnly: true,
	})

	return writeJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

func (s *Server) handleLogout(w http.ResponseWriter, _ *http.Request) error {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
	})
	return writeJSON(w, http.StatusOK, map[string]string{"status": "success"})
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

	return writeJSON(w, http.StatusOK, map[string]string{"status": "success"})
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
	pw, err := HashPassword(req.Password)
	if err != nil {
		return err
	}
	account := db.Account{
		Id:                claims.Id,
		FirstName:         req.FirstName,
		LastName:          req.LastName,
		EncryptedPassword: pw,
		Email:             req.Email,
		Credit:            req.Credit,
	}
	if err := s.store.UpdateAccount(&account); err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, map[string]string{"status": "success"})
}
