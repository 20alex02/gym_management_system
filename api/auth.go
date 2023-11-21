package api

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gym_management_system/db"
	"gym_management_system/errors"
	"log"
	"net/http"
	"os"
	"time"
)

func validPassword(reqPw, encPw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(encPw), []byte(reqPw)) == nil
}

func HashPassword(password string) (string, error) {
	pw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(pw), err
}

func createJWT(account *db.Account, exp time.Time) (string, error) {
	claims := &Claims{
		Id: account.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	return token.SignedString([]byte(secret))
}

func validateJWT(tokenString string, claims jwt.Claims) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	log.Println(secret)
	log.Println("validating jwt")
	return jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
}

func withJWTAuth(next http.Handler) http.Handler {
	log.Println("jwt auth middleware")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")
		if err != nil {
			log.Println("cookie error", err)
			writeErrorJSON(w, errors.PermissionDenied{})
			return
		}

		claims := &Claims{}
		token, err := validateJWT(c.Value, claims)
		if err != nil || !token.Valid {
			log.Println("validation of jwt failed")
			writeErrorJSON(w, errors.PermissionDenied{})
			return
		}

		ctx := context.WithValue(r.Context(), "claims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) isAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value("claims").(*Claims)
		if !ok {
			log.Println("cannot get claims")
			writeErrorJSON(w, errors.PermissionDenied{})
			return
		}

		account, err := s.store.GetAccountById(claims.Id)
		if err != nil {
			log.Println("error getting account:", err)
			writeErrorJSON(w, errors.PermissionDenied{})
			return
		}

		if account.Role != db.ADMIN {
			log.Println("user is not admin")
			writeErrorJSON(w, errors.PermissionDenied{})
			return
		}

		next.ServeHTTP(w, r)
	})
}

/*
func Refresh(w http.ResponseWriter, r *http.Request) {
	// (BEGIN) The code until this point is the same as the first part of the `Welcome` route
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tknStr := c.Value
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (any, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// (END) The code until this point is the same as the first part of the `Welcome` route

	// We ensure that a new token is not issued until enough time has elapsed
	// In this case, a new token will only be issued if the old token is within
	// 30 seconds of expiry. Otherwise, return a bad request status
	if time.Until(claims.ExpiresAt.Time) > 30*time.Second {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Now, create a new token for the current use, with a renewed expiration time
	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the new token as the users `token` cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
}

*/
