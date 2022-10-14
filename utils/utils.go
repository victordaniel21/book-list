package utils

import (
	"book-list/goconf"
	"book-list/models"
	"encoding/json"
	"log"
	"net/http"

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
