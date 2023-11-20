package api

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"gym_management_system/db"
	customErr "gym_management_system/errors"
	"log"
	"net/http"
	"strconv"
	"time"
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
	r := mux.NewRouter()

	// Define the API router
	apiRouter := r.PathPrefix("/api").Subrouter()

	// Public Endpoints
	apiRouter.HandleFunc("/sign-up", makeHTTPHandleFunc(s.handleCreateAccount)).Methods("POST")
	apiRouter.HandleFunc("/login", makeHTTPHandleFunc(s.handleLogin)).Methods("POST")
	apiRouter.HandleFunc("/logout", makeHTTPHandleFunc(s.handleLogout)).Methods("POST")
	apiRouter.HandleFunc("/events", makeHTTPHandleFunc(s.handleGetEvents)).Methods("GET")
	apiRouter.HandleFunc("/memberships", makeHTTPHandleFunc(s.handleGetMemberships)).Methods("GET")

	// Authenticated User Endpoints with JWT authentication
	authRouter := apiRouter.PathPrefix("/").Subrouter()
	authRouter.Use(withJWTAuth)

	authRouter.HandleFunc("/account", makeHTTPHandleFunc(s.handleGetAccount)).Methods("GET")
	authRouter.HandleFunc("/account", makeHTTPHandleFunc(s.handleUpdateAccount)).Methods("PUT")
	authRouter.HandleFunc("/account", makeHTTPHandleFunc(s.handleDeleteAccount)).Methods("DELETE")
	authRouter.HandleFunc("/account/memberships", makeHTTPHandleFunc(s.handleGetAccountMemberships)).Methods("GET")
	authRouter.HandleFunc("/account/entries", makeHTTPHandleFunc(s.handleGetAccountEntries)).Methods("GET")
	authRouter.HandleFunc("/memberships/{membershipId}/purchase", makeHTTPHandleFunc(s.handleCreateAccountMembership)).Methods("POST")
	authRouter.HandleFunc("/events/{eventId}/entries", makeHTTPHandleFunc(s.handleCreateEntry)).Methods("POST")
	authRouter.HandleFunc("/events/{eventId}/entries", makeHTTPHandleFunc(s.handleGetEventEntries)).Methods("GET")
	authRouter.HandleFunc("/entries/{entryId}", makeHTTPHandleFunc(s.handleDeleteEntry)).Methods("DELETE")

	// Admin Endpoints with JWT authentication
	//adminRouter := apiRouter.PathPrefix("/admin").Subrouter()
	//adminRouter.Use(withJWTAuth)

	//adminRouter.HandleFunc("/events", makeHTTPHandleFunc(handleCreateEvent)).Methods("POST")
	//adminRouter.HandleFunc("/events/{eventId}", makeHTTPHandleFunc(handleDeleteEvent)).Methods("DELETE")
	//adminRouter.HandleFunc("/memberships", makeHTTPHandleFunc(handleCreateMembership)).Methods("POST")
	//adminRouter.HandleFunc("/memberships/{membershipId}", makeHTTPHandleFunc(handleDeleteMembership)).Methods("DELETE")

	log.Println("JSON API server running on port: ", s.listenAddr)

	if err := http.ListenAndServe(s.listenAddr, r); err != nil {
		log.Fatal(err)
	}
}

func writeErrorJSON(w http.ResponseWriter, e error) {
	var status int
	var errorMessage string
	switch {
	case errors.As(e, &customErr.ConflictingRecord{}),
		errors.As(e, &customErr.InvalidRequest{}),
		errors.As(e, &customErr.InsufficientResources{}):
		status = http.StatusBadRequest
	case errors.As(e, &customErr.PermissionDenied{}):
		status = http.StatusForbidden
	case errors.As(e, &customErr.RecordNotFound{}), errors.As(e, &customErr.DeletedRecord{}):
		status = http.StatusNotFound
	default:
		status = http.StatusInternalServerError
		errorMessage = "unknown error"
	}
	log.Println(e)
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
			log.Println("nested function failed")
			writeErrorJSON(w, err)
		}
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

func getId(r *http.Request, key string) (int, error) {
	idStr, ok := mux.Vars(r)[key]
	if !ok {
		return 0, customErr.InvalidRequest{Message: "missing " + key}
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, customErr.InvalidRequest{Message: err.Error()}
	}
	return id, nil
}

func getTime(r *http.Request, key string) (time.Time, error) {
	return time.Parse(time.RFC3339, r.URL.Query().Get(key))
}
