package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"gym_management_system/db"
	customErr "gym_management_system/errors"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Server struct {
	listenAddr string
	store      db.Storage
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type Error struct {
	Error string `json:"error"`
}

func NewServer(listenAddr string, store db.Storage) *Server {
	return &Server{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *Server) Run() {
	router := mux.NewRouter()

	//router.HandleFunc("/login", makeHTTPHandleFunc(s.handleLogin))
	router.HandleFunc("/sign-up", makeHTTPHandleFunc(s.handleCreateAccount))
	router.HandleFunc("/accounts", makeHTTPHandleFunc(s.handleGetAccounts))
	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleAccountById))

	//router.HandleFunc("/account/{id}/memberships", makeHTTPHandleFunc(s.handleGetAccountMemberships))
	// get all available, create new
	//router.HandleFunc("/memberships", makeHTTPHandleFunc(s.handleMemberships))
	// get all
	//router.HandleFunc("/events", makeHTTPHandleFunc(s.handleEvents))
	//router.HandleFunc("/event", makeHTTPHandleFunc(s.handleCreateEvent))
	//router.HandleFunc("/event/{id}", makeHTTPHandleFunc(s.handleGetEvent))
	// create delete
	//router.HandleFunc("/event/{id}/entry", makeHTTPHandleFunc(s.handleEntry))

	log.Println("JSON API server running on port: ", s.listenAddr)

	if err := http.ListenAndServe(s.listenAddr, router); err != nil {
		log.Fatal(err)
	}
}

func writeErrorJSON(w http.ResponseWriter, e error) {
	var status int
	var errorMessage string
	switch {
	case errors.Is(e, customErr.RecordNotFound{}), errors.Is(e, customErr.AlreadyDeleted{}):
		status = http.StatusNotFound
	case errors.Is(e, customErr.PermissionDenied{}):
		status = http.StatusForbidden
	default:
		status = http.StatusInternalServerError
		errorMessage = "unknown error"
	}
	if errorMessage == "" {
		errorMessage = e.Error()
	}
	err := writeJSON(w, status, Error{Error: errorMessage})
	if err != nil {
		log.Println(err)
	}
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			writeErrorJSON(w, err)
		}
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

func getId(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("invalid id given %s", idStr)
	}
	return id, nil
}

//func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
//	if r.Method != "POST" {
//		return fmt.Errorf("method not allowed %s", r.Method)
//	}
//
//	var req LoginRequest
//	if errors := json.NewDecoder(r.Body).Decode(&req); errors != nil {
//		return errors
//	}
//
//	acc, errors := s.store.GetAccountByNumber(int(req.Number))
//	if errors != nil {
//		return errors
//	}
//
//	if !acc.ValidPassword(req.Password) {
//		return fmt.Errorf("not authenticated")
//	}
//
//	token, errors := createJWT(acc)
//	if errors != nil {
//		return errors
//	}
//
//	resp := LoginResponse{
//		Token:  token,
//		Number: acc.Number,
//	}
//
//	return writeJSON(w, http.StatusOK, resp)
//}

func createJWT(account *db.Account) (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt": 15000,
		"id":        account.Id,
	}

	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

//func withJWTAuth(handlerFunc http.HandlerFunc, s db.Storage) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		fmt.Println("calling JWT auth middleware")
//
//		tokenString := r.Header.Get("x-jwt-token")
//		token, err := validateJWT(tokenString)
//		if err != nil || !token.Valid {
//			writeErrorJSON(w, customErr.PermissionDenied{})
//			return
//		}
//
//		userID, err := getId(r)
//		if err != nil {
//			writeErrorJSON(w, customErr.PermissionDenied{})
//			return
//		}
//		account, err := s.GetAccountByID(userID)
//		if err != nil {
//			writeErrorJSON(w, customErr.PermissionDenied{})
//			return
//		}
//
//		claims := token.Claims.(jwt.MapClaims)
//		if account.Number != int64(claims["accountNumber"].(float64)) {
//			writeErrorJSON(w, customErr.PermissionDenied{})
//			return
//		}
//
//		handlerFunc(w, r)
//	}
//}

func validateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})
}
