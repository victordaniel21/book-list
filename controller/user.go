package controller

import (
	"book-list/models"
	userRepo "book-list/repository/user"
	"book-list/utils"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (c Controller) Signup(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			user   models.User
			errMsg models.Error
		)

		// decode request body to JSON to User struct
		json.NewDecoder(r.Body).Decode(&user)

		// validation field email should not be empty
		if user.Email == "" {
			errMsg.Message = "Email should not be empty"
			utils.SendError(w, http.StatusBadRequest, errMsg)
			return
		}

		// validation field password should not be empty
		if user.Password == "" {
			errMsg.Message = "Password should not be empty"
			utils.SendError(w, http.StatusBadRequest, errMsg)
			return
		}

		// Hashing password
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
		if err != nil {
			log.Fatal(err)
		}

		// populated data hash to user.password
		user.Password = string(hash)

		// calling repository
		userRep := userRepo.UserRepository{}
		user, errUser := userRep.Signup(db, user)
		if errUser != nil {
			errMsg.Message = "Soemthing gone wrogn with sql" + errUser.Error()
			utils.SendError(w, http.StatusInternalServerError, errMsg)
			return
		}

		utils.SendSuccess(w, user)
	}
}

func (c Controller) Signin(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			user   models.User
			errMsg models.Error
			jwt    models.JWT
		)

		// decode request body from JSON to User struct
		json.NewDecoder(r.Body).Decode(&user)

		// validation user email should not be empty
		if user.Email == "" {
			errMsg.Message = "Email is missing"
			utils.SendError(w, http.StatusBadRequest, errMsg)
			return
		}

		// validation field password should not be empty
		if user.Password == "" {
			errMsg.Message = "Password is missing"
			utils.SendError(w, http.StatusBadRequest, errMsg)
			return
		}

		password := user.Password

		// calling repository
		userRep := userRepo.UserRepository{}
		user, errUser := userRep.Signin(db, user)
		if errUser != nil {
			if errUser == sql.ErrNoRows {
				errMsg.Message = "The user does not exist"
				utils.SendError(w, http.StatusInternalServerError, errMsg)
				return
			} else {
				errMsg.Message = "something gone wrong with sql" + errUser.Error()
				utils.SendError(w, http.StatusInternalServerError, errMsg)
				return
			}
		}

		hashedPassword := user.Password

		err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		if err != nil {
			errMsg.Message = "Invalid Password"
			utils.SendError(w, http.StatusUnauthorized, errMsg)
			return
		}

		// generate token
		token, err := utils.GenerateToken(user)
		if err != nil {
			errMsg.Message = "Failed generate token"
			utils.SendError(w, http.StatusInternalServerError, errMsg)
			return
		}

		jwt.Token = token
		utils.SendSuccess(w, jwt)
	}

}
