package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/marceaudavid/learn-go/src/db"
	"github.com/marceaudavid/learn-go/src/models"
	"github.com/marceaudavid/learn-go/src/utils"
)

// Register ...
func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var body models.User
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&body)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		_, err = db.UserExist(body)
		if err == nil {
			http.Error(w, "Username is already taken", http.StatusUnauthorized)
			return
		}
		hash, _ := utils.Hash(body.Password, os.Getenv("KEY"))
		_, err = db.CreateUser(body, hash)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "OK")
	}
}

// Login ...
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var body models.User
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&body)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		user, err := db.GetUser(body)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		bool := utils.Compare(user.Password, body.Password, os.Getenv("KEY"))
		if bool {
			token, err := db.CreateToken(user, 61)
			if err != nil {
				http.Error(w, "Something went wrong", http.StatusUnauthorized)
				return
			}
			rows, err := db.DeleteExpiredTokens()
			if err != nil {
				fmt.Fprintf(w, "%d %s", &rows, err.Error())
				return
			}
			http.SetCookie(w, &http.Cookie{Name: "token", Value: *token, HttpOnly: true, Expires: time.Now().Add(61 * time.Minute)})
			fmt.Fprintf(w, "OK")
			return
		}
		http.Error(w, "Password doesn't match", http.StatusUnauthorized)
	}
}

// Logout ...
func Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "No valid token", http.StatusUnauthorized)
			return
		}
		_, err = db.DeleteTokens(cookie)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{Name: "token", HttpOnly: true, Expires: time.Now()})
		fmt.Fprintf(w, "You're logged out")
	}
}

// Load ...
func Load(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "No valid token", http.StatusUnauthorized)
			return
		}
		user := db.FromToken(cookie.Value)
		json, _ := db.GetData(user.ID)
		fmt.Fprintf(w, json)
	}
}

// Save ...
func Save(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		json := buf.String()
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "No valid token", http.StatusUnauthorized)
			return
		}
		user := db.FromToken(cookie.Value)
		_, err = db.SaveData(json, user.ID)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, json)
	}
}

// func WebsocketTicket(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == "GET" {
// 		cookie, err := r.Cookie("token")
// 		if err != nil {
// 			http.Error(w, "No valid token", http.StatusUnauthorized)
// 			return
// 		}
// 		ticket := uuid.NewV4().String()

// 	}
// }
