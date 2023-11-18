package api

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gym_management_system/db"
	"gym_management_system/errors"
	"net/http"
	"os"
	"time"
)

func validPassword(reqPw, encPw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(encPw), []byte(reqPw)) == nil
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
	return jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	/*
		return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return []byte(secret), nil
		})
	*/
}

func withJWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")
		if err != nil {
			writeErrorJSON(w, errors.PermissionDenied{})
			return
		}

		claims := &Claims{}
		token, err := validateJWT(c.Value, claims)
		if err != nil || !token.Valid {
			writeErrorJSON(w, errors.PermissionDenied{})
			return
		}

		ctx := context.WithValue(r.Context(), "claims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
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
