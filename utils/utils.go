package utils

import (
	"book-list/goconf"
	"book-list/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

func SendError(w http.ResponseWriter, status int, err models.Error) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(err)
}

func SendSuccess(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(data)
}

func GenerateToken(user models.User) (string, error) {
	secret := goconf.Config().GetString("jwtsecret")

	// create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"iss":   "learning",
	})

	// sign and get the complete encoded token  as astring  using the secret
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Fatal(err)
	}

	return tokenString, nil

}

func TokenVerifyMW(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		var errMsg models.Error
		authHeader := req.Header.Get("Authorization")
		bearerToken := strings.Split(authHeader, " ")
		secret := goconf.Config().GetString("jwtsecret")

		if len(bearerToken) == 2 {
			authToken := bearerToken[1]

			token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error while parsing token.")
				}
				return []byte(secret), nil
			})
			if err != nil {
				errMsg.Message = err.Error()
				SendError(resp, http.StatusUnauthorized, errMsg)
				return
			}

			if token.Valid {
				next.ServeHTTP(resp, req)
			} else {
				errMsg.Message = err.Error()
				SendError(resp, http.StatusUnauthorized, errMsg)
				return
			}
		} else {
			errMsg.Message = "Invalid Token"
			SendError(resp, http.StatusUnauthorized, errMsg)
			return

		}
	})
}
